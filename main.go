package main

import (
	"fmt"
	"html/template"

	"net/http"

	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Server ...
func main() {
	host := "127.0.0.1"
	port := "8081"

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.PathPrefix("/js/").Handler(fs)
	r.PathPrefix("/css/").Handler(fs)
	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/login", showLoginHandler).Methods("GET")

	r.HandleFunc("/signup", showSignupHandler).Methods("GET")
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

func showSignupHandler(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFiles("views/layout.html", "views/head.html", "views/index.html")
	var tmpl = template.Must(template.ParseFiles("views/template.html", "views/signup.html"))

	tmpl.ExecuteTemplate(w, "template", struct {
		Title string
	}{Title: "聊天室註冊"})
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
