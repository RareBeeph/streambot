package commands

import "github.com/bwmarrin/discordgo"

var helloWorldCmd = &Definition{
	Name: "helloworld",
	Base: &discordgo.ApplicationCommand{
		Description: "hello world",
	},
	handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "hello world",
			},
		})
	},
}

func init() {
	Register(helloWorldCmd)
}
