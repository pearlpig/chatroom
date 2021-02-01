package model

import (
	"log"
	"strings"
)

// Member ...
type Member struct {
	ID       int
	Email    string
	Pwd      string
	Nickname string
}

// LoginStatus ...
type LoginStatus struct {
	Status     int     `json:"status"`
	Msg        string  `json:"msg,omitempty"`
	MemberInfo *Member `json:"member_info,omitempty"`
}

// CheckLogin ...
func CheckLogin(email, pwd string) *LoginStatus {
	if checkEmailExist(email) {
		if r, n := checkPwdCorrect(email, pwd); r {
			return &LoginStatus{Status: 0, Msg: "login", MemberInfo: n}
		}
		return &LoginStatus{Status: 2, Msg: "password is not correct!"}
	}
	return &LoginStatus{Status: 1, Msg: "email is not exist, please sign up!"}
}

func checkPwdCorrect(email string, pwd string) (bool, *Member) {
	db, err := Connect()
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select ID, email, password, nickname from member where email=?", email)
	if err != nil {
		log.Fatal(err)
	}
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
		log.Fatal("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where email=?)", email)
	if err != nil {
		log.Fatal(err)
	}
	exist := -1
	if rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 0 {
		return false
	} else if exist == 1 {
		return true
	} else {
		log.Fatal("This email is not unique in db!")
		return true
	}

}

// func checkNicknameExist(nickname string) {

// }
