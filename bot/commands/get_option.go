package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/teris-io/shortid"
)

func get_option(s *discordgo.Session, i *discordgo.InteractionCreate, prompt string, options []discordgo.SelectMenuOption) (optionId string, iOut *discordgo.InteractionCreate) {
	selectId := shortid.MustGenerate()
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: prompt,
			Flags:   discordgo.MessageFlagsEphemeral,
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

	unregisterHandler := s.AddHandler(func(s *discordgo.Session, i2 *discordgo.InteractionCreate) {
		if i2.Type != discordgo.InteractionMessageComponent {
			return
		}

		if id := i2.MessageComponentData().CustomID; id == selectId {
			s.InteractionResponseDelete(i.Interaction)
			iOut = i2
			c <- i2.MessageComponentData().Values[0]
		}
	})
	defer unregisterHandler()

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		iOut = i // already (attempted to be) responded to, but at least it's not nil
	} else {
		optionId = <-c
	}

	return
}
