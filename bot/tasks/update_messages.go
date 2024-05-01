package tasks

import (
	"streambot/models"
	"streambot/query"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func updateMessages(s *discordgo.Session, streams []*models.Stream) {
	m := query.Message
	qs := query.Subscription

	messages, err := m.Find()
	if err != nil {
		log.Err(err).Msg("Failed to find messages")
	}

	subscriptions, err := qs.Find()
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	for _, sub := range subscriptions {
		for _, st := range streams {
			if st.GameID == sub.GameID && strings.Contains(strings.ToLower(st.Title), strings.ToLower(sub.Filter)) {
				hasMessage := false

				for _, me := range sub.Messages {
					if me.UserID == st.UserID {
						hasMessage = true
						s.ChannelMessageEdit(sub.ChannelID, me.MessageID, "TestEdit: "+st.Title)
					}
				}

				if !hasMessage {
					log.Print("session: ", s)
					newMessage, _ := s.ChannelMessageSend(sub.ChannelID, "TestSend: "+st.Title)
					m.Create(&models.Message{MessageID: newMessage.ID, UserID: st.UserID})
				}
			}
		}
	}

	for _, me := range messages {
		hasStream := false
		for _, st := range streams {
			if me.UserID == st.UserID {
				hasStream = true
			}
		}

		if !hasStream {
			sub, _ := qs.Where(qs.ID.Eq(me.SubscriptionID)).First()

			s.ChannelMessageDelete(sub.ChannelID, me.MessageID)
			m.Delete(me)
		}
	}
}
