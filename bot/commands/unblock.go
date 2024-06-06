package commands

import (
	"fmt"
	"strconv"
	"streambot/models"
	"streambot/query"
	"streambot/util"

	"github.com/bwmarrin/discordgo"
)

// mostly copied from unsubscribe
var unblockCmd = &Definition{
	Name: "unblock",
	Base: &discordgo.ApplicationCommand{
		Description: "Remove a user from this channe's blacklist.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		qb := query.BlacklistEntry

		allblocks, err := qb.Where(qb.ChannelID.Eq(i.ChannelID)).Find()
		if len(allblocks) == 0 || err != nil {
			msg := "Blacklist is empty."
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

		selectedBlock, i := get_option(s, i, "Which blacklist entry would you like to remove?",
			util.Map(allblocks, func(block *models.BlacklistEntry, _ int) discordgo.SelectMenuOption {
				label := fmt.Sprintf("%s (ID: %s)", block.UserLogin, block.UserID)

				return discordgo.SelectMenuOption{
					Label: label,
					Value: fmt.Sprint(block.ID),
				}
			}),
		)

		// This func just exists as a layer from which to only partially return on error
		msg, err := (func() (string, error) {
			// Ignoring error as we generated these ourselves
			blockid, _ := strconv.ParseUint(selectedBlock, 10, 32)

			// Get the blacklist entry (for info display purposes)
			block, err := qb.Where(qb.ChannelID.Eq(i.ChannelID), qb.ID.Eq(uint(blockid))).First()
			if err != nil {
				return "", err
			}

			// Delete the blacklist entry
			_, err = qb.Delete(block)
			if err != nil {
				return "", err
			}

			q := tickQuoteHelper
			return fmt.Sprintf(`Unblocked user %s (ID: %s)`, q(block.UserLogin), q(block.UserID)), nil
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
	Register(unblockCmd)
}
