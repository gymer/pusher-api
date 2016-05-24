package main

import (
	"log"

	"github.com/gymer/pusher-api/controllers"
	_ "github.com/gymer/pusher-api/docs"
	"github.com/gymer/pusher-api/models"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

const defaultPort = "3001"

var (
	env, port string
	err       error
)

func init() {
	flag.StringVar(&env, "env", "dev", "set app environment")
	flag.StringVar(&port, "port", defaultPort, "listening port")
}

func main() {
	flag.Parse()

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	// r.

	// s := r.Host("www.example.com").Subrouter()
	r.HandleFunc("/ws/app/{key}", controllers.Join)

	startServer()

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}

func startServer() {
	models.ConnectDB()
	// routers.Config()
	controllers.AppStart()
}
