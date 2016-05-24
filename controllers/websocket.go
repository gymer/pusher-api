package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/gymer/pusher-api/models"
	"github.com/pborman/uuid"
)

func Join(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appKey := vars["key"]

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Cannot setup WebSocket connection: %s", err), 400)
		return
	}

	// Validate client_access_token
	app, err := findAppByAccessToken(appKey)
	if err != nil {
		closeWS(ws, InvalidAppCode, "Invalid app code")
		return
	}

	uuid := uuid.New()
	client := models.WSClient{Uuid: uuid, Conn: ws, AppID: app.ID}
	Logger.Warn("New websocket connection: %s \n", uuid)

	connect(&client)
	defer disconnect(&client)

	// Message receive loop.
	for {
		var e models.Event
		message := models.WSMessage{Client: &client}

		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}

		json.Unmarshal(p, &e)
		message.Event = e
		wsMessage <- message
	}
}

func findAppByAccessToken(key string) (models.App, error) {
	var app models.App
	var err error

	err = models.DB.Where("client_access_token = ?", key).First(&app).Error

	return app, err
}
