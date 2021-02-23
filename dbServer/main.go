package main

//go:generate docker rm -f chatroom_server_db
//go:generate docker build -t chatroom/db:latest .
//go:generate docker run -d --name=chatroom_server_db chatroom/db:latest

func main() {

}
