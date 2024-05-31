package tasks

import (
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/copier"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

var queryStreamsAndUpdateHealthy = Task{
	Spec:         "*/5 * * * *",
	runOnStartup: true,
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
			First:   100,
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

		_, err = qst.Where(qst.GameID.In(gameIDs...)).Delete()
		if err != nil {
			log.Err(err).Msg("Failed to properly clear stale streams.")
			return
		}

		err = qst.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).Create(streams...)
		if err != nil {
			log.Err(err).Msg("Failed to update stream list.")
		}

		updateMessages(s, models.SubHealths.Healthy, models.SubHealths.Stale)
	},
}

func init() {
	All = append(All, &queryStreamsAndUpdateHealthy)
}
