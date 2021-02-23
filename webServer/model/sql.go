package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Connect ...
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "asd:1234@tcp(chatroom_server_db)/cyberon_chatroom?charset=utf8mb4&parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}
