package tasks

import (
	"errors"
	"fmt"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type msgAction struct {
	target  *models.Message
	content []*discordgo.MessageEmbed
}

func updateMessages(s *discordgo.Session) {
	m := query.Message
	qs := query.Subscription

	// Load all active subscriptions
	subscriptions, err := qs.
		Preload(qs.Messages.Order(m.PostOrder.Asc())).
		GetByHealth(models.MaxTimesFailed)
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	for _, sub := range subscriptions {
		orphanedMsgs, err := performUpdates(s, sub)
		if err != nil {
			continue
		}

		// bulk delete unneeded messages and update database
		messagesToDelete := util.Map(orphanedMsgs, func(message models.Message, idx int) string {
			return message.MessageID
		})
		err = s.ChannelMessagesBulkDelete(sub.ChannelID, messagesToDelete)
		if err != nil {
			log.Err(err).Msg("Failed to bulk delete messages")
			// TODO: testing and recovery logic
			qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
			continue
		}
		m.Where(m.MessageID.In(messagesToDelete...)).Delete()

		// if all our posting/editing/deleting succeeded
		sub.TimesFailed = 0
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
	}
}

func performUpdates(s *discordgo.Session, sub *models.Subscription) ([]models.Message, error) {
	m := query.Message
	qs := query.Subscription
	qst := query.Stream

	matchingStreams, err := qst.Where(qst.GameID.Eq(sub.GameID), qst.Title.Lower().Like(fmt.Sprintf("%%%s%%", sub.Filter))).Find()
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
	messageChunks := util.Chunk(embeds, 10)

	// Determine what action needs to be taken to post each chunk
	messageCount := len(sub.Messages)
	for idx, embed := range messageChunks {
		if idx < messageCount {
			actions = append(actions, &msgAction{target: &sub.Messages[idx], content: embed})
		} else {
			actions = append(actions, &msgAction{content: embed})
		}
	}

	errored := false
	for idx, action := range actions {
		// Perform post/edit and update database
		if action.target == nil {
			// no target => post
			message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
				// Content: "placeholder",
				Embeds: action.content,
			})
			m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: idx})

			if err != nil {
				log.Err(err).Msg("Failed to send message.")
				errored = true
			}
		} else {
			// yes target => edit
			_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Embeds: action.content,

				ID:      action.target.MessageID,
				Channel: sub.ChannelID,
			})
			if err != nil {
				// janky check to catch message records that no longer refer to existent messages
				if strings.Contains(err.Error(), fmt.Sprint(discordgo.ErrCodeUnknownMessage)) {
					m.Where(m.ID.Eq(action.target.ID)).Delete()
				}
				log.Err(err).Msg("Failed to edit messages.")
				errored = true
			}

		}

	}

	if errored {
		// Propagate failure count
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
		return []models.Message{}, errors.New("failed to update at least one message")
	} else {
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, 0)
	}

	if len(sub.Messages) > len(actions) {
		return sub.Messages[len(actions):], nil
	}

	return []models.Message{}, nil
}

func StreamsToEmbedFields(streams ...*models.Stream) []*discordgo.MessageEmbedField {
	out := []*discordgo.MessageEmbedField{}

	for _, s := range streams {
		out = append(out, &discordgo.MessageEmbedField{
			Name:   s.Title,
			Value:  fmt.Sprintf("https://twitch.tv/%s", s.UserName),
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
