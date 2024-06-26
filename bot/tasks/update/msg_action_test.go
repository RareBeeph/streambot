package update

import (
	"os"
	"streambot/models"
	"streambot/query"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/dismock/pkg/dismock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	db.AutoMigrate(models.All...)
	query.SetDefault(db)

	query.Subscription.Create(&models.Subscription{ChannelID: "123"})

	e := m.Run()

	os.Exit(e)
}

func TestPost(t *testing.T) {
	m := dismock.New(t)

	s, _ := discordgo.New("Bot abc") // the token doesn't have to be valid
	s.Client = m.Client
	s.StateEnabled = false

	// m.Messages(discord.ChannelID(123), 0, nil)
	// m.SendText(discord.Message{ChannelID: discord.ChannelID(123)})

	// TODO: Requires embed data
	m.SendEmbed(discord.Message{ChannelID: discord.ChannelID(123), Embeds: []discord.Embed{}})

	a := msgAction{content: []*discordgo.MessageEmbed{{}}}
	sub, _ := query.Subscription.First()

	// should be a valid post
	a.postContent(s, sub, 0)

	// expect: sub still corresponds to a db row (should be implicit later)
	subs, _ := query.Subscription.Find()
	assert.NotZero(t, len(subs))

	// expect: one new Message row in db table
	msgs, _ := query.Message.Find()
	assert.NotZero(t, len(msgs))

	/*
		// nonexistent channel post
		sub.ChannelID = "ech"
		a.postContent(s, sub, 0)
		// expect: sub no longer corresponds to a db row
		// expect: corresponding Message rows are gone
		// expect: no new Message row in db table
	*/
}
