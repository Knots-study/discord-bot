package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

func main() {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}
	dg.AddHandler(onMessageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages
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

	td := new(todo)
	order := m.Content

	switch m.Content {
	case "登録":
		s.ChannelMessageSend(m.ChannelID, "登録したいタスクを言ってね\n例：2 部屋の掃除　9")
	case "削除":
		s.ChannelMessageSend(m.ChannelID, "削除したいidを言ってね")
	case "更新":
		s.ChannelMessageSend(m.ChannelID, "更新したいidを言ってね")
	case "表示":
		s.ChannelMessageSend(m.ChannelID, "タスクの一覧を表示するよ")
		operateData(m.Content, td, s, m)
	}

	switch oldmessage {
	case "登録":
		arr := strings.Split(order, " ")
		td.id, td.name, td.level = arr[0], arr[1], arr[2]
		operateData(oldmessage, td, s, m)
	case "削除":
		arr := strings.Split(order, " ")
		td.id = arr[0]
		operateData(oldmessage, td, s, m)
	case "更新":
		arr := strings.Split(order, " ")
		td.id, td.level = arr[0], arr[1]
		operateData(oldmessage, td, s, m)
	}
	oldmessage = m.Content
}

func operateData(order string, data *todo, s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Operate DB")
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

	switch order {

	case "登録":
		cmd = "INSERT INTO todo (id, name, level) VALUES (?, ?, ?)" //SQLインジェクション回避用
		_, err = db.Exec(cmd, data.id, data.name, data.level)
		if err != nil {
			fmt.Println("Fail to insert DB", err)
			s.ChannelMessageSend(m.ChannelID, "idが重複しています")
		} else {
			s.ChannelMessageSend(m.ChannelID, "登録したよ！")
		}

	case "削除":
		cmd = "DELETE FROM todo WHERE id = ?"
		_, err = db.Exec(cmd, data.id)
		if err != nil {
			fmt.Println("Fail to delete DB", err)
		}
		s.ChannelMessageSend(m.ChannelID, "削除したよ")

	case "更新":
		cmd = "UPDATE todo SET level = ? WHERE id = ?"
		_, err = db.Exec(cmd, data.level, data.id)
		if err != nil {
			fmt.Println("Fail to update DB", err)
		}
		s.ChannelMessageSend(m.ChannelID, "更新したよ")

	case "表示":
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
		fmt.Println(count - 1)
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "1⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "2⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "3⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "4⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "5⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "6⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "7⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "8⃣")
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "9⃣")
	default:
		s.ChannelMessageSend(m.ChannelID, "きちんとつぶやいてね!")
	}
}
