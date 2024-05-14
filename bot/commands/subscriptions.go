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
				// Note: no longer displays game ID.
				// This might be unideal for users who wish to confirm externally that they queried the right game.
				content += fmt.Sprintf(`%d) %s`, i+1, sub)
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
