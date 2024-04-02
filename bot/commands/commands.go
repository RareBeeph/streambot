package commands

import "github.com/bwmarrin/discordgo"

type handler = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Definition struct {
	Base    *discordgo.ApplicationCommand
	Name    string
	Handler handler
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "helloworld",
		Description: "hello world",
	},
}

var handlers = map[string]handler{
	"helloworld": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "hello world",
			},
		})
	},
}

func GetCommands() []*discordgo.ApplicationCommand {
	return commands
}

func Register(cmd *Definition) {
	cmd.Base.Name = cmd.Name
	commands = append(commands, cmd.Base)
	handlers[cmd.Name] = cmd.Handler
}

func SlashCommandRouter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
