package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/teris-io/shortid"
)

func get_option(s *discordgo.Session, i *discordgo.InteractionCreate, prompt string, options []discordgo.SelectMenuOption) string {
	selectId := shortid.MustGenerate()
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: prompt,
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID: selectId,
					Options:  options,
				},
			},
		},
	}

	// TODO: Create channel

	unregisterHandler := s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		if id := i.MessageComponentData().CustomID; id == selectId {
			// TODO: send this through a channel
		}
	})
	defer unregisterHandler()

	s.InteractionRespond(i.Interaction, response)

	// TODO: Wait on channel, return value from it
	return ""
}
