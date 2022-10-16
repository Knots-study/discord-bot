package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var (
	Token = "Bot " + os.Getenv("Discord-Bot-Token")
)

type todo struct {
	id    string
	name  string
	level string
}

func main() {
	td := new(todo)
	fmt.Println("登録するデータを入力してください")
	_, err := fmt.Scan(&td.id, &td.name, &td.level)
	if err != nil {
		fmt.Println(err)
		return
	}
	//登録するユーザ情報

	//データベースに情報を登録
	err = saveData(td)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Succeeded to save your task")
}

func saveData(data *todo) error {
	fmt.Println("Save data to DB")
	db, err := sql.Open("sqlite3", "todo_database.db") //データベースに接続
	if err != nil {
		fmt.Println("open sql")
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db) //データベースの接続解除

	todo, err := db.Prepare("CREATE TABLE IF NOT EXISTS todo (id text PRIMARY KEY, name text, level text)") //データベースの準備
	if err != nil {
		log.Println("prepare sql(create table)")
		return err
	}
	defer func(todo *sql.Stmt) {
		err := todo.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(todo) //データベースを閉じる
	_, err = todo.Exec()
	if err != nil {
		fmt.Println(err)
	}

	todo, err = db.Prepare("INSERT INTO todo (id, name, level) VALUES ($1, $2, $3)")
	if err != nil {
		log.Println("[-]sql.Prepare (INSERT INTO todo)")
		return err
	}
	defer func(todo *sql.Stmt) {
		err := todo.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(todo)

	_, err = todo.Exec(data.id, data.name, data.level)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
