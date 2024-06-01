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

	"sync"
)

var queryStreamsAndUpdateHealthy = Task{
	Spec:         "*/5 * * * *",
	runOnStartup: true,
	handler: func(s *discordgo.Session) {
		var wg sync.WaitGroup

		qs := query.Subscription

		// Get a full list of GameIDs
		gameIDs := make([]string, 0)
		err := qs.Distinct(qs.GameID).Pluck(qs.GameID, &gameIDs)
		if err != nil {
			log.Err(err).Msg("Failed to pluck game IDs.")
			return
		}

		// Query 100 streams for each of them in parallel
		wg.Add(len(gameIDs))
		// We use a map instead of a slice as our intermediary
		// results store for thread safety
		queryResults := make(map[string][]helix.Stream)
		for _, gameID := range gameIDs {
			go func(gameID string) {
				defer wg.Done()
				streams, _ := twitch.FetchStreams(gameID)
				queryResults[gameID] = streams
			}(gameID)
		}
		wg.Wait()

		// Flatten that map into a slice
		rawStreams := []helix.Stream{}
		for _, streams := range queryResults {
			rawStreams = append(rawStreams, streams...)
		}

		// Convert from Helix streams to model instances
		streams := make([]*models.Stream, len(rawStreams))
		err = copier.Copy(&streams, &rawStreams)
		if err != nil {
			log.Err(err).Msg("Failed to copy streams.")
			return
		}

		// Delete old streams and create new ones in one atomic transaction
		err = query.Q.Transaction(func(tx *query.Query) error {
			qst := tx.Stream

			_, err = qst.Where(qst.GameID.In(gameIDs...)).Delete()
			if err != nil {
				return err
			}

			err = qst.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}},
				UpdateAll: true,
			}).Create(streams...)

			return err
		})
		if err != nil {
			log.Err(err).Msg("Failed to update stream list")
		}

		updateMessages(s, models.SubHealths.Healthy, models.SubHealths.Stale)
	},
}

func init() {
	All = append(All, &queryStreamsAndUpdateHealthy)
}
