package tasks

import (
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/jinzhu/copier"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

var queryStreams = Task{
	Spec: "*/5 * * * *",
	Handler: func() {
		subscriptions, _ := query.Subscription.Find() // err unhandled
		if len(subscriptions) == 0 {
			return
		}

		GameIDs := *util.Map(subscriptions, func(s *models.Subscription, i int) string {
			return s.GameID
		})

		resp, err := twitch.Client.GetStreams(&helix.StreamsParams{
			GameIDs: GameIDs,
		})
		if err != nil {
			log.Print(err)
			return // idk what to do here yet
		}

		for _, s := range resp.Data.Streams {
			stream := &models.Stream{}
			copier.Copy(stream, s)      // err unhandled
			query.Stream.Create(stream) // err unhandled
			log.Print(stream)           // making sure i know what's going on
			query.Stream.Delete(stream) // temp cleanup
		}
	},
}

func init() {
	All = append(All, queryStreams)
}
