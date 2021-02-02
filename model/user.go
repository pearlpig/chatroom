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
func CheckLogin(form LoginForm) *LoginStatus {

	if checkEmailExist(form.Email) {
		if r, n := checkPwdCorrect(form.Email, form.Pwd); r {
			return &LoginStatus{Status: 0, Msg: "Login", MemberInfo: n}
		}
		return &LoginStatus{Status: 2, Msg: "Password is not correct!"}
	}
	return &LoginStatus{Status: 1, Msg: "Email is not exist, please sign up!"}
}

// CheckSignup ...
func CheckSignup(form SignupForm) *LoginStatus {
	log.Println(form)
	// log.Println(checkEmailExist(form.Email))
	if checkEmailExist(form.Email) {
		return &LoginStatus{Status: 1, Msg: "Email is exist, please try another one."}
	}
	if checkNicknameExist(form.Nickname) {
		return &LoginStatus{Status: 2, Msg: "Nickname is exist, please try another one."}
	}
	id := doSignup(&form)
	return &LoginStatus{Status: 0, Msg: "Signup!", MemberInfo: &Member{ID: id, Email: form.Email, Pwd: form.Pwd1, Nickname: form.Nickname}}
}
func doSignup(form *SignupForm) int {
	db, err := Connect()
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer db.Close()

	res, err := db.Exec("insert into member(ID, email, password, nickname) values(?, ?, ?, ?)", nil, form.Email, form.Pwd1, form.Nickname)
	if err != nil {
		log.Fatal("insert: ", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return int(id)

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

	// log.Println(email)
	// log.Println(exist)
	if exist == 0 {
		return false
	} else if exist == 1 {
		return true
	} else {
		log.Fatal("This email is not unique in db!")
		return true
	}

}

func checkNicknameExist(nickname string) bool {
	db, err := Connect()
	if err != nil {
		log.Fatal("connect: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select exists (select password from member where nickname=?)", nickname)
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
		log.Fatal("This nickname is not unique in db!")
		return true
	}
}
