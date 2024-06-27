package tasks

import (
	"streambot/bot/tasks/update"
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/copier"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"

	mapset "github.com/deckarep/golang-set/v2"

	"gorm.io/gorm/clause"

	"sync"
)

var queryStreamsAndUpdateHealthy = Task{
	Spec:         "*/5 * * * *",
	runOnStartup: true,
	handler: func(s *discordgo.Session) {
		var wg sync.WaitGroup

		qs := query.Subscription

		// Get every distinct pair of (Language, GameID)
		gamesByLanguage, err := qs.Distinct(qs.Language, qs.GameID).Select(qs.ALL).Find()
		if err != nil {
			log.Err(err).Msg("Failed to get language/gameID pairs.")
			return
		}

		// Store a deduplicated list of all game IDs we have fetched streams for
		gameIDs := mapset.NewSet[string]()

		// Query 100 streams for each of them in parallel
		wg.Add(len(gamesByLanguage))
		// We use a map instead of a slice as our intermediary
		// results store for thread safety
		queryResults := sync.Map{}
		for _, sub := range gamesByLanguage {
			go func(sub *models.Subscription) {
				defer wg.Done()
				defer gameIDs.Add(sub.GameID)
				streams, _ := twitch.FetchStreams(sub.GameID, sub.Language)
				queryResults.Store(sub, streams)
			}(sub)
		}
		wg.Wait()

		// Flatten that map into a slice
		rawStreams := []helix.Stream{}
		queryResults.Range(func(_, streams any) bool {
			rawStreams = append(rawStreams, streams.([]helix.Stream)...)
			return true
		})

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

			_, err = qst.Where(qst.GameID.In(gameIDs.ToSlice()...)).Delete()
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

		update.UpdateSubscriptions(s, models.SubHealths.Healthy, models.SubHealths.Stale)
	},
}

func init() {
	All = append(All, &queryStreamsAndUpdateHealthy)
}
