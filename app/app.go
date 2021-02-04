package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"chatroom/model"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Cookies ...
type Cookies struct {
	MemberID int    `json:"member_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	RoomID   int    `json:"room_id,omitempty"`
}

// Server ...
func Server() {
	host := "127.0.0.1"
	port := "8080"

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/css/").Handler(fs)

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/", getRoomListHandler).Methods("POST")

	r.HandleFunc("/checkMember", checkCookieHandler)

	r.HandleFunc("/login", showLoginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/signup", showSignupHandler).Methods("GET")
	r.HandleFunc("/signup", doSignupHandler).Methods("POST")

	r.HandleFunc("/create", showCreateRoomHandler).Methods("GET")
	r.HandleFunc("/create", doCreateRoomHandler).Methods("POST")

	r.HandleFunc("/room/{[0~9]+}", showChatRoomHandler).Methods("GET")
	r.HandleFunc("/enterRoom", getRoomListHandler).Methods("GET")

	/* Create the logger for the web application. */
	l := log.New()

	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "web"))
	n.UseHandler(r)
	// Set the parameters for a HTTP server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: n,
	}

	log.Fatal(server.ListenAndServe())
}

func checkCookieHandler(w http.ResponseWriter, r *http.Request) {

	cookie := &Cookies{MemberID: -1, Nickname: ""}

	id, err := r.Cookie("memberID")
	if err != nil {
		log.Println("Error:", err)
		userInfo, _ := json.Marshal(cookie)
		w.Write(userInfo)
		return
	}
	i, _ := strconv.Atoi(id.Value)
	cookie.MemberID = i
	if err != nil {
		log.Println("Error:", err)
		userInfo, _ := json.Marshal(cookie)
		w.Write(userInfo)
		return
	}
	nickname, err := r.Cookie("nickname")
	if err != nil {
		log.Println("Error:", err)
		userInfo, _ := json.Marshal(cookie)
		w.Write(userInfo)
		return
	}
	cookie.Nickname = nickname.Value
	userInfo, _ := json.Marshal(cookie)
	w.Write(userInfo)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/index.html"))

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{
		"首頁",
	})
}

func showLoginHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/login.html"))

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{Title: "聊天室登入"})
}
func doLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	var form model.LoginForm
	if err != nil {
		log.Println("Error:", err)
	}
	err = schema.NewDecoder().Decode(&form, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
	res := model.CheckLogin(form)
	if res.Status.Code == 0 {
		c1 := http.Cookie{Name: "memberID", Value: strconv.Itoa(res.Data.ID), MaxAge: 365 * 24 * 60 * 60}
		c2 := http.Cookie{Name: "nickname", Value: res.Data.Nickname, MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &c1)
		http.SetCookie(w, &c2)
	}

	result, err := json.Marshal(res.Status)
	w.Write(result)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("memberID")
	if err != nil {
		log.Println("Error:", err)
		return
	}
	cookie := http.Cookie{Name: "memberID", MaxAge: -1}
	http.SetCookie(w, &cookie)
	return
}

func showSignupHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/signup.html"))
	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{Title: "聊天室註冊"})
}

func doSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	var form model.SignupForm
	if err != nil {
		log.Println("Error:", err)
	}
	err = schema.NewDecoder().Decode(&form, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println(form)
	res := model.CheckSignup(form)
	if res.Status.Code == 0 {
		c1 := http.Cookie{Name: "memberID", Value: strconv.Itoa(res.Data.ID), MaxAge: 365 * 24 * 60 * 60}
		c2 := http.Cookie{Name: "nickname", Value: res.Data.Nickname, MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &c1)
		http.SetCookie(w, &c2)
	}
	result, err := json.Marshal(res.Status)
	w.Write(result)

}

// create room
func showCreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/create_room.html"))

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{Title: "建立聊天室"})
}

func doCreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Cookie("memberID")
	i, err := strconv.Atoi(id.Value)
	if err != nil {
		log.Println("Error:", err)
	}
	var form model.CreateRoomForm
	err = r.ParseForm()

	err = schema.NewDecoder().Decode(&form, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
	form.MemberID = i
	res := model.CheckCreate(form)
	if res.Status.Code == 0 {
		c := http.Cookie{Name: "roomID", Value: strconv.Itoa(res.Data.ID), MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &c)
	}

	result, err := json.Marshal(res)
	if err != nil {
		log.Println("Error: ", err)
	}
	w.Write(result)

}

func showChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/chat.html"))
	var roomName string = "room"
	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{
		Title: roomName,
	})
}

// Show ...
type Show struct {
	Page int
}

func getRoomListHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error:", err)
	}

	var reqRoom Show
	// reqRoom := struct{ page int }{}
	err = schema.NewDecoder().Decode(&reqRoom, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
	// log.Println(reqRoom)
	roomList := model.GetRoom(reqRoom.Page)
	result, err := json.Marshal(roomList)
	if err != nil {
		log.Println("Error: ", err)
	}
	// log.Println(roomList)
	w.Write(result)
}
