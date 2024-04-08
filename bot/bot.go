package bot

import (
	"os"
	"os/signal"
	"streambot/bot/commands"
	"streambot/config"
	"streambot/models"
	"streambot/query"
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
	session  *discordgo.Session
	conf     *config.Config
	channel  chan os.Signal
	initOnce sync.Once
}

func New(conf *config.Config) (Bot, error) {
	session, err := discordgo.New("Bot " + conf.Token)
	return &bot{session: session, conf: conf}, err
}

func (b *bot) Start() error {
	b.init()

	// cleanup. should do nothing on a clean start,
	// but unregister any lingering commands from a dirty start
	b.unregisterCommands()

	err := b.registerCommands()
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

	b.unregisterCommands()

	log.Info().Msg("Closing session...")
	b.session.Close() // Technically returns an error but the app is closing anyways
	log.Info().Msg("Shutting down.")
}

func (b *bot) registerCommands() error {
	commandList := commands.GetCommands()

	for _, c := range commandList {
		reg, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.conf.GuildID, c)
		if err != nil {
			return err
		}

		err = query.RegisteredCommand.Create(&models.RegisteredCommand{ID: reg.ID, GuildID: reg.GuildID})
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) unregisterCommands() error {
	rc := query.RegisteredCommand

	allcommands, err := rc.Find()
	if err != nil {
		return err
	}

	for _, c := range allcommands {
		err = b.session.ApplicationCommandDelete(b.session.State.User.ID, c.GuildID, c.ID)
		if err != nil {
			return err
		}

		_, err = rc.Where(rc.ID.Eq(c.ID)).Delete()
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
