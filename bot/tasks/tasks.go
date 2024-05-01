package tasks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Task struct {
	Spec string

	handler func(*discordgo.Session)
	dg      *discordgo.Session
}

func (t *Task) BindSession(s *discordgo.Session) {
	t.dg = s
	log.Print("input: ", s)
	log.Print("output: ", t.dg)
}

func (t *Task) Run() {
	t.handler(t.dg)
}

var All = []*Task{}
