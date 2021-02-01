package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func insert() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	insert, err := db.Exec("INSERT INTO member (id,email,password,nickname,created,updated) VALUES(?,?,?,?,now(),now())", nil, "123@123", "asd", "hi")
	if err != nil {
		panic(err)
	}
	fmt.Println(insert)
}

// Connect ...
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cyberon_chatroom?charset=utf8mb4&parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDB() {
	createDB()
	createMenberTable()
	createChatroomTable()
}

func createDB() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("create database if not exists `cyberon_chatroom` character set 'utf8mb4';")
	if err != nil {
		panic(err)
	}

}
func createMenberTable() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cyberon_chatroom")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("create table if not exists member (id int primary key auto_increment, email varchar(100) unique key not null,password char(129) not null,nickname varchar(255) not null,created datetime default current_timestamp,updated datetime default current_timestamp);")
	if err != nil {
		panic(err)
	}
}
func createChatroomTable() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cyberon_chatroom")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("create table if not exists chatroom (id int primary key auto_increment, title varchar(100) unique key not null, member_id int, created datetime, updated datetime);")
	if err != nil {
		panic(err)
	}
}
