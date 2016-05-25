package router

import (
	"net/http"

	"github.com/gymer/pusher-api/controllers"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var ns1 = namespace("/v1")

type Routes []Route

// var routes = Routes{
// 	Route{
// 		"WSConnect",
// 		"GET",
// 		"/ws/app/{key}",
// 		controllers.Join,
// 	},
// 	Route{
// 		"CreateEvent",
// 		"POST",
// 		"/events",
// 		controllers.CreateEvent,
// 	},
// }

func namespace(namespace string) func(path string) string {
	return func(path string) string {
		return namespace + path
	}
}

func Create() http.Handler {
	router := httprouter.New()
	router.GET(ns1("/ws/app/:key"), controllers.Join)
	router.POST(ns1("/apps/:appId/events"), controllers.CreateEvent)

	// handler := AddFilter(ns1("/apps/"), AuthMiddleware)(router)

	return AuthMiddleware(router)
}
