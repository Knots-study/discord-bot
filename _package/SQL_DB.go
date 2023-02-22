package _package

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
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

func ConnectDB() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("環境変数が読み込めませんでした")
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("error to Connect DB", err)
	}
	c := mysql.Config{
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      "localhost:3306",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		fmt.Println("Fail to open DB", err)
	}

	// テーブルがない場合，作成する
	cmd := "CREATE TABLE IF NOT EXISTS todobot.Todo (id integer AUTO_INCREMENT, name text, deadline text, level text, UnixDead int, primary key(id))"
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB", err)
	}
	cmd = "CREATE TABLE IF NOT EXISTS todobot.TodoNotice (id integer AUTO_INCREMENT, noticeTime integer, name text, primary key(id))" //通知用のデータベース
	_, err = db.Exec(cmd)
	if err != nil {
		fmt.Println("Fail to create DB(todo_notice)", err)
	}
	return db
}

func InsertDB(words []string, db *sql.DB) {
	cmd := "INSERT INTO todobot.Todo (id, name, deadline, level, UnixDead) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(cmd, nil, words[0], words[1], words[2], CalcTime(words[1])) //deadlineからUnixdeadを計算
	if err != nil {
		fmt.Println("Fail to insert DB", err)
	}
	cmd = "INSERT INTO todobot.TodoNotice (id, noticeTime, name) VALUES (?, ?, ?)"
	var PushTime = InformCnt(words[1], words[2]) //deadline, level
	for _, tm := range PushTime {
		_, err = db.Exec(cmd, nil, words[0], tm, words[2])
	}
}

func SelectDB(db *sql.DB) (int, discordgo.MessageEmbed) {
	cmd := "DELETE FROM todobot.Todo WHERE UnixDead < ?" //締め切りが過ぎたタスクを自動削除
	_, err := db.Exec(cmd, int(time.Now().Unix()))
	if err != nil {
		fmt.Println("Fail to Delete contents", err)
	}

	cmd = "SELECT * FROM todobot.Todo order by UnixDead" // 残り時間が少ないタスクを上に表示させる
	rows, err := db.Query(cmd)                           //複数の検索結果を取得するため，Query
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
	cmd := "DELETE FROM todobot.Todo WHERE id = (select id from ( select id from todobot.Todo limit 1 offset ?) temp)"
	fmt.Println(Name)
	_, err := db.Exec(cmd, Name)
	if err != nil {
		fmt.Println("Fail to delete DB", err)
	}
}

func UpdateDB(db *sql.DB, Name string) { //一旦保留した関数(後で作る)
	cmd := "UPDATE todobot.Todo SET level = ? WHERE id = (select id from todobot.Todo limit 1 offset ?)" // サブクエリのfrom句と最新のターゲットが両方同じテーブルを指定するとエラー(上のように書く)
	_, err := db.Exec(cmd, Name)
	if err != nil {
		fmt.Println("Fail to delete DB", err)
	}
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println("Fail to close DB", err)
	} //データベースの接続解除
}
