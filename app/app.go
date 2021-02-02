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

var userList = make(map[int]*model.Member)

// Server ...
func Server() {
	host := "127.0.0.1"
	port := "8080"

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/css/").Handler(fs)

	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/checkMember", checkMemberHandler)
	r.HandleFunc("/login", showLoginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/signup", showSignupHandler).Methods("GET")
	r.HandleFunc("/signup", doSignupHandler).Methods("POST")

	r.HandleFunc("/createroom", showCreateRoomHandler).Methods("GET")
	r.HandleFunc("/chatroom", showChatRoomHandler).Methods("GET")

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

func checkMemberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("userId")
	if err != nil {
		log.Println("Error:", err)
		userInfo, _ := json.Marshal(&model.Member{ID: -1})
		w.Write(userInfo)
		return
	}
	i, _ := strconv.Atoi(id.Value)
	userInfo, err := json.Marshal(userList[i])
	if err != nil {
		log.Println("Error:", err)
		userInfo, _ := json.Marshal(&model.Member{ID: -1})
		w.Write(userInfo)
		return
	}
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
	log.Println(form)
	checkMemberInfo := model.CheckLogin(form)
	if checkMemberInfo.Status == 0 {
		cookie := http.Cookie{Name: "userId", Value: strconv.Itoa(checkMemberInfo.MemberInfo.ID), MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &cookie)
		userList[checkMemberInfo.MemberInfo.ID] = checkMemberInfo.MemberInfo
	}

	result, err := json.Marshal(checkMemberInfo)

	w.Write(result)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := r.Cookie("userId")
	if err != nil {
		log.Println("Error:", err)
		return
	}
	i, _ := strconv.Atoi(id.Value)
	userList[i] = nil
	cookie := http.Cookie{Name: "userId", MaxAge: -1}
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
	signupInfo := model.CheckSignup(form)
	if signupInfo.Status == 0 {
		cookie := http.Cookie{Name: "userId", Value: strconv.Itoa(signupInfo.MemberInfo.ID), MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &cookie)
		userList[signupInfo.MemberInfo.ID] = signupInfo.MemberInfo
	}

	result, err := json.Marshal(signupInfo)

	w.Write(result)

}
func showCreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/create_room.html"))

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{Title: "建立聊天室"})
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
