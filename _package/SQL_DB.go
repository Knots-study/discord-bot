package _package

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

/*
データベースを管理するためのパッケージ
*/

var (
	emojis   = InitEmojis()
	flag_new = 0
)

type todo struct {
	name  string
	level string
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

	cmd := "CREATE TABLE IF NOT EXISTS todo (name text primary key, level text)" //AUTOINCREMENT制約(セキュリティ的にヤバめ，どうするか)
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB", err)
	}
}

func InsertDB(words []string, db *sql.DB) {
	cmd := "INSERT INTO todo (name, level) VALUES (?, ?)"
	_, err := db.Exec(cmd, words[0], words[1])
	if err != nil {
		fmt.Println("Fail to insert DB", err)
	}
}

func SelectDB(db *sql.DB) (int, discordgo.MessageEmbed) {
	cmd := "SELECT * FROM todo"
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

	td := new(todo)
	comment := ""
	count := 0
	for rows.Next() {
		err := rows.Scan(&td.name, &td.level)
		if err != nil {
			fmt.Println(err)
		}
		comment += emojis.Numbers[count] + " タスク名: " + td.name + " 優先度: " + td.level + "\n"
		count += 1
	}
	embed := discordgo.MessageEmbed{Title: "ToDoリスト(※10個まで登録可)", Description: comment, Color: 1752220}
	return count, embed
}

func DeleteDB(db *sql.DB, Name string) {
	cmd := "DELETE FROM todo WHERE name = (select name from todo limit 1 offset ?-1)"
	_, err := db.Exec(cmd, Name)
	if err != nil {
		fmt.Println("Fail to delete DB", err)
	}
}

func UpdateDB(db *sql.DB, Name string) { //一旦保留した関数
	cmd := "UPDATE todo SET level = ? WHERE name = (select name from todo limit 1 offset ?-1)"
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
