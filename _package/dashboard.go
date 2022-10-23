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

func main() {
	fmt.Println("aiueo")

}

func makeDashboard(session *discordgo.Session, message *discordgo.MessageCreate, title string, description string) *Dashboard {
	ds := new(Dashboard)
	ds.title = title
	ds.description = description
	ds.session = session
	ds.message = message
	return ds
}

func (ds *Dashboard) reload() {
	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedType("rich"),
		Title:       ds.title,
		Description: ds.description,
		Color:       32768,
	}
	// post
	fmt.Println(embed)
}
