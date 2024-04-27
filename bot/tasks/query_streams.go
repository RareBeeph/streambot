package tasks

import (
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"

	"github.com/jinzhu/copier"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

var queryStreams = Task{
	Spec: "*/5 * * * *",
	Handler: func() {
		qs := query.Subscription

		gameIDs := make([]string, 0)
		qs.Distinct(qs.GameID).Pluck(qs.GameID, gameIDs)

		resp, err := twitch.Client.GetStreams(&helix.StreamsParams{
			GameIDs: gameIDs,
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
