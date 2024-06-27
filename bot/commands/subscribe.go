package commands

import (
	"fmt"
	"streambot/bot/tasks/update"
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"strings"

	"github.com/bwmarrin/discordgo"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
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
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "optional_language",
				Description:  "Only subscribe to streams in the specified language",
				Required:     false,
				Autocomplete: true,
			},
		},
	},
	handler:      subscribeHandler,
	autocomplete: subscribeAutocompleteHandler,
}

func subscribeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// this is a mess
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

	var sub *models.Subscription
	if (game != helix.Game{}) {
		qs := query.Subscription

		sub = &models.Subscription{
			GameName:  game.Name,
			GameID:    game.ID,
			Filter:    optionValues["optional_filter"],
			GuildID:   i.GuildID,
			ChannelID: i.ChannelID,
		}

		if optionValues["optional_language"] != "" {
			sub.Language = optionValues["optional_language"]
		}

		err = qs.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "game_id"}, {Name: "filter"}, {Name: "channel_id"}, {Name: "language"}},
			UpdateAll: true, // Specifically, reset TimesFailed to 0 on failure of this index
		}).Create(sub)
		if err != nil {
			content = err.Error()
		} else {
			content = fmt.Sprintf(`Subscription added--%s`, sub)
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})

	if sub != nil {
		update.UpdateSubscription(s, sub)
	}
}

func subscribeAutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: generate dynamically according to user input from ISO 639-1 specs,
	// with language names written in their language
	var input string
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "optional_language" {
			input = strings.ToLower(opt.StringValue())
		}
	}

	// should we prioritize languages with more user representation on discord?
	matchingCodes := util.Filter(iso6391.Codes, func(c string, idx int) bool {
		return strings.HasPrefix(c, input) ||
			strings.HasPrefix(strings.ToLower(iso6391.Name(c)), input) ||
			strings.HasPrefix(strings.ToLower(iso6391.NativeName(c)), input)
	})
	choices := util.Map(matchingCodes, func(c string, idx int) *discordgo.ApplicationCommandOptionChoice {
		return &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("%s (%s)", c, iso6391.NativeName(c)),
			Value: c,
		}
	})
	if len(choices) > 25 {
		choices = choices[:25]
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		log.Err(err).Msg("Failed to respond with autocomplete choices")
	}
}

func init() {
	Register(subscribeCmd)
}
