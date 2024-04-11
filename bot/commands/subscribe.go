package commands

import "github.com/bwmarrin/discordgo"

func init() {
	Register(&Definition{
		Name: "subscribe",
		Base: &discordgo.ApplicationCommand{
			Description: "hello world",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "game_name",
					Description: "Game name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "optional_filter",
					Description: "Only subscribe to streams containing the filter string in their titles",
					Required:    false,
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			// very much not permanent handling behavior
			filter := " "
			if len(options) == 2 {
				filter += options[1].StringValue()
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: options[0].StringValue() + filter,
				},
			})
		},
	})
}
