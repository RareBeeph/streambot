package commands

import "github.com/bwmarrin/discordgo"

var subscribeCmd = &Definition{
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

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: optionMap["game_name"].StringValue() + " " + optionMap["optional_filter"].StringValue(),
			},
		})
	},
}

func init() {
	Register(subscribeCmd)
}
