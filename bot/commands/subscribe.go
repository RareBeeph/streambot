package commands

import (
	"log"
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix/v2"
)

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

		optionValues := make(map[string]string, len(options))
		for _, opt := range options {
			optionValues[opt.Name] = opt.StringValue()
		}

		gamesResponse, err := twitch.Client.GetGames(&helix.GamesParams{
			Names: []string{optionValues["game_name"]},
		})
		if err != nil {
			log.Println("I have no idea what to do here") // Temp
		} else if len(gamesResponse.Data.Games) == 0 {
			log.Println("No games found") // Temp
		} else if len(gamesResponse.Data.Games) > 1 {
			log.Println("More than one game with the same title???") // Temp
		} else {
			qs := query.Subscription

			game := gamesResponse.Data.Games[0]

			sub := &models.Subscription{
				GameName:  game.Name,
				GameID:    game.ID,
				Filter:    optionValues["optional_filter"],
				GuildID:   i.GuildID,
				ChannelID: i.ChannelID,
			}

			qs.Create(sub) // TODO: handle potential error
			qs.Delete(sub) // Temp because I don't have an unsubscribe command yet
		}

		// Temp; echo
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: optionValues["game_name"] + " " + optionValues["optional_filter"],
			},
		})
	},
}

func init() {
	Register(subscribeCmd)
}
