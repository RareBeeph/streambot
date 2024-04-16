package commands

import (
	"fmt"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
)

var subscriptionsCmd = &Definition{
	Name: "subscriptions",
	Base: &discordgo.ApplicationCommand{
		Description: "Displays currently active subscriptions for the current Discord channel.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		qs := query.Subscription

		var content string
		subscriptions, err := qs.Find()
		if err != nil {
			content = err.Error()
		} else if len(subscriptions) > 0 {
			for i, sub := range subscriptions {
				q := tickQuoteHelper
				content += fmt.Sprintf(`%d)  Game Name: %s  |  GameID: %s`, i+1, q(sub.GameName), q(sub.GameID))
				if sub.Filter != "" {
					content += "  |  Filter: " + q(sub.Filter)
				}
				content += "\n"
			}
		} else {
			content = "No subscriptions currently active."
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	},
}

func init() {
	Register(subscriptionsCmd)
}
