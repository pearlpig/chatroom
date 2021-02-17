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
	"github.com/gorilla/securecookie"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	hashKey    = "qtm^ac',{6t=&_mK9v%ngG/P'B?A!.(?)c^CQr-s{QZCF>Y4c=R&eyPqagBH(.9d"
	blockKey   = "m5cV7vY[4*5%7b{_"
	cookieName = "chatroom"
)

var secureC *securecookie.SecureCookie

// Cookies ...
type Cookies struct {
	MemberID int    `json:"member_id"`
	Nickname string `json:"nickname"`
	RoomID   int    `json:"room_id"`
}

// Server ...
func Server() {
	host := "127.0.0.1"
	port := "8080"

	secureC = securecookie.New([]byte(hashKey), []byte(blockKey))

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/css/").Handler(fs)

	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/", getRoomListHandler).Methods("POST")

	r.HandleFunc("/login", showLoginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/signup", showSignupHandler).Methods("GET")
	r.HandleFunc("/signup", doSignupHandler).Methods("POST")

	r.HandleFunc("/create", showCreateRoomHandler).Methods("GET")
	r.HandleFunc("/create", doCreateRoomHandler).Methods("POST")

	// s := r.PathPrefix("/room/{[0~9]+}").Subrouter()
	s := r.PathPrefix("/room/{roomID:[0-9]+}").Subrouter()
	s.HandleFunc("", chatRoomHandler).Methods("GET")

	// s.HandleFunc("/echo", wsHandler)
	s.HandleFunc("/echo", wsHandler)
	s.HandleFunc("/connRoom", connRoomHandler)
	s.HandleFunc("/disconnRoom", disconnRoomHandler)

	//
	r.HandleFunc("/check", memberAuthHandler)

	/* Create the logger for the web application. */
	l := log.New()
	// r.Use(memberAuthHandler)
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "web"))
	n.UseHandler(r)
	// Set the parameters for a HTTP server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: n,
	}

	go broker()
	go dispensor()

	log.Fatal(server.ListenAndServe())
}

func chatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/chat.html"))

	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		log.Println("Error: ", err, " roomID is not number.")
	}
	cookie := getCookie(w, r)
	cookie.RoomID = roomID
	setCookie(&w, r, cookie)
	roomName := model.GetRoomName(roomID)

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{
		Title: roomName,
	})
}
func memberAuthHandler(w http.ResponseWriter, r *http.Request) {
	auth := &Cookies{MemberID: -1, Nickname: "", RoomID: -1}
	i, err := r.Cookie(cookieName)
	if err != nil {
		log.Println("Error:", err)
		result, _ := json.Marshal(auth)
		w.Write(result)
		return
	}
	value := &Cookies{}
	if err := secureC.Decode(cookieName, i.Value, value); err != nil {
		log.Println("decode secure cookie:", err)
		result, _ := json.Marshal(auth)
		w.Write(result)
		return
	}
	if value.MemberID == -1 {
		result, _ := json.Marshal(auth)
		w.Write(result)
		return
	}
	result, _ := json.Marshal(value)
	w.Write(result)
}

func redirect(w http.ResponseWriter, target string) {
	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusFound)
}

func getCookie(w http.ResponseWriter, r *http.Request) *Cookies {
	i, err := r.Cookie(cookieName)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	value := &Cookies{}
	if err := secureC.Decode(cookieName, i.Value, value); err != nil {
		log.Println("decode secure cookie:", err)
		return nil
	}
	if err != nil {
		log.Println("Error:", err)
		return nil
	}
	return value
}
func setCookie(w *http.ResponseWriter, r *http.Request, cookie *Cookies) {
	tmp, err := secureC.Encode(cookieName, cookie)
	if err != nil {
		log.Println("encode secure cookie:", err)
	}
	c := http.Cookie{Name: cookieName, Value: tmp, MaxAge: 365 * 24 * 60 * 60, Path: "/"}
	http.SetCookie(*w, &c)
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
		setCookie(&w, r, &Cookies{MemberID: res.Data.ID, Nickname: res.Data.Nickname})
	}

	result, err := json.Marshal(res.Status)
	w.Write(result)
	redirect(w, "/")

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: cookieName, Value: "", MaxAge: -1}
	http.SetCookie(w, &cookie)
	redirect(w, "/")
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

	if err != nil {
		log.Println("Error:", err)
	}
	var form model.SignupForm
	err = schema.NewDecoder().Decode(&form, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
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
