package commands

import (
	"fmt"
	"strconv"
	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/bwmarrin/discordgo"
)

var unsubscribeCmd = &Definition{
	Name: "unsubscribe",
	Base: &discordgo.ApplicationCommand{
		Description: "Remove a subscription to a game/filter pair.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		qs := query.Subscription
		allsubs, err := qs.Find()
		var content string
		if err != nil {
			content = err.Error()
		}
		if len(allsubs) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No active subscriptions.",
				},
			})
			return
		}

		selectedSub := get_option(s, i, "Which active subscription would you like to remove?",
			*util.Map(allsubs, func(sub *models.Subscription, _ int) discordgo.SelectMenuOption {
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
		content, err = (func() (string, error) {
			// Ignoring error as we generated these ourselves
			subid, _ := strconv.ParseUint(selectedSub, 10, 32)

			sub, err := qs.Where(qs.ID.Eq(uint(subid))).First()
			if err != nil {
				return "", err
			}

			_, err = qs.Delete(sub)
			if err != nil {
				return "", err
			}

			q := tickQuoteHelper
			out := fmt.Sprintf(`Unsubscribed from subscription--Game: %s`, q(sub.GameName))
			if sub.Filter != "" {
				out += q(sub.Filter)
			}

			return out, nil
		})()

		if err != nil {
			content = err.Error()
		}

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
	},
}

func init() {
	Register(unsubscribeCmd)
}
