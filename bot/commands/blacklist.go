package commands

import (
	"fmt"
	"streambot/bot/twitch"
	"streambot/models"
	"streambot/query"
	"streambot/util"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix/v2"
	"gorm.io/gorm/clause"
)

var blacklistCmd = &Definition{
	Name: "blacklist",
	Base: &discordgo.ApplicationCommand{
		Description: "Prevents a Twitch user's streams from appearing in the current channel.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user_name",
				Description: "User login name",
				Required:    true,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		input := i.ApplicationCommandData().Options[0].StringValue()

		usersResponse, err := twitch.Client.GetUsers(&helix.UsersParams{
			Logins: []string{strings.ToLower(input)},
		})

		// copying my content/error handling structure from subscribe. still probably not optimal but we can fix it later
		var content string
		var userID string
		if err != nil {
			content = err.Error()
		} else {
			users := usersResponse.Data.Users

			if len(users) == 0 {
				content = "No matching users found."
			} else if len(users) == 1 {
				userID = users[0].ID
			} else {
				userID, i = get_option(s, i,
					"Which of these users did you mean?",
					util.Map(users, func(u helix.User, _ int) discordgo.SelectMenuOption {
						return discordgo.SelectMenuOption{
							Label: u.Login + " (ID: " + u.ID + ")",
							Value: u.ID,
						}
					}),
				)
			}
		}

		qb := query.BlacklistEntry

		blk := &models.BlacklistEntry{
			UserID:    userID,
			UserLogin: input, // possibly sketchy but this is just for easy display later
			ChannelID: i.ChannelID,
		}

		err = qb.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "channel_id"}},
			UpdateAll: true,
		}).Create(blk)
		if err != nil {
			content = err.Error()
		} else {
			q := tickQuoteHelper

			content = fmt.Sprintf(`Blacklist entry added for user: %s (ID: %s)`, q(strings.ToLower(input)), q(blk.UserID))
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
	Register(blacklistCmd)
}
