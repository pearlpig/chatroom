all: 
	-docker-compose up

network: 
	-docker network create -d bridge chatroom
db: 
	-make network
	-docker rm -f chatroom_server_db
	-docker build -t chatroom/db:latest ./dbServer
	-docker run -d --name=chatroom_server_db --network=chatroom chatroom/db:latest

web:
	-go build ./webServer/.
	-make network
	-docker rm -f chatroom_server_web
	-docker build -t chatroom/web:latest ./webServer
	-docker run -d --name=chatroom_server_web --network=chatroom -p 8080:8080 chatroom/web:latest
web-run:
	-go run ./webServer/main.go
clean:
	-docker stop chatroom_chatroom_server_db_1
	-docker stop chatroom_chatroom_server_web_1
	-docker rm chatroom_chatroom_server_db_1
	-docker rm chatroom_chatroom_server_web_1
	-docker network rm chatroom