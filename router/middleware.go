package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gymer/pusher-api/controllers"
	"github.com/gymer/pusher-api/models"
)

const Body int = 0

func BaseMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("Request body error: %s \n", err)
	}

	context.Set(r, Body, body)

	next(w, r)

	context.Clear(r)
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var app models.App
	vars := mux.Vars(r)
	appId := vars["appId"]
	u, p, ok := r.BasicAuth()

	if !ok {
		controllers.HttpResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := models.DB.Where("id = ? and client_access_token = ? and server_access_token = ?", appId, u, p).First(&app).Error
	if err != nil {
		controllers.HttpResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	next(w, r)
}

func LoggerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	body := context.Get(r, Body)

	fmt.Printf("Instance: %+v \n", serverPort)
	fmt.Printf("Started %s %s at %s \n", r.Method, r.URL.Path, start)

	fmt.Printf("Body: %s \n", body)

	next(w, r)

	fmt.Printf("Finished: %s \n\n\n", time.Since(start))
}
