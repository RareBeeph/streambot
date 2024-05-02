package tasks

import (
	"github.com/bwmarrin/discordgo"
)

type Task struct {
	Spec string

	handler func(*discordgo.Session)
	dg      *discordgo.Session
}

func (t *Task) BindSession(s *discordgo.Session) {
	t.dg = s
}

func (t *Task) Run() {
	t.handler(t.dg)
}

var All = []*Task{}
