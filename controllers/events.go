package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gymer/pusher-api/models"
	"github.com/julienschmidt/httprouter"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var event models.Event
	app := findOrAddApp(ps.ByName("appId"))

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
