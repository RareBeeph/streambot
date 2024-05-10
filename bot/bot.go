package bot

import (
	"os"
	"os/signal"
	"streambot/bot/tasks"
	"streambot/bot/twitch"
	"streambot/config"
	"sync"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
)

type Bot interface {
	Start() error
	Wait()
	Stop()
}

type bot struct {
	session   *discordgo.Session
	conf      *config.Config
	channel   chan os.Signal
	scheduler *cron.Cron
	initOnce  sync.Once
}

func New(conf *config.Config) (b Bot, err error) {
	discord, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		return
	}

	twitch.LoadConfig(conf)

	b = &bot{session: discord, conf: conf}
	return
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

	tasks.All[0].Run() // temp; we are gonna want to run some of our tasks on startup though
	b.scheduler.Start()

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
	b.scheduler.Stop()

	log.Info().Msg("Closing session...")
	b.session.Close() // Technically returns an error but the app is closing anyways
}
