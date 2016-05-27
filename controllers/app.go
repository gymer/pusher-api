package controllers

import (
	"container/list"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gymer/pusher-api/models"
)

// Close codes
const (
	InvalidAppCode = 4001
)

type DataStore struct {
	Apps map[string]*models.App
}

var (
	store        = DataStore{Apps: make(map[string]*models.App)}
	wsConnect    = make(chan *models.WSClient)
	wsDisconnect = make(chan *models.WSClient)
	wsMessage    = make(chan models.WSMessage, 10)
)

func connect(client *models.WSClient) {
	wsConnect <- client
}

func disconnect(client *models.WSClient) {
	wsDisconnect <- client
}

func closeWS(ws *websocket.Conn, code int, reason string) {
	ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, reason), time.Now().Add(time.Second))
	ws.Close()
}

func newEvent(name string, channel string, data map[string]interface{}) models.Event {
	if data == nil {
		data = make(map[string]interface{})
	}

	return models.Event{Name: name, Channel: channel, Data: data}
}

func newServiceEvent(name string, channel string, data map[string]interface{}) models.Event {
	return newEvent("gymer:"+name, channel, data)
}

func broadcastEvent(app *models.App, event models.Event) int {
	pushCount := 0

	if app.Subscriptions[event.Channel] == nil {
		return pushCount
	}

	for client, subscribed := range app.Subscriptions[event.Channel] {
		if subscribed {
			if err := client.Push(event); err == nil {
				pushCount += 1
			} else {
				disconnect(client)
			}
		}
	}

	return pushCount
}

func AppStart() {
	go appDispatcher()
}

func findOrAddApp(appID string) *models.App {
	if getApp(appID) == nil {
		store.Apps[appID] = &models.App{ID: appID, Clients: list.New(), Subscriptions: make(map[string]map[*models.WSClient]bool)}
	}

	return store.Apps[appID]
}

func getApp(appID string) *models.App {
	return store.Apps[appID]
}

func appDispatcher() {
	for {
		select {
		case connect_client := <-wsConnect:
			var app *models.App = findOrAddApp(connect_client.AppID)

			app.AddClient(connect_client)
			data := map[string]interface{}{
				"socket_id": connect_client.Uuid,
			}
			connect_client.Push(newServiceEvent("connection_established", "", data))

		case disconnect_client := <-wsDisconnect:
			var app *models.App = findOrAddApp(disconnect_client.AppID)

			if disconnect_client.Conn != nil {
				disconnect_client.Conn.Close()
				log.Printf("WebSocket closed: %s \n", disconnect_client.Uuid)
			}
			app.RemoveClient(disconnect_client)

		case message := <-wsMessage:
			client := message.Client
			app := getApp(client.AppID)
			eventName := message.Event.GetName()
			// Logger.Info("App before: %+v", app)

			switch eventName {
			case "subscribe":
				app.SubscribeToChannel(client, message.Event.Channel)
				client.Push(newServiceEvent("subscription_success", message.Event.Channel, nil))
			case "unsubscribe":
				app.UnsubscribeToChannel(client, message.Event.Channel)
				client.Push(newServiceEvent("unsubscription_success", message.Event.Channel, nil))
			default:
				log.Printf("Unknown event type: %+v", eventName)

			}

			// Logger.Info("App after: %+v", app)
		}

	}
}
