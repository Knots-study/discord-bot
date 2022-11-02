package _package

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Dashboard struct {
	title       string
	description string
	session     *discordgo.Session
	message     *discordgo.MessageCreate
	channelID   string
	messageID   string
}

func MakeDashboard(session *discordgo.Session, message *discordgo.MessageCreate, title string, description string) *Dashboard {
	ds := new(Dashboard)
	ds.title = title
	ds.description = description
	ds.session = session
	ds.message = message
	return ds
}

func (ds *Dashboard) Reload() {
	//var fields []*discordgo.MessageEmbedField

	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedType("rich"),
		Title:       ds.title,
		Description: ds.description,
		Color:       32768,
	}
	// post
	fmt.Println(embed)
}

func (ds *Dashboard) Renew() {
	err := ds.session.ChannelMessageDelete(ds.channelID, ds.messageID)
	if err != nil {
		//Log.Err(err)
	}
	ds.messageID = ""
	ds.Reload()
}
