package controllers

import (
	"container/list"
	"log"
	"net/http"

	"github.com/Iverson/pusher-api/models"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"

	"github.com/astaxie/beego"
)

var (
	// Channel for new join users.
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	publish = make(chan models.Event, 10)
	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()
)

// Operations about Users
type WebsocketController struct {
	beego.Controller
}

// @Title Get
// @Description get all Users
// @Success 200 {object} models.User
// @router /app/:key [get]
func (this *WebsocketController) Connect() {
	appKey := this.Ctx.Input.Params[":key"]
	uuid := uuid.New()

	if !this.checkAppKey(appKey) {
		http.Error(this.Ctx.ResponseWriter, "Wrong app key", 400)
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join chat room.
	log.Printf("New websocket connection: %s \n", uuid)
	defer log.Println("Close websocket connection")
	subscribers.PushBack(models.Subscriber{Uuid: uuid, Conn: ws})

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- models.Event{User: uuid, Content: string(p)}
	}
}

func (this *WebsocketController) checkAppKey(key string) bool {
	return key == "secret-app-key"
}

// @Title createUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *WebsocketController) Post() {
	log.Println(subscribers.Len())

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(models.Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, u.Ctx.Input.RequestBody) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(models.Subscriber).Uuid
			}
		}
	}
}

func init() {
	go websocketDispatcher()
}

func websocketDispatcher() {
	for {
		select {
		case unsub_uuid := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(models.Subscriber).Uuid == unsub_uuid {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(models.Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub_uuid)
					}
					// publish <- newEvent(models.EVENT_LEAVE, unsub_uuid, "")
					break
				}
			}
		}
	}
}
