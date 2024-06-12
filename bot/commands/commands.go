package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type handler = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Definition struct {
	Base         *discordgo.ApplicationCommand
	Name         string
	handler      handler
	autocomplete handler
}

func (d *Definition) Interact(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// d.handler should never be nil but it can't hurt to check
		if d.handler != nil {
			d.handler(s, i)
		} else {
			log.Error().Msg("Definitions must have a non-nil primary handler")
		}
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		if d.autocomplete != nil {
			d.autocomplete(s, i)
		}
		return
	default:
	}
}

var commands = []*discordgo.ApplicationCommand{}
var handlers = map[string]*Definition{}

func All() []*discordgo.ApplicationCommand {
	return commands
}

func Register(cmd *Definition) {
	cmd.Base.Name = cmd.Name
	commands = append(commands, cmd.Base)
	handlers[cmd.Name] = cmd
}

func SlashCommandRouter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// nearly redundant type check, but we need to know we can legally call ApplicationCommandData()
	if i.Type != discordgo.InteractionApplicationCommand && i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	if d, ok := handlers[i.ApplicationCommandData().Name]; ok {
		d.Interact(s, i)
	}
}

func tickQuoteHelper(s string) string {
	return "`" + s + "`"
}
