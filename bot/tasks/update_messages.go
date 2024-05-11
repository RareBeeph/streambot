package tasks

import (
	"fmt"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const maxTimesFailed = 5

type msgAction struct {
	target  *models.Message
	content []*discordgo.MessageEmbed
}

func updateMessages(s *discordgo.Session) {
	m := query.Message
	qs := query.Subscription

	subscriptions, err := qs.Preload(qs.Messages.Order(m.PostOrder.Asc())).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

subscriptionLoop:
	for _, sub := range subscriptions {
		if sub.TimesFailed >= maxTimesFailed {
			continue
		}

		actions := getActions(sub)

		for idx, action := range actions {
			// requirements: idx, action, sub

			// Perform post/edit and update database
			if action.target == nil {
				// no target => post
				err = postMessage(s, sub, action, idx)
				if err != nil {
					log.Err(err).Msg("Failed to send message.")
					qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, sub.TimesFailed+1)
					continue subscriptionLoop
				}
			} else {
				// yes target => edit
				_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					Embeds: action.content,

					ID:      action.target.MessageID,
					Channel: sub.ChannelID,
				})
				if err != nil {
					log.Err(err).Msg("Failed to edit message.")
					qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, sub.TimesFailed+1)
					continue subscriptionLoop
				}
			}
		}

		if len(sub.Messages) > len(actions) {
			// bulk delete unneeded messages and update database
			messagesToDelete := util.Map(sub.Messages[len(actions):], func(message models.Message, idx int) string {
				return message.MessageID
			})
			err := s.ChannelMessagesBulkDelete(sub.ChannelID, messagesToDelete)
			if err != nil {
				log.Err(err).Msg("Failed to bulk delete messages")
				qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, sub.TimesFailed+1)
				continue subscriptionLoop // could just be "continue" but might as well extend the pattern
			}

			m.Where(m.MessageID.In(messagesToDelete...)).Delete()
		}

		// if all our posting/editing/deleting succeeded
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, 0)
	}
}

func getActions(sub *models.Subscription) []*msgAction {
	qst := query.Stream

	matchingStreams, err := qst.Where(qst.GameID.Eq(sub.GameID), qst.Title.Lower().Like(fmt.Sprintf("%%%s%%", sub.Filter))).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find matching streams.")
	}

	// The modifications we're going to make to our messages
	actions := []*msgAction{}

	embedFields := util.Chunk(StreamsToEmbedFields(matchingStreams...), 25)
	embeds := util.Map(embedFields, StreamsMessageEmbed)
	messageChunks := util.Chunk(embeds, 10)

	messageCount := len(sub.Messages)
	for idx, embed := range messageChunks {
		// Determine what action needs to be taken to post this chunk
		if idx < messageCount {
			actions = append(actions, &msgAction{target: &sub.Messages[idx], content: embed})
		} else {
			actions = append(actions, &msgAction{content: embed})
		}
	}

	return actions
}

func postMessage(s *discordgo.Session, sub *models.Subscription, action *msgAction, postOrderIdx int) error {
	m := query.Message

	message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
		Content: "placeholder",
		Embeds:  action.content,
	})
	m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: postOrderIdx})

	return err
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
		Color:       0x9922cc,
		Author:      &discordgo.MessageEmbedAuthor{},
		Title:       "placeholder",
		Description: "placeholder",
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}
