package tasks

import (
	"fmt"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func updateMessages(s *discordgo.Session, streams []*models.Stream) {
	m := query.Message
	qs := query.Subscription

	subscriptions, err := qs.Preload(qs.Messages).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	for _, sub := range subscriptions {
		matchingStreams := []*models.Stream{}
		for _, st := range streams {
			if StreamMatchesSubscription(st, sub) {
				matchingStreams = append(matchingStreams, st)
			}
		}

		// naive edit
		if len(sub.Messages) != 0 {
			embedfields := util.Chunk(StreamsToEmbedFields(matchingStreams...), 25)
			embeds := util.Map(embedfields, StreamsMessageEmbed)
			messagechunks := util.Chunk(embeds, 10)

			// empty case
			if len(matchingStreams) == 0 {
				s.ChannelMessageEdit(sub.ChannelID, sub.Messages[0].MessageID, "No streams currently active.")
			}

			// remove excessive messages, else post needed ones
			if len(sub.Messages) > len(messagechunks) {
				s.ChannelMessagesBulkDelete(sub.ChannelID, util.Map(sub.Messages[len(messagechunks):], func(message models.Message, i int) string {
					id := message.MessageID
					m.Delete(&message)
					return id
				}))
			} else if len(sub.Messages) < len(messagechunks) {
				for idx, me := range messagechunks[len(sub.Messages):] {
					newMessage, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
						Content: "Test Edit Send",
						Embeds:  me,
					})
					if err != nil {
						log.Err(err).Msg("Failed to send message in edit branch.")
					}

					m.Create(&models.Message{MessageID: newMessage.ID, SubscriptionID: sub.ID, PostOrder: idx + len(sub.Messages)})
				}
			}
			sub.Messages = sub.Messages[:len(messagechunks)]

			// edit existing messages
			for i, c := range sub.Messages {
				content := "Test Edit"
				_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					Content: &content,
					Embeds:  messagechunks[i],

					ID:      c.MessageID,
					Channel: sub.ChannelID,
				})
				if err != nil {
					log.Err(err).Msg("Failed to edit message.")
				}
			}

			// we're done here
			return
		}

		// don't try to post nothing
		if len(matchingStreams) == 0 {
			return
		}

		// naive post
		embedfields := util.Chunk(StreamsToEmbedFields(matchingStreams...), 25)
		embeds := util.Map(embedfields, StreamsMessageEmbed)
		messagechunks := util.Chunk(embeds, 10)

		for idx, me := range messagechunks {
			newMessage, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
				Content: "Test Send",
				Embeds:  me,
			})
			if err != nil {
				log.Err(err).Msg("Failed to send message.")
			}

			m.Create(&models.Message{MessageID: newMessage.ID, SubscriptionID: sub.ID, PostOrder: idx})
		}
	}
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

func StreamMatchesSubscription(st *models.Stream, sub *models.Subscription) bool {
	return st.GameID == sub.GameID && strings.Contains(strings.ToLower(st.Title), strings.ToLower(sub.Filter))
}
