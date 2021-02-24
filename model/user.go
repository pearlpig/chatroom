package model

import (
	"crypto/sha512"
	"fmt"
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
func CheckLogin(form LoginForm) (*ResMember, error) {
	// check format
	if len([]rune(form.Pwd)) < 8 {
		return &ResMember{Status: &ErrStatus{Code: 3, Msg: "Password length should at least 8."}, Data: nil}, nil
	}
	// check email exist
	if exist, err := checkEmailExist(form.Email); err != nil {
		return nil, err
	} else if !exist {
		return &ResMember{Status: &ErrStatus{Code: 1, Msg: "Email is not exist, please sign up!"}, Data: nil}, nil
	}
	// check password correct
	if data, err := checkPwdCorrect(form.Email, form.Pwd); err != nil {
		return nil, err
	} else if data != nil {
		return &ResMember{Status: &ErrStatus{Code: 0, Msg: "Login"}, Data: data}, nil
	}
	return &ResMember{Status: &ErrStatus{Code: 2, Msg: "Password is not correct!"}, Data: nil}, nil
}

// CheckSignup ...
func CheckSignup(form SignupForm) (*ResMember, error) {
	// check format
	if len([]rune(form.Nickname)) > 20 {
		return &ResMember{Status: &ErrStatus{Code: 3, Msg: "Nickname length should at most 20 character."}, Data: nil}, nil
	}
	if len([]rune(form.Nickname)) < 1 {
		return &ResMember{Status: &ErrStatus{Code: 4, Msg: "Nickname should not be empty."}, Data: nil}, nil
	}
	if len([]rune(form.Pwd1)) < 8 {
		return &ResMember{Status: &ErrStatus{Code: 5, Msg: "Password length should at least 8."}, Data: nil}, nil
	}
	if strings.Compare(form.Pwd1, form.Pwd2) != 0 {
		return &ResMember{Status: &ErrStatus{Code: 6, Msg: "Please check the confirmed password."}, Data: nil}, nil
	}
	// check if email is exist or not
	if exist, err := checkEmailExist(form.Email); err != nil {
		return nil, err
	} else if exist {
		return &ResMember{Status: &ErrStatus{Code: 1, Msg: "Email is exist, please try another one."}, Data: nil}, nil
	}
	// check nickname is exist or not
	if exist, err := checkNicknameExist(form.Email); err != nil {
		return nil, err
	} else if exist {
		return &ResMember{Status: &ErrStatus{Code: 2, Msg: "Nickname is exist, please try another one."}, Data: nil}, nil
	}

	id, err := doSignup(&form)
	if err != nil {
		return nil, err
	}
	return &ResMember{Status: &ErrStatus{Code: 0, Msg: "Signup!"}, Data: &Member{ID: id, Email: form.Email, Pwd: form.Pwd1, Nickname: form.Nickname}}, nil

}

func doSignup(form *SignupForm) (int, error) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
		return -1, err
	}
	defer db.Close()
	pwd := fmt.Sprintf("%x", sha512.Sum512([]byte(form.Pwd1)))
	res, err := db.Exec("insert into member(ID, email, password, nickname) values(?, ?, ?, ?)", nil, form.Email, pwd, form.Nickname)
	if err != nil {
		log.Println("insert: ", err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error: ", err)
		return -1, err
	}
	return int(id), nil

}
func checkPwdCorrect(email string, pwd string) (*Member, error) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select ID, email, password, nickname from member where email=?", email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	dbMember := &Member{}
	if rows.Next() {
		rows.Scan(&dbMember.ID, &dbMember.Email, &dbMember.Pwd, &dbMember.Nickname)
	}
	log.Println(fmt.Sprintf("%x", sha512.Sum512([]byte(pwd))))
	if strings.Compare(fmt.Sprintf("%x", sha512.Sum512([]byte(pwd))), dbMember.Pwd) == 0 {
		return dbMember, nil
	}
	log.Println("password is incorrect")
	return nil, nil
}

func checkEmailExist(email string) (bool, error) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where email=?)", email)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false, nil
	} else if exist == 1 {
		return true, nil
	} else {
		return false, fmt.Errorf("this email is not unique in db")
	}
}

func checkNicknameExist(nickname string) (bool, error) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where nickname=?)", nickname)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false, nil
	} else if exist == 1 {
		return true, nil
	} else {
		return false, fmt.Errorf("this nickname is not unique in db")
	}
}

// CheckCreate return room info if it is successfully created
func CheckCreate(form CreateRoomForm) (*ResCreateRoom, error) {
	// check format
	if len([]rune(form.RoomName)) > 20 {
		return &ResCreateRoom{Status: &ErrStatus{Code: 3, Msg: "Room name length should at most 20 character."}, Data: nil}, nil
	} else if len([]rune(form.RoomName)) < 0 {
		return &ResCreateRoom{Status: &ErrStatus{Code: 2, Msg: "Room name should not be empty."}, Data: nil}, nil
	}

	// check if room is exist
	if exist, err := checkRoomExist(form.RoomName); err != nil {
		return nil, err
	} else if exist {
		return &ResCreateRoom{Status: &ErrStatus{Code: 1, Msg: "Room is exist, please try another one."}, Data: nil}, nil
	}

	// create room
	roomID, err := doCreateRoom(&form)
	if err != nil {
		log.Println("Error: create room:", err)
		return nil, err
	}
	return &ResCreateRoom{Status: &ErrStatus{Code: 0, Msg: "Create room!"}, Data: &Room{ID: roomID, Title: form.RoomName, MemberID: form.MemberID, Nickname: ""}}, nil
}

func doCreateRoom(form *CreateRoomForm) (int, error) {
	db, err := Connect()
	if err != nil {
		log.Println("connect: ", err)
		return -1, err
	}
	defer db.Close()

	res, err := db.Exec("insert into chatroom(ID, title, member_id) values(?, ?, ?)", nil, form.RoomName, form.MemberID)
	if err != nil {
		log.Println("insert: ", err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error: ", err)
		return -1, err
	}
	return int(id), nil
}

func checkRoomExist(roomName string) (bool, error) {
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("select exists (select id from chatroom where title=?)", roomName)
	if err != nil {
		log.Println("Error: ", err)
		return false, err
	}
	defer rows.Close()
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false, nil
	} else if exist == 1 {
		return true, nil
	} else {
		return false, fmt.Errorf("this roomName is not unique in db")
	}
}

// GetRoom ...
func GetRoom(page int) ([]*Room, error) {
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select c.id,title,member_id,nickname from chatroom c join member m on c.member_id=m.id where floor((c.id-1)/10)=?", page-1)
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}
	defer rows.Close()

	roomList := []*Room{}
	for rows.Next() {
		room := &Room{}
		rows.Scan(&room.ID, &room.Title, &room.MemberID, &room.Nickname)
		roomList = append(roomList, room)
	}

	return roomList, nil
}

// GetRoomName ...
func GetRoomName(roomNum int) (string, error) {
	roomName := ""
	db, err := Connect()
	if err != nil {
		log.Println("Connected Error: ", err)
		return roomName, err
	}
	defer db.Close()

	// log.Println(page)
	rows, err := db.Query("select title from chatroom where id=?", roomNum)
	if err != nil {
		log.Println("Error: ", err)
		return roomName, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&roomName)
	}
	if roomName == "" {
		return roomName, fmt.Errorf("room is not exist")
	}
	return roomName, nil
}
