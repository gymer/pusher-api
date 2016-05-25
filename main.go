package main

import (
	"fmt"
	"log"

	"github.com/gymer/pusher-api/controllers"
	_ "github.com/gymer/pusher-api/docs"
	"github.com/gymer/pusher-api/models"
	"github.com/gymer/pusher-api/router"

	"net/http"

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

	router := router.Create()
	// loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	startApp()
	fmt.Printf("Running on port: %+v \n", defaultPort)
	fmt.Printf("Environment: %+v \n", env)
	log.Fatal(http.ListenAndServe(":"+defaultPort, router))
}

func startApp() {
	models.ConnectDB(env)
	// routers.Config()
	controllers.AppStart()
}
