package messages

import (
	"fmt"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func UpdateMessages(s *discordgo.Session, minHealth int, maxHealth int) {
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
		go func(sub *models.Subscription) {
			defer wg.Done()
			performUpdates(s, sub)
		}(sub)
	}

	wg.Wait()
}

func performUpdates(s *discordgo.Session, sub *models.Subscription) {
	qs := query.Subscription

	matchingStreams, err := relevantStreams(sub)
	if err != nil {
		return
	}

	actions := arrangeMsgActions(matchingStreams, sub)

	var errored bool
	for idx, action := range actions {
		// Perform post/edit/delete and update database
		if action.target == nil {
			errored = action.postContent(s, sub, idx)
		} else if action.content == nil {
			errored = action.deleteTarget(s, sub)
		} else {
			// if action has a target and content, edit that target with that content
			errored = action.editMessage(s, sub)
		}
	}

	if errored {
		// propagate failure count
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
		return
	}

	// on success, reset failure count
	qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, 0)
}

func relevantStreams(sub *models.Subscription) ([]*models.Stream, error) {
	qb := query.BlacklistEntry
	qst := query.Stream

	matchingBlacklists, err := qb.Where(qb.ChannelID.Eq(sub.ChannelID)).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find matching blacklists.")
		return []*models.Stream{}, err
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
		return []*models.Stream{}, err
	}

	return matchingStreams, nil
}

func arrangeMsgActions(streams []*models.Stream, sub *models.Subscription) (actions []*msgAction) {
	// The actual content of the messages we intend to post
	embedFields := util.Chunk(streamsToEmbedFields(streams...), 25)
	embeds := util.Map(embedFields, func(fields []*discordgo.MessageEmbedField, idx int) *discordgo.MessageEmbed {
		out := fieldsToMessageEmbed(fields)
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

	return
}

func streamsToEmbedFields(streams ...*models.Stream) []*discordgo.MessageEmbedField {
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

func fieldsToMessageEmbed(fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     0x9922cc,
		Author:    &discordgo.MessageEmbedAuthor{},
		Title:     "placeholder",
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
