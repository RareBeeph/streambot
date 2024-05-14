package commands

import (
	"fmt"
	"strconv"

	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/bwmarrin/discordgo"
)

var reactivateCmd = &Definition{
	Name: "reactivate",
	Base: &discordgo.ApplicationCommand{
		Description: "Reactivates a deactivated subscription.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// mostly copied from unsubscribe.go

		qs := query.Subscription

		deactivatedsubs, err := qs.Where(qs.TimesFailed.Gte(models.MaxTimesFailed)).Find()
		if len(deactivatedsubs) == 0 || err != nil {
			msg := "No deactivated subscriptions"
			if err != nil {
				msg = err.Error()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
			return
		}

		selectedSub, i := get_option(s, i, "Which subscription would you like to reactivate?",
			util.Map(deactivatedsubs, func(sub *models.Subscription, _ int) discordgo.SelectMenuOption {
				// maybe replace label with sub.String()
				label := sub.GameName + " (ID: " + sub.GameID + ")"
				if sub.Filter != "" {
					label += " with filter: " + sub.Filter
				}

				return discordgo.SelectMenuOption{
					Emoji: discordgo.ComponentEmoji{
						Name: "🦦", // temp emoji
					},
					Label: label,
					Value: fmt.Sprint(sub.ID), // maybe don't use db id here
				}
			}),
		)

		// This func just exists as a layer from which to only partially return on error
		msg, err := (func() (string, error) {
			// Ignoring error as we generated these ourselves
			subid, _ := strconv.ParseUint(selectedSub, 10, 32)

			sub, err := qs.Where(qs.ID.Eq(uint(subid))).First()
			if err != nil {
				return "", err
			}

			_, err = qs.Where(qs.ID.Eq(uint(subid))).Update(qs.TimesFailed, 0)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf(`Reactivated subscription--%s`, sub), nil
		})()

		if err != nil {
			msg = err.Error()
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	},
}

func init() {
	Register(reactivateCmd)
}
