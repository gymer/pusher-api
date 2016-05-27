package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gymer/pusher-api/controllers"
	"github.com/urfave/negroni"
)

func namespace(namespace string) func(path string) string {
	return func(path string) string {
		return namespace + path
	}
}

func Create() http.Handler {
	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	r.Methods("GET").Path("/ws/app/{key}").HandlerFunc(controllers.Join)

	appsRouter := mux.NewRouter().PathPrefix("/v1/apps/{appId}").Subrouter()
	appsRouter.Methods("POST").Path("/events").HandlerFunc(controllers.CreateEvent)

	r.PathPrefix("/apps/{appId}").Handler(negroni.New(
		negroni.HandlerFunc(AuthMiddleware),
		negroni.Wrap(appsRouter),
	))

	return r
}
