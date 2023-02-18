package _package

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

/*
データベースを管理するためのパッケージ
*/

var (
	emojis   = InitEmojis()
	flag_new = 0
)

type Todo struct {
	id       int
	name     string
	deadline string
	level    string
	UnixDead int
}

type TodoNotice struct {
	Todo
	noticeTime int
}

func CreateDB() {
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

	cmd := "CREATE TABLE IF NOT EXISTS Todo (id integer primary key, name text, deadline text, level text, UnixDead int)"
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB", err)
	}
	cmd = "CREATE TABLE IF NOT EXISTS TodoNotice (id integer primary key, noticeTime integer, name text)" //通知用のデータベース
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB(todo_notice)", err)
	}
}

func InsertDB(words []string, db *sql.DB) {
	cmd := "INSERT INTO Todo (id, name, deadline, level, UnixDead) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(cmd, nil, words[0], words[1], words[2], CalcTime(words[1])) //deadlineからUnixdeadを計算
	if err != nil {
		fmt.Println("Fail to insert DB", err)
	}
	cmd = "INSERT INTO TodoNotice (id, noticeTime, name) VALUES (?, ?, ?)"
	var PushTime = InformCnt(words[1], words[2]) //deadline, level
	for _, tm := range PushTime {
		_, err = db.Exec(cmd, nil, words[0], tm, words[2])
	}
}

func SelectDB(db *sql.DB) (int, discordgo.MessageEmbed) {
	cmd := "DELETE FROM Todo WHERE UnixDead < ?" //締め切りが過ぎたタスクを自動削除
	_, err := db.Exec(cmd, int(time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
	}

	cmd = "SELECT * FROM Todo order by UnixDead" // 残り時間が少ないタスクを上に表示させる
	rows, err := db.Query(cmd)                   //複数の検索結果を取得するため，Query
	if err != nil {
		fmt.Println("Fail to select DB", err)
	}
	defer func(rows *sql.Rows) { //絶対に閉じる
		err := rows.Close()
		if err != nil {
			fmt.Println("Fail to close selecting DB", err)
		}
	}(rows)

	td := new(Todo)
	comment := ""
	count := 0
	for rows.Next() {
		err := rows.Scan(&td.id, &td.name, &td.deadline, &td.level, &td.UnixDead)
		if err != nil {
			fmt.Println(err)
		}
		comment += emojis.Numbers[count] + " タスク名: " + td.name + " 締め切り: " + td.deadline + " 優先度: " + td.level + "\n"
		count += 1
	}
	embed := discordgo.MessageEmbed{Title: "ToDoリスト(※10個まで登録可)", Description: comment, Color: 1752220}
	return count, embed
}

func DeleteStampDB(db *sql.DB, Name string) {
	cmd := "DELETE FROM todo WHERE id = (select id from todo limit 1 offset ?-1)" //バグはココ
	fmt.Println(Name)
	_, err := db.Exec(cmd, Name)
	if err != nil {
		fmt.Println("Fail to delete DB", err)
	}
}

func UpdateDB(db *sql.DB, Name string) { //一旦保留した関数(後で作る)
	cmd := "UPDATE todo SET level = ? WHERE id = (select id from todo limit 1 offset ?-1)"
	_, err := db.Exec(cmd, Name)
	if err != nil {
		fmt.Println("Fail to delete DB", err)
	}
}

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "todo_database.db") //データベースに接続
	if err != nil {
		fmt.Println("Fail to open DB", err)
	}
	return db
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println("Fail to close DB", err)
	} //データベースの接続解除
}
