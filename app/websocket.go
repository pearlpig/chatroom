package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// var upgrader = &websocket.Upgrader{
// 	//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// Message ...
type Message struct {
	Nickname string `json:"nickname,omitempty"`
	Msg      string `json:"msg,omitempty"`
	RoomID   int    `json:"room_id,omitempty"`
	Status   int    `json:"status"`
}

// SocketConn ...
type SocketConn struct {
	Conn   *websocket.Conn
	Cookie *Cookies
}

var i = 0
var connList = make(map[int]map[int]*SocketConn)
var number = make(chan int)
var cMsg = make(chan Message)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(w, r)
	nickname := cookie.Nickname
	roomID := cookie.RoomID
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	ticket := <-number
	log.Println("門票號碼：", ticket)
	connList[roomID] = make(map[int]*SocketConn)
	connList[roomID][ticket] = &SocketConn{Conn: conn, Cookie: cookie}
	defer func(ticket int, conn *websocket.Conn) {
		log.Println("disconnect !!")
		conn.Close()
		conn = nil
		delete(connList[roomID], ticket)
	}(ticket, conn)
	cMsg <- Message{RoomID: roomID, Nickname: nickname, Status: 1}
	for {
		log.Println("listening socket")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			fmt.Println(msg)
			break
		}
		fmt.Printf("Got message: %s\n", msg)
		cMsg <- Message{RoomID: roomID, Msg: string(msg), Nickname: nickname, Status: 2}
	}
	cMsg <- Message{Nickname: nickname, Status: 0}
}
func connRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("connect room handler")
	cookie := getCookie(w, r)
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])

	if err != nil {
		log.Println("Error: ", err, " roomID is not number.")
	}
	cookie.RoomID = roomID
	setCookie(&w, r, cookie)
}
func disconnRoomHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("disconnect room handler")
	cookie := getCookie(w, r)
	cookie.RoomID = -1
	setCookie(&w, r, cookie)
}
func dispensor() {
	for {
		number <- i
		i++
	}
}
func broker() {
	for {
		msg := <-cMsg
		log.Println(msg.RoomID)
		for i, conn := range connList[msg.RoomID] {
			if conn != nil {
				log.Println("send", i)
				if err := conn.Conn.WriteJSON(msg); err != nil {
					fmt.Println("write err: ", err)
				}
			} else {
				log.Println("disconnected.................")
			}
		}
	}
}
