package commands

import (
	"fmt"
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix/v2"
	"gorm.io/gorm/clause"
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

		var content string
		var game helix.Game
		if err != nil {
			content = err.Error()
		} else if len(gamesResponse.Data.Games) == 0 {
			content = "No matching games found."
		} else if len(gamesResponse.Data.Games) > 1 {
			var selectedGameID string
			selectedGameID, i = get_option(s, i,
				"Which of these games did you mean?",
				util.Map(gamesResponse.Data.Games, func(game helix.Game, _ int) discordgo.SelectMenuOption {
					return discordgo.SelectMenuOption{
						Label: game.Name + " (ID: " + game.ID + ")",
						Value: game.ID,
					}
				}),
			)
			for _, g := range gamesResponse.Data.Games {
				if g.ID == selectedGameID {
					game = g
				}
			}
		} else {
			game = gamesResponse.Data.Games[0]
		}

		if (game != helix.Game{}) {
			qs := query.Subscription

			sub := &models.Subscription{
				GameName:  game.Name,
				GameID:    game.ID,
				Filter:    optionValues["optional_filter"],
				GuildID:   i.GuildID,
				ChannelID: i.ChannelID,
			}

			err = qs.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "game_id"}, {Name: "filter"}, {Name: "channel_id"}},
				UpdateAll: true, // Specifically, reset TimesFailed to 0 on failure of this index
			}).Create(sub)
			if err != nil {
				content = err.Error()
			} else {
				q := tickQuoteHelper

				content = fmt.Sprintf(`Subscription added for game: %s (ID: %s)`, q(sub.GameName), q(sub.GameID))
				if sub.Filter != "" {
					content += " with filter: " + q(sub.Filter)
				}
			}
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: content,
			},
		})
	},
}

func init() {
	Register(subscribeCmd)
}
