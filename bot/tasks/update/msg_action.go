package update

import (
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type msgAction struct {
	target  *models.Message
	content []*discordgo.MessageEmbed
}

func (a *msgAction) perform(s *discordgo.Session, sub *models.Subscription, postIdx int) (err error) {
	// Perform post/edit/delete and update database
	if a.target == nil {
		err = a.postContent(s, sub, postIdx)
	} else if a.content == nil {
		err = a.deleteTarget(s, sub)
	} else {
		// if action has a target and content, edit that target with that content
		err = a.editMessage(s, sub)
	}
	return
}

func (a *msgAction) postContent(s *discordgo.Session, sub *models.Subscription, postIdx int) (err error) {
	m := query.Message
	qs := query.Subscription

	message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
		Embeds: a.content,
	})

	// if the target channel or server no longer exists, remove this subscription
	if channelNoLongerValid(err) {
		qs.Select(qs.Messages.Field()).Delete(sub)
	}

	if err != nil {
		log.Err(err).Msg("Failed to send message.")
		return
	}

	m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: postIdx})
	return
}

func (a *msgAction) deleteTarget(s *discordgo.Session, sub *models.Subscription) (err error) {
	m := query.Message
	qs := query.Subscription

	err = s.ChannelMessageDelete(sub.ChannelID, a.target.MessageID)

	// if the target channel or server no longer exists, remove this subscription
	if channelNoLongerValid(err) {
		qs.Select(qs.Messages.Field()).Delete(sub)
	}

	if err != nil && !messageUnavailable(err) {
		log.Err(err).Msg("Failed to delete message.")
		return
	}

	// if we successfully deleted the message, or if it had already been deleted, remove our record of it
	m.Where(m.MessageID.Eq(a.target.MessageID)).Delete()
	return
}

func (a *msgAction) editMessage(s *discordgo.Session, sub *models.Subscription) (err error) {
	m := query.Message
	qs := query.Subscription

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds: &a.content,

		ID:      a.target.MessageID,
		Channel: sub.ChannelID,
	})

	// if the target channel or server no longer exists, remove this subscription
	if channelNoLongerValid(err) {
		qs.Select(qs.Messages.Field()).Delete(sub)
	}

	// if the target message no longer exists, remove our record of it
	if messageUnavailable(err) {
		m.Where(m.MessageID.Eq(a.target.MessageID)).Delete()
	}

	// not an else if like in the delete case, because the intended edit never makes it to the user
	if err != nil {
		log.Err(err).Msg("Failed to edit message.")
	}
	return
}

func messageUnavailable(err error) bool {
	codes := []int{
		discordgo.ErrCodeUnknownMessage,
		// discordgo.ErrCodeUnknownChannel,
		// discordgo.ErrCodeUnknownGuild,
	}
	return isRestError(err, codes...)
}

func channelNoLongerValid(err error) bool {
	codes := []int{
		discordgo.ErrCodeUnknownChannel,
		discordgo.ErrCodeUnknownGuild,
	}
	return isRestError(err, codes...)
}
