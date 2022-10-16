package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token = "Bot " + os.Getenv("Discord-Bot-Token")
)

func main() {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			return
		}
	}(dg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("test")
	fmt.Println(m.Content, m.Author.Username)
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			return
		}
	}
	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			return
		}
	}
}
