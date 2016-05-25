package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gymer/pusher-api/models"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	vars := mux.Vars(r)
	app := findOrAddApp(vars["appId"])

	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	pushedClient := broadcastEvent(app, event)
	resp := make(map[string]interface{})
	resp["pushed_clients"] = pushedClient

	httpResponseJson(w, 200, resp)
}
