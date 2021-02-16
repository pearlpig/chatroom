package model

import (
	"log"
	"strings"
)

// LoginForm ...
type LoginForm struct {
	Email string
	Pwd   string
}

// SignupForm ...
type SignupForm struct {
	Email    string
	Nickname string
	Pwd1     string
	Pwd2     string
}

// CreateRoomForm ...
type CreateRoomForm struct {
	RoomName string
	MemberID int
}

// Member ...
type Member struct {
	ID       int
	Email    string
	Pwd      string
	Nickname string
}

// Room ...
type Room struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	MemberID int    `json:"member_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

// ErrStatus ...
type ErrStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// ResCreateRoom ...
type ResCreateRoom struct {
	Status *ErrStatus `json:"status,omitempty"`
	Data   *Room      `json:"data,omitempty"`
}

// ResMember ...
type ResMember struct {
	Status *ErrStatus `json:"status,omitempty"`
	Data   *Member    `json:"data,omitempty"`
}

// CheckLogin ...
func CheckLogin(form LoginForm) *ResMember {

	if checkEmailExist(form.Email) {
		if r, n := checkPwdCorrect(form.Email, form.Pwd); r {
			return &ResMember{Status: &ErrStatus{Code: 0, Msg: "Login"}, Data: n}
		}
		return &ResMember{Status: &ErrStatus{Code: 2, Msg: "Password is not correct!"}, Data: nil}
	}

	return &ResMember{Status: &ErrStatus{Code: 1, Msg: "Email is not exist, please sign up!"}, Data: nil}
}

// CheckSignup ...
func CheckSignup(form SignupForm) *ResMember {
	log.Println(form)
	// log.Println(checkEmailExist(form.Email))
	if checkEmailExist(form.Email) {
		return &ResMember{Status: &ErrStatus{Code: 1, Msg: "Email is exist, please try another one."}, Data: nil}
	}
	if checkNicknameExist(form.Nickname) {
		return &ResMember{Status: &ErrStatus{Code: 2, Msg: "Nickname is exist, please try another one."}, Data: nil}
	}
	id := doSignup(&form)
	return &ResMember{Status: &ErrStatus{Code: 0, Msg: "Signup!"}, Data: &Member{ID: id, Email: form.Email, Pwd: form.Pwd1, Nickname: form.Nickname}}
}

func doSignup(form *SignupForm) int {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
	}
	defer db.Close()

	res, err := db.Exec("insert into member(ID, email, password, nickname) values(?, ?, ?, ?)", nil, form.Email, form.Pwd1, form.Nickname)
	if err != nil {
		log.Println("insert: ", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error: ", err)
	}
	return int(id)

}

// CheckCreate ...
func CheckCreate(form CreateRoomForm) *ResCreateRoom {
	log.Println(form)
	// log.Println(checkEmailExist(form.Email))
	if checkRoomExist(form.RoomName) {
		return &ResCreateRoom{Status: &ErrStatus{Code: 1, Msg: "Room is exist, please try another one."}, Data: nil}
	}
	roomID := doCreateRoom(&form)
	return &ResCreateRoom{Status: &ErrStatus{Code: 0, Msg: "Create room!"}, Data: &Room{ID: roomID, Title: form.RoomName, MemberID: form.MemberID, Nickname: ""}}
}

func doCreateRoom(form *CreateRoomForm) int {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
	}
	defer db.Close()

	res, err := db.Exec("insert into chatroom(ID, title, member_id) values(?, ?, ?)", nil, form.RoomName, form.MemberID)
	if err != nil {
		log.Println("insert: ", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error: ", err)
	}
	return int(id)
}
func checkPwdCorrect(email string, pwd string) (bool, *Member) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select ID, email, password, nickname from member where email=?", email)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	dbMember := &Member{}
	if rows.Next() {
		rows.Scan(&dbMember.ID, &dbMember.Email, &dbMember.Pwd, &dbMember.Nickname)
	}
	log.Println(dbMember)
	if strings.Compare(pwd, dbMember.Pwd) == 0 {
		return true, dbMember
	}
	return false, nil
}

func checkEmailExist(email string) bool {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where email=?)", email)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false
	} else if exist == 1 {
		return true
	} else {
		log.Println("This email is not unique in db!")
		return true
	}

}

func checkNicknameExist(nickname string) bool {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where nickname=?)", nickname)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false
	} else if exist == 1 {
		return true
	} else {
		log.Println("This nickname is not unique in db!")
		return true
	}
}

func checkRoomExist(roomName string) bool {
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select exists (select id from chatroom where title=?)", roomName)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false
	} else if exist == 1 {
		return true
	} else {
		log.Println("This roomName is not unique in db!")
		return true
	}
}

// GetRoom ...
func GetRoom(page int) []*Room {
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
	}
	defer db.Close()
	// log.Println(page)
	rows, err := db.Query("select c.id,title,member_id,nickname from chatroom c join member m on c.member_id=m.id where floor((c.id-1)/10)=?", page-1)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer rows.Close()
	roomList := []*Room{}
	for rows.Next() {
		room := &Room{}
		rows.Scan(&room.ID, &room.Title, &room.MemberID, &room.Nickname)
		roomList = append(roomList, room)
	}

	return roomList
}

// GetRoomName ...
func GetRoomName(roomNum int) string {
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
	}
	defer db.Close()
	// log.Println(page)
	rows, err := db.Query("select title from chatroom where id=?", roomNum)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer rows.Close()
	roomName := ""
	for rows.Next() {
		rows.Scan(&roomName)
	}

	return roomName
}
