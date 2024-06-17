/*
Copyright © 2024 Aria Taylor <ari@aricodes.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"streambot/config"
	"streambot/query"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Flushes our past 100 messages from known channels",
	Long: `Flushes our past from known channels.

Intended for use during data migration from an older iteration of the bot.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := query.Subscription
		flushdb, _ := cmd.Flags().GetBool("flushdb")
		searchMax, _ := cmd.Flags().GetInt("count")
		conf := config.Load()

		discord, err := discordgo.New("Bot " + conf.Token)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to discord!")
		}
		discord.Open()

		channels := []string{}
		s.Pluck(s.ChannelID.Distinct(), &channels)

		for _, id := range channels {
			before := ""
			searched := 0

			for searched <= searchMax {
				messages, _ := discord.ChannelMessages(id, 100, before, "", "")

				for _, msg := range messages {
					if msg.Author.ID == discord.State.User.ID {
						discord.ChannelMessageDelete(id, msg.ID)
					}

					before = msg.ID
					searched += 1

					if searched%100 == 0 {
						log.Info().Msgf("Searched through %v messages...", searched)
					}
				}
			}

			if flushdb {
				s.Where(s.ChannelID.Eq(id)).Delete()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolP("flushdb", "f", false, "Clear the database at the same time")
	cleanCmd.Flags().IntP("count", "n", 200, "How many messages to search back through in each channel")
}
