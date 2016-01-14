package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/gymer/pusher-api/models"
	"github.com/pborman/uuid"

	"github.com/astaxie/beego"
)

type WebsocketController struct {
	beego.Controller
}

// @Title Get
// @Description get all Users
// @Success 200 {object} models.User
// @router /app/:key [get]
func (w *WebsocketController) Connect() {
	appKey := w.Ctx.Input.Params[":key"]
	uuid := uuid.New()

	// Validate client_access_token
	app, err := w.findAppByAccessToken(appKey)

	if err != nil {
		http.Error(w.Ctx.ResponseWriter, "Wrong app key", 400)
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

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

func (w *WebsocketController) findAppByAccessToken(key string) (models.App, error) {
	var app models.App
	var err error

	err = models.DB.Where("client_access_token = ?", key).First(&app).Error

	return app, err
}
