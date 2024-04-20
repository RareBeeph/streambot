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
		func() {
			if selectedSub == "" {
				content = "Error: Somehow, the selected subscription had no ID."
				return
			}

			subid, err := strconv.ParseUint(selectedSub, 10, 32)
			if err != nil {
				content = err.Error()
				return
			}

			temp, err := qs.Where(qs.ID.Eq(uint(subid))).First()
			if err != nil {
				content = err.Error()
				return
			}

			_, err = qs.Delete(temp)
			if err != nil {
				content = err.Error()
				return
			}

			q := tickQuoteHelper
			content = fmt.Sprintf(`Unsubscribed from subscription--Game: %s`, q(temp.GameName))
			if temp.Filter != "" {
				content += q(temp.Filter)
			}
		}()

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content:    &content,
			Components: &[]discordgo.MessageComponent{},
		})
	},
}

func init() {
	Register(unsubscribeCmd)
}
