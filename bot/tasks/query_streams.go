package tasks

import (
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/copier"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

var queryStreams = Task{
	Spec: "*/5 * * * *",
	handler: func(s *discordgo.Session) {
		qs := query.Subscription
		qst := query.Stream

		gameIDs := make([]string, 0)
		err := qs.Distinct(qs.GameID).Pluck(qs.GameID, &gameIDs)
		if err != nil {
			log.Err(err).Msg("Failed to pluck game IDs.")
			return
		}

		resp, err := twitch.Client.GetStreams(&helix.StreamsParams{
			GameIDs: gameIDs,
		})
		if err != nil {
			log.Err(err).Msg("Failed to fetch streams")
			return
		}

		streams := make([]*models.Stream, len(resp.Data.Streams))
		err = copier.Copy(&streams, &resp.Data.Streams)
		if err != nil {
			log.Err(err).Msg("Failed to copy streams.")
			return
		}

		_, err = qst.Unscoped().Where(qst.GameID.In(gameIDs...)).Delete()
		if err != nil {
			log.Err(err).Msg("Failed to properly clear stale streams.")
			return
		}

		err = qst.Create(streams...)
		if err != nil {
			log.Err(err).Msg("Failed to update stream list.")
		}

		// updateMessages(s, streams)
	},
}

func init() {
	All = append(All, &queryStreams)
}
