package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	Token = "Bot " + os.Getenv("Discord-Bot-Token")
)

type todo struct {
	id    string
	name  string
	level string
}

var oldmessage = " "
var Bot_Message_ID = ""

func main() {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	Emojis.InitEmojis()
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
	if m.ChannelID != os.Getenv("Discord-Bot-Todo-ChannelID") { //チャンネル外での発言
		return
	}
	if m.Author.ID == s.State.User.ID { //本人の発言
		return
	}

	db, err := sql.Open("sqlite3", "todo_database.db") //データベースに接続
	if err != nil {
		fmt.Println("Fail to open DB", err)
	}
	defer func(db *sql.DB) { //必ず閉じる
		err := db.Close()
		if err != nil {
			fmt.Println("Fail to close DB", err)
		}
	}(db) //データベースの接続解除

	cmd := "CREATE TABLE IF NOT EXISTS todo (id text PRIMARY KEY, name text, level text)" //DBが存在しない場合，新たにCREATE
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB", err)
	}

	cmd = "SELECT * FROM todo"
	rows, err := db.Query(cmd) //複数の検索結果を取得するため，Query
	if err != nil {
		fmt.Println("Fail to select DB", err)
	}
	defer func(rows *sql.Rows) { //絶対に閉じる
		err := rows.Close()
		if err != nil {
			fmt.Println("Fail to close selecting DB", err)
		}
	}(rows)

	var td todo
	comment := ""
	count := 1
	for rows.Next() {
		err := rows.Scan(&td.id, &td.name, &td.level)
		if err != nil {
			fmt.Println(err)
		}
		comment += "ID: " + strconv.Itoa(count) + " タスク名: " + td.name + " 優先度: " + td.level + "\n"
		count += 1
	}
	embed := discordgo.MessageEmbed{Title: "ToDoリスト", Description: comment, Color: 1752220}
	message, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error")
	}
	emoji := [...]string{"1⃣", "2⃣", "3⃣", "4⃣", "5⃣", "6⃣", "7⃣", "8⃣", "9⃣"}
	for i := 0; i < count-1; i++ {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, emoji[i])
	}
	_ = s.MessageReactionAdd(m.ChannelID, message.ID, "📝")
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	if m.ChannelID != os.Getenv("Discord-Bot-Todo-ChannelID") { //チャンネル外での発言
		return
	}
	if m.UserID == s.State.User.ID {
		Bot_Message_ID = m.MessageID
		return
	}
	if m.MessageID == Bot_Message_ID {
		db, err := sql.Open("sqlite3", "todo_database.db") //データベースに接続
		if err != nil {
			fmt.Println("Fail to open DB", err)
		}
		defer func(db *sql.DB) { //必ず閉じる
			err := db.Close()
			if err != nil {
				fmt.Println("Fail to close DB", err)
			}
		}(db) //データベースの接続解除

		fmt.Println(m.Emoji.Name)
		fmt.Println(m.Emoji.ID)
		//絵文字を押した時に削除か更新か選びたい
		cmd := "DELETE FROM todo WHERE id = (select id from todo limit 1 offset ?-1)"
		_, err = db.Exec(cmd, m.Emoji.Name)
		if err != nil {
			fmt.Println("Fail to delete DB", err)
		}
		s.ChannelMessageSend(m.ChannelID, "削除したよ")
	}
}
