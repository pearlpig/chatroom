package main

import (
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {
	// set the server parameter
	host := "127.0.0.1"
	port := "8080"

	r := mux.NewRouter()

	// set the file server
	fs := http.FileServer(http.Dir("./public"))
	// listen js file
	r.PathPrefix("/js/").Handler(fs)
	// listen css file
	r.PathPrefix("/css/").Handler(fs)

	r.HandleFunc("/", indexHandler)
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
	fmt.Fprintln(w, "Welcome!")
}
