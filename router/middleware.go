package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gymer/pusher-api/controllers"
	"github.com/gymer/pusher-api/models"
)

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
