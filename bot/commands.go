package bot

import (
	"github.com/rs/zerolog/log"
	"streambot/bot/commands"
	"streambot/models"
	"streambot/query"
)

func (b *bot) registerCommands() error {
	commandList := commands.All()

	for _, c := range commandList {
		reg, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.conf.GuildID, c)
		if err != nil {
			log.Print(err)
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
