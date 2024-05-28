package tasks

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Task struct {
	Spec string

	handler      func(*discordgo.Session)
	runOnStartup bool
	dg           *discordgo.Session
	mu           sync.Mutex
}

func (t *Task) BindSession(s *discordgo.Session) {
	t.dg = s
}

func (t *Task) Run() {
	if t.dg == nil {
		log.Error().Msg("Cannot run tasks before binding a session to them.")
		return
	}

	if locked := t.mu.TryLock(); !locked {
		// Maybe we should add a Name field so that we can improve this output?
		log.Error().Msg("Skipping task execution, prior run still going")
		return
	}
	defer t.mu.Unlock()

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
