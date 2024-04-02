package bot

import (
	"os"
	"os/signal"
	"streambot/bot/commands"
	"streambot/config"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
)

type Bot interface {
	Start() error
	Wait()
	Stop()
}

type bot struct {
	session            *discordgo.Session
	channel            chan os.Signal
	initOnce           sync.Once
	registeredCommands []*discordgo.ApplicationCommand
}

func New(conf *config.Config) (Bot, error) {
	session, err := discordgo.New("Bot " + conf.Token)
	return &bot{session: session}, err
}

func (b *bot) Start() error {
	b.init()

	guildID := "" // TODO: get this from the yml
	err := b.registerCommands(guildID)
	if err != nil {
		return err
	}

	signal.Notify(b.channel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return nil
}

func (b *bot) Wait() {
	b.init()
	<-b.channel
}

func (b *bot) Stop() {
	b.init()

	guildID := "" // TODO: get this from the yml
	b.unregisterCommands(guildID)

	log.Info().Msg("Closing session...")
	b.session.Close() // Technically returns an error but the app is closing anyways
	log.Info().Msg("Shutting down.")
}

func (b *bot) registerCommands(guildID string) error {
	commandList := commands.GetCommands()
	b.registeredCommands = make([]*discordgo.ApplicationCommand, len(commandList))

	for i, c := range commandList {
		reg, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, c)
		if err != nil {
			return err
		}
		b.registeredCommands[i] = reg
	}
	return nil
}

func (b *bot) unregisterCommands(guildID string) error {
	// unsafe:
	// registeredCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, guildID)
	// if err != nil {
	// 	return err
	// }

	for _, c := range b.registeredCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, guildID, c.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bot) init() {
	b.initOnce.Do(func() {
		// Register the messageCreate func as a callback for MessageCreate events.
		b.session.AddHandler(messageCreate)
		b.session.AddHandler(commands.SlashCommandRouter)

		// In this example, we only care about receiving message events.
		b.session.Identify.Intents = discordgo.IntentsGuildMessages

		// Open a websocket connection to Discord and begin listening.
		err := b.session.Open()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to discord")
		}

		b.channel = make(chan os.Signal, 1)
	})
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
