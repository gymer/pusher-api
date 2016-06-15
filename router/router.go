package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gymer/pusher-api/controllers"
	"github.com/urfave/negroni"
)

var serverPort string

func namespace(namespace string) func(path string) string {
	return func(path string) string {
		return namespace + path
	}
}

func Create(port string, env string) http.Handler {
	serverPort = port

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	r.Methods("GET").Path("/ws/app/{key}").HandlerFunc(controllers.Join)

	appsRouter := mux.NewRouter().PathPrefix("/v1/apps/{appId}").Subrouter()
	appsRouter.Methods("POST").Path("/events").HandlerFunc(controllers.CreateEvent)

	r.PathPrefix("/apps/{appId}").Handler(negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
		negroni.Wrap(appsRouter),
	))

	n := negroni.New(negroni.HandlerFunc(BaseMiddleware))
	if env == "dev" {
		n.Use(negroni.HandlerFunc(LoggerMiddleware))
	}

	n.UseHandler(r)

	return n
}
