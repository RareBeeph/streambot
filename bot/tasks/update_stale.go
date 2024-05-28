package tasks

import (
	"streambot/models"

	"github.com/bwmarrin/discordgo"
)

var updateStale = Task{
	Spec:         "*/30 * * * *",
	runOnStartup: false,
	handler: func(s *discordgo.Session) {
		updateMessages(s, models.SubHealths.Stale, models.SubHealths.Orphaned)
	},
}

func init() {
	All = append(All, &updateStale)
}
