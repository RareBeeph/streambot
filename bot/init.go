package bot

import (
	"os"
	"streambot/bot/commands"
	"streambot/bot/tasks"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
)

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

		b.scheduler = cron.New()
		for _, t := range tasks.All {
			t.BindSession(b.session)
			b.scheduler.AddFunc(t.Spec, t.Run)
		}
	})
}
