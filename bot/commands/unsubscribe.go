package commands

import (
	"fmt"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
)

var unsubscribeCmd = &Definition{
	Name: "unsubscribe",
	Base: &discordgo.ApplicationCommand{
		Description: "Unsubscribe from a command by index (use /subscriptions to check indices)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "sub_idx",
				Description: "The index of the subscription to unsubscribe from",
				Required:    true,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		qs := query.Subscription
		subIdx := i.ApplicationCommandData().Options[0].IntValue()

		var content string
		subscriptions, err := qs.Find()
		if err != nil {
			content = err.Error()
		} else if subIdx < 0 {
			content = "Sub_idx must be greater than 0."
		} else if subIdx > int64(len(subscriptions)) {
			content = "Sub_idx too large."
		} else {
			temp := *subscriptions[subIdx-1]
			_, err = qs.Delete(subscriptions[subIdx-1])
			if err != nil {
				content = err.Error()
			} else {
				q := tickQuoteHelper
				content = fmt.Sprintf(`Unsubscribed from subscription %s (Game: %s`, q(fmt.Sprint(subIdx)), q(temp.GameName))
				if temp.Filter != "" {
					content += q(temp.Filter)
				}
				content += ")"
			}
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
	Register(unsubscribeCmd)
}
