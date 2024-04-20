package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

func get_option(s *discordgo.Session, i *discordgo.InteractionCreate, prompt string, options []discordgo.SelectMenuOption) string {
	selectId := shortid.MustGenerate()
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: prompt,
			// Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID: selectId,
							Options:  options,
						},
					},
				},
			},
		},
	}

	c := make(chan string)

	unregisterHandler := s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		if id := i.MessageComponentData().CustomID; id == selectId {
			c <- i.MessageComponentData().Values[0]
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
		}
	})
	defer unregisterHandler()

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		// temp
		log.Print(err)
	} else {
		v := <-c
		return v
	}

	return ""
}
