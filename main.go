package main

import (
	"fmt"
	. "github.com/Knots-study/discord-bot/_package"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	Token          = "Bot " + os.Getenv("Discord-Bot-Token")
	Bot_Message_ID = ""
	emojis         = InitEmojis()
	flag_new       = 0
)

func main() {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	dg.AddHandler(onMessageCreate)
	dg.AddHandler(messageReactionAdd)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			fmt.Println("error closing connection,", err)
		}
	}(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != os.Getenv("Discord_Bot_Todo_ChannelID") { //チャンネル外での発言
		return
	}
	if m.Author.ID == s.State.User.ID { //本人の発言
		return
	}

	db := OpenDB()    //DBを起動
	defer CloseDB(db) //DBは必ず閉じる

	CreateDB() //起動時にDBのテーブルが未作成の場合，作成する

	if flag_new == 1 { //Newボタンが押された時のみ，値を追加する
		arr := strings.Split(m.Content, " ")
		if len(arr) != 2 {
			s.ChannelMessageSend(m.ChannelID, "きちんと入力してください")
			s.ChannelMessageSend(m.ChannelID, "登録したいタスクを言ってね(例:部屋の掃除 9)")
			return
		}
		InsertDB(arr, db)
		flag_new = 0
	}

	count, embed := SelectDB(db) //一覧表示
	message, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error")
	}
	for i := 0; i < count; i++ {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, emojis.Numbers[i])
	}
	_ = s.MessageReactionAdd(m.ChannelID, message.ID, emojis.New)
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != os.Getenv("Discord_Bot_Todo_ChannelID") { //チャンネル外での発言
		return
	}
	if m.UserID == s.State.User.ID {
		Bot_Message_ID = m.MessageID
		return
	}
	if m.MessageID == Bot_Message_ID { //絵文字を押した時に削除か更新か選びたい
		name := m.Emoji.Name
		switch name {
		case emojis.New: //登録
			flag_new = 1
			s.ChannelMessageSend(m.ChannelID, "登録したいタスクを言ってね(例:部屋の掃除 9)")
		default:
			db := OpenDB()    //DBを起動
			defer CloseDB(db) //DBは必ず閉じる
			DeleteDB(db, name)
			s.ChannelMessageSend(m.ChannelID, "削除したよ")
			//UpdateDB(db, name)　//updateをすると，通知時間を計算し直す必要がある為，一旦保留

			count, embed := SelectDB(db) //一覧表示
			message, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error")
			}
			for i := 0; i < count; i++ {
				_ = s.MessageReactionAdd(m.ChannelID, message.ID, emojis.Numbers[i])
			}
			_ = s.MessageReactionAdd(m.ChannelID, message.ID, emojis.New)
		}
	}
}
