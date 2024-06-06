package commands

import (
	"fmt"
	"strconv"
	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var unsubscribeCmd = &Definition{
	Name: "unsubscribe",
	Base: &discordgo.ApplicationCommand{
		Description: "Remove a subscription to a game/filter pair.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		qs := query.Subscription

		allsubs, err := qs.Where(qs.ChannelID.Eq(i.ChannelID)).Find()
		if len(allsubs) == 0 || err != nil {
			msg := "No active subscriptions"
			if err != nil {
				msg = err.Error()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: msg,
				},
			})
			return
		}

		selectedSub, i := get_option(s, i, "Which active subscription would you like to remove?",
			util.Map(allsubs, func(sub *models.Subscription, _ int) discordgo.SelectMenuOption {
				label := sub.GameName + " (ID: " + sub.GameID + ")"
				if sub.Filter != "" {
					label += " with filter: " + sub.Filter
				}

				return discordgo.SelectMenuOption{
					Label: label,
					Value: fmt.Sprint(sub.ID),
				}
			}),
		)

		// This func just exists as a layer from which to only partially return on error
		msg, err := (func() (string, error) {
			// Ignoring error as we generated these ourselves
			subid, _ := strconv.ParseUint(selectedSub, 10, 32)

			// Fetch the subscription
			sub, err := qs.Preload(qs.Messages).Where(qs.ID.Eq(uint(subid))).First()
			if err != nil {
				return "", err
			}

			// Cleanup the posted messages from the channel
			// This would be a bulk delete, but that only works on
			// messages younger than 14 days
			for _, message := range sub.Messages {
				err := s.ChannelMessageDelete(sub.ChannelID, message.MessageID)
				if err != nil {
					log.Err(err).Msg("Failed to delete subscription message")
				}
			}

			// If that was successful, clean the objects out of our database
			_, err = qs.Select(qs.Messages.Field()).Delete(sub)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf(`Unsubscribed from subscription--%s`, sub), nil
		})()

		if err != nil {
			msg = err.Error()
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: msg,
			},
		})
	},
}

func init() {
	Register(unsubscribeCmd)
}
