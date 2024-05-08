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

func updateMessages(s *discordgo.Session) {
	type msgAction struct {
		target  *models.Message
		content []*discordgo.MessageEmbed
	}

	m := query.Message
	qs := query.Subscription
	qst := query.Stream

	subscriptions, err := qs.Preload(qs.Messages).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	for _, sub := range subscriptions {
		matchingStreams, err := qst.Where(qst.GameID.Eq(sub.GameID), qst.Title.Lower().Like(fmt.Sprintf("%%%s%%", sub.Filter))).Find()
		if err != nil {
			log.Err(err).Msg("Failed to find matching streams.")
		}

		// The modifications we're going to make to our messages
		actions := []*msgAction{}

		embedFields := util.Chunk(StreamsToEmbedFields(matchingStreams...), 25)
		embeds := util.Map(embedFields, StreamsMessageEmbed)
		messageChunks := util.Chunk(embeds, 10)

		// we could manully sort sub.Messages instead
		sortedMessages, err := m.Where(m.SubscriptionID.Eq(sub.ID)).Order(m.PostOrder.Asc()).Find()
		if err != nil {
			log.Err(err).Msg("Failed to find or sort messages.")
		}
		if len(sortedMessages) != len(sub.Messages) {
			log.Err(err).Msg("DEBUG: didn't find the expected number of messages")
		}

		// messageCount := len(sub.Messages)
		messageCount := len(sortedMessages)
		for idx, embed := range messageChunks {
			// Determine what action needs to be taken to post this chunk
			if idx < messageCount {
				actions = append(actions, &msgAction{target: sortedMessages[idx], content: embed})
			} else {
				actions = append(actions, &msgAction{content: embed})
			}
		}

		// we're gonna handle bulk deletion later

		// if messageCount > len(messageChunks) {
		// 	// Issue a delete action
		// 	for idx, mess := range sortedMessages[messageCount:] {
		// 		actions = append(actions, &msgAction{target: mess})
		// 	}
		// }

		for idx, action := range actions {
			// Perform post/edit and update database
			if action.target == nil {
				// no target => post
				message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
					Content: "placeholder",
					Embeds:  action.content,
				})
				if err != nil {
					log.Err(err).Msg("Failed to send message.")
				}

				m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: idx})
			} else {
				// yes target => edit
				textContent := "placeholder"
				_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					Content: &textContent,
					Embeds:  action.content,

					ID:      action.target.MessageID,
					Channel: sub.ChannelID,
				})
				if err != nil {
					log.Err(err).Msg("Failed to edit message.")
				}
			}
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
