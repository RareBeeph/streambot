package update

import (
	"errors"
	"fmt"
	"slices"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func msgActionsForSub(sub *models.Subscription) ([]*msgAction, error) {
	// The actual content of the messages we intend to post
	embeds, err := embedsForSub(sub)
	if err != nil {
		return nil, err
	}
	messageChunks := util.Chunk(embeds, 2)

	// Temp: limit to 1 message per subscription for aesthetic purposes
	// TODO: reconsider this decision.
	if len(messageChunks) > 1 {
		messageChunks = messageChunks[:1]
	}

	// Determine what action needs to be taken to post each chunk
	actions := []*msgAction{}
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

	return actions, nil
}

func embedsForSub(sub *models.Subscription) ([]*discordgo.MessageEmbed, error) {
	streams, err := streamsForSub(sub)
	if err != nil {
		return nil, err
	}

	embedFields := util.Chunk(streamsToEmbedFields(streams...), 25)
	embeds := util.Map(embedFields, func(fields []*discordgo.MessageEmbedField, idx int) *discordgo.MessageEmbed {
		return &discordgo.MessageEmbed{
			Color:     0x9922cc,
			Author:    &discordgo.MessageEmbedAuthor{},
			Title:     fmt.Sprintf("Streams for %s", sub),
			Fields:    fields,
			Timestamp: time.Now().Format(time.RFC3339),
		}
	})

	return embeds, nil
}

func streamsForSub(sub *models.Subscription) ([]*models.Stream, error) {
	qb := query.BlacklistEntry
	qst := query.Stream

	// Determine which user IDs are blacklisted here
	matchingBlacklists, err := qb.Where(qb.ChannelID.Eq(sub.ChannelID)).Find()
	if err != nil {
		log.Err(err).Msg("Failed to find matching blacklists.")
		return []*models.Stream{}, err
	}
	blacklistUserIDs := util.Map(matchingBlacklists, func(bl *models.BlacklistEntry, idx int) string {
		return bl.UserID
	})

	// Query streams from our database, respecting the blacklist
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

func isRestError(err error, codes ...int) bool {
	var resterr *discordgo.RESTError
	return errors.As(err, &resterr) && slices.Contains(codes, resterr.Message.Code)
}
