package main

import (
	"chatroom/app"
	"chatroom/model"
)

func main() {
	model.InitDB()
	app.Server()

}
