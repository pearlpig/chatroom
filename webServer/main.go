package main

import (
	"chatroom/app"
)

//go:generate docker rm -f chatroom_server_web
//go:generate docker build -t chatroom/web:latest .
//go:generate docker run -d --name=chatroom_server_web --link chatroom_server_db -p 8080:8080 chatroom/web:latest
func main() {
	app.Server()
}
