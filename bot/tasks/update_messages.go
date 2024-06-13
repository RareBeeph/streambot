package tasks

import (
	"errors"
	"fmt"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type msgAction struct {
	target  *models.Message
	content []*discordgo.MessageEmbed
}

func updateMessages(s *discordgo.Session, minHealth int, maxHealth int) {
	m := query.Message
	qs := query.Subscription

	// Load all active subscriptions
	subscriptions, err := qs.
		Preload(qs.Messages.Order(m.PostOrder.Asc())).
		GetByHealth(minHealth, maxHealth)
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	var wg sync.WaitGroup
	wg.Add(len(subscriptions))
	for _, sub := range subscriptions {
		go func() {
			defer wg.Done()
			performUpdates(s, sub)
		}()
	}

	wg.Wait()
}

func performUpdates(s *discordgo.Session, sub *models.Subscription) {
	m := query.Message
	qs := query.Subscription
	qst := query.Stream
	qb := query.BlacklistEntry

	matchingBlacklists, err := qb.Where(qb.ChannelID.Eq(sub.ChannelID)).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find matching blacklists.")
	}

	blacklistUserIDs := util.Map(matchingBlacklists, func(bl *models.BlacklistEntry, idx int) string {
		return bl.UserID
	})

	streamquery := qst.Where(
		qst.GameID.Eq(sub.GameID),
		qst.Title.Lower().Like("%"+sub.Filter+"%"),
		qst.UserID.NotIn(blacklistUserIDs...))
	if sub.Language != "" {
		streamquery = streamquery.Where(qst.Language.Eq(sub.Language))
	}
	matchingStreams, err := streamquery.Find()

	if err != nil {
		log.Err(err).Msg("Failed to find matching streams.")
	}

	// The modifications we're going to make to our messages
	actions := []*msgAction{}

	// The actual content of the messages we intend to post
	embedFields := util.Chunk(StreamsToEmbedFields(matchingStreams...), 25)
	embeds := util.Map(embedFields, func(fields []*discordgo.MessageEmbedField, idx int) *discordgo.MessageEmbed {
		out := StreamsMessageEmbed(fields, idx)
		out.Title = fmt.Sprintf("Streams for %s", sub)
		return out
	})
	messageChunks := util.Chunk(embeds, 2)

	// Temp: limit to 1 message per subscription for aesthetic purposes
	// TODO: reconsider this decision.
	if len(messageChunks) > 1 {
		messageChunks = messageChunks[:1]
	}

	// Determine what action needs to be taken to post each chunk
	messageCount := len(sub.Messages)
	for idx, embed := range messageChunks {
		if idx < messageCount {
			actions = append(actions, &msgAction{target: &sub.Messages[idx], content: embed})
		} else {
			actions = append(actions, &msgAction{content: embed})
		}
	}
	// Determine which messages need to be deleted
	if len(messageChunks) < messageCount {
		for _, message := range sub.Messages[len(messageChunks):] {
			actions = append(actions, &msgAction{target: &message})
		}
	}

	errored := false
	for idx, action := range actions {
		var err error
		// Perform post/edit/delete and update database
		if action.target == nil {
			// no target => post
			message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
				Embeds: action.content,
			})
			m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: idx})

			if err != nil {
				log.Err(err).Msg("Failed to send message.")
			}
		} else if action.content == nil {
			// no content => delete
			err = s.ChannelMessageDelete(sub.ChannelID, action.target.MessageID)

			var resterr *discordgo.RESTError
			if err == nil || (errors.As(err, &resterr) && resterr.Message.Code == discordgo.ErrCodeUnknownMessage) {
				// If we successfully deleted the message, or if it had already been deleted, remove our record of it
				m.Where(m.MessageID.Eq(action.target.MessageID)).Delete()
			} else {
				log.Err(err).Msg("Failed to delete message.")
			}
		} else {
			// yes taget and content => edit
			_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Embeds: &action.content,

				ID:      action.target.MessageID,
				Channel: sub.ChannelID,
			})

			var resterr *discordgo.RESTError
			if errors.As(err, &resterr) && resterr.Message.Code == discordgo.ErrCodeUnknownMessage {
				// if the target message no longer exists, remove our record of it
				m.Where(m.MessageID.Eq(action.target.MessageID)).Delete()
			}

			// not an else if like in the delete case, because the intended edit never makes it to the user
			if err != nil {
				errored = true
				log.Err(err).Msg("Failed to edit message.")
			}
		}
	}

	if errored {
		// propagate failure count
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
		return
	}

	// if all our posting/editing/deleting succeeded
	qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, 0)
}

func StreamsToEmbedFields(streams ...*models.Stream) []*discordgo.MessageEmbedField {
	out := []*discordgo.MessageEmbedField{}

	for _, s := range streams {
		link := fmt.Sprintf("https://twitch.tv/%s", s.UserName)
		title := util.TruncateString(s.Title, 100-len([]rune(link)))

		out = append(out, &discordgo.MessageEmbedField{
			Name:   title,
			Value:  link,
			Inline: true,
		})
	}

	return out
}

func StreamsMessageEmbed(fields []*discordgo.MessageEmbedField, idx int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     0x9922cc,
		Author:    &discordgo.MessageEmbedAuthor{},
		Title:     "placeholder",
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
