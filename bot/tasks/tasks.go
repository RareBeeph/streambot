package tasks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Task struct {
	Spec string

	handler      func(*discordgo.Session)
	runOnStartup bool
	dg           *discordgo.Session
}

func (t *Task) BindSession(s *discordgo.Session) {
	t.dg = s
}

func (t *Task) Run() {
	if t.dg == nil {
		log.Error().Msg("Cannot run tasks before binding a session to them.")
	}

	t.handler(t.dg)
}

var All = []*Task{}

func Startup() {
	for _, a := range All {
		if a.runOnStartup {
			a.Run()
		}
	}
}
