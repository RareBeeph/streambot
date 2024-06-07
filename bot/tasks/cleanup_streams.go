package tasks

import (
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var cleanupStreams = &Task{
	Spec:         "*/30 * * * *",
	runOnStartup: true, // why not
	handler: func(s *discordgo.Session) {
		qst := query.Stream
		qs := query.Subscription

		gameids := make([]string, 0)
		err := qs.Pluck(qs.GameID, &gameids)
		if err != nil {
			log.Err(err).Msg("Failed to pluck game IDs")
		}

		// note: does not consider streams of a subscribed game but in the wrong language to be orphaned
		// (mostly because i don't know how to phrase a condition for "these two columns both don't match any model instance in this array")
		_, err = qst.Unscoped().Where(qst.GameID.NotIn(gameids...)).Delete()
		if err != nil {
			log.Err(err).Msg("Failed to delete orphaned stream records")
		}
	},
}

func init() {
	All = append(All, cleanupStreams)
}
