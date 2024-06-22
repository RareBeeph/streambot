package messages

import (
	"errors"
	"slices"
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type msgAction struct {
	target  *models.Message
	content []*discordgo.MessageEmbed
}

func (a *msgAction) postContent(s *discordgo.Session, sub *models.Subscription, idx int) (errored bool) {
	m := query.Message
	qs := query.Subscription

	message, err := s.ChannelMessageSendComplex(sub.ChannelID, &discordgo.MessageSend{
		Embeds: a.content,
	})

	if channelNoLongerValid(err) {
		// if the target channel or server no longer exists, remove this subscription
		qs.Select(qs.Messages.Field()).Delete(sub)
		log.Err(err).Msg("Failed to send message.")
		return
	}

	if err != nil {
		log.Err(err).Msg("Failed to send message.")
		errored = true
		return
	}

	m.Create(&models.Message{MessageID: message.ID, SubscriptionID: sub.ID, PostOrder: idx})
	return
}

func (a *msgAction) deleteTarget(s *discordgo.Session, sub *models.Subscription) (errored bool) {
	m := query.Message
	qs := query.Subscription

	err := s.ChannelMessageDelete(sub.ChannelID, a.target.MessageID)

	if channelNoLongerValid(err) {
		// if the target channel or server no longer exists, remove this subscription
		qs.Select(qs.Messages.Field()).Delete(sub)
		return
	}

	if err != nil && !messageUnavailable(err) {
		log.Err(err).Msg("Failed to delete message.")
		errored = true
		return
	}

	// If we successfully deleted the message, or if it had already been deleted, remove our record of it
	m.Where(m.MessageID.Eq(a.target.MessageID)).Delete()
	return
}

func (a *msgAction) editMessage(s *discordgo.Session, sub *models.Subscription) (errored bool) {
	m := query.Message
	qs := query.Subscription

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds: &a.content,

		ID:      a.target.MessageID,
		Channel: sub.ChannelID,
	})

	if channelNoLongerValid(err) {
		// if the target channel or server no longer exists, remove this subscription
		qs.Select(qs.Messages.Field()).Delete(sub)
		log.Err(err).Msg("Failed to edit message.")
		return
	}

	if messageUnavailable(err) {
		// if the target message no longer exists, remove our record of it
		m.Where(m.MessageID.Eq(a.target.MessageID)).Delete()
	}

	// not an else if like in the delete case, because the intended edit never makes it to the user
	if err != nil {
		errored = true
		log.Err(err).Msg("Failed to edit message.")
	}
	return
}

func messageUnavailable(err error) bool {
	var resterr *discordgo.RESTError
	var errorcodes = []int{
		discordgo.ErrCodeUnknownMessage,
		// discordgo.ErrCodeUnknownChannel,
		// discordgo.ErrCodeUnknownGuild,
	}
	return errors.As(err, &resterr) && slices.Contains(errorcodes, resterr.Message.Code)
}

func channelNoLongerValid(err error) bool {
	var resterr *discordgo.RESTError
	var errorcodes = []int{
		discordgo.ErrCodeUnknownChannel,
		discordgo.ErrCodeUnknownGuild,
	}
	return errors.As(err, &resterr) && slices.Contains(errorcodes, resterr.Message.Code)
}
