package commands

import (
	"github.com/bwmarrin/discordgo"
)

type handler = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Definition struct {
	Base    *discordgo.ApplicationCommand
	Name    string
	Handler handler
}

var commands = []*discordgo.ApplicationCommand{}
var handlers = map[string]handler{}

func GetCommands() []*discordgo.ApplicationCommand {
	return commands
}

func Register(cmd *Definition) {
	cmd.Base.Name = cmd.Name
	commands = append(commands, cmd.Base)
	handlers[cmd.Name] = cmd.Handler
}

func SlashCommandRouter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func tickQuoteHelper(s string) string {
	return "`" + s + "`"
}
