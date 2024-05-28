package tasks

import (
	"streambot/models"

	"github.com/bwmarrin/discordgo"
)

var updateStale = Task{
	Spec:         "*/30 * * * *",
	runOnStartup: false,
	handler: func(s *discordgo.Session) {
		updateMessages(s, models.SubHealths.Stale, models.SubHealths.Orphaned) // maybe rename that const or add a new one for the real max
	},
}

func init() {
	All = append(All, &updateStale)
}
