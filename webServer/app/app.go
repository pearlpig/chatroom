package app

import (
	"chatroom/model"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	// hashKey    = "qtm^ac',{6t=&_mK9v%ngG/P'B?A!.(?)c^CQr-s{QZCF>Y4c=R&eyPqagBH(.9d"
	// blockKey   = "m5cV7vY[4*5%7b{_"
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
	host := getIP()
	port := "8080"

	secureC = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(16))

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/css/").Handler(fs)
	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/images/").Handler(fs)

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

func memberAuthHandler(w http.ResponseWriter, r *http.Request) {
	auth := &Cookies{MemberID: -1, Nickname: "", RoomID: -1}
	cookie := getCookie(w, r)
	if cookie == nil || cookie.MemberID == -1 {
		result, _ := json.Marshal(auth)
		w.Write(result)
		return
	}
	result, _ := json.Marshal(cookie)
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
	var tmpl = template.Must(template.ParseFiles("views/index.html"))

	tmpl.ExecuteTemplate(w, "index", struct {
		Title string
	}{
		"首頁",
	})
}

func showLoginHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/login.html"))

	tmpl.ExecuteTemplate(w, "login", struct {
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
func showSignupHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("views/signup.html"))
	tmpl.ExecuteTemplate(w, "signup", struct {
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
		setCookie(&w, r, &Cookies{MemberID: res.Data.ID, Nickname: res.Data.Nickname})
	}
	result, err := json.Marshal(res.Status)
	w.Write(result)

}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: cookieName, MaxAge: -1}
	http.SetCookie(w, &cookie)
	redirect(w, "/")
	return
}
func chatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/chat.html"))

	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		log.Println("Error: ", err, " roomID is not number.")
	}
	cookie := getCookie(w, r)
	if cookie != nil {
		cookie.RoomID = roomID
		setCookie(&w, r, cookie)
	}
	roomName := model.GetRoomName(roomID)

	tmpl.ExecuteTemplate(w, "chatroom", struct {
		Title string
	}{
		Title: roomName,
	})
}

// create room
func showCreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/create_room.html"))

	tmpl.ExecuteTemplate(w, "create", struct {
		Title string
	}{Title: "建立聊天室"})
}

func doCreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	cookie := getCookie(w, r)

	var form model.CreateRoomForm
	err := r.ParseForm()
	if err != nil {
		log.Println("Error:", err)
	}
	err = schema.NewDecoder().Decode(&form, r.PostForm)
	if err != nil {
		log.Println("Error:", err)
	}
	form.MemberID = cookie.MemberID
	res := model.CheckCreate(form)
	log.Println(res)
	if res.Status.Code == 0 {
		setCookie(&w, r, &Cookies{RoomID: res.Data.ID, MemberID: cookie.MemberID, Nickname: cookie.Nickname})
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

func getIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}
