package tasks

import (
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"gorm.io/gen"
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

		_, err = qst.Unscoped().Not(gen.Exists(
			qs.Limit(1).Where(
				qs.GameID.EqCol(qst.GameID),
				qs.Where(qs.Language.EqCol(qst.Language)).Or(qs.Language.Eq("")),
			),
		)).Delete()
		if err != nil {
			log.Err(err).Msg("Failed to delete orphaned stream records")
		}
	},
}

func init() {
	All = append(All, cleanupStreams)
}
