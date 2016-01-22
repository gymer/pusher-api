package controllers

import (
	"container/list"
	"log"

	"github.com/astaxie/beego/logs"

	"github.com/gymer/pusher-api/models"
)

type DataStore struct {
	Apps map[string]*models.App
}

var (
	Logger       *logs.BeeLogger
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

func newEvent(name string, channel string, data map[string]string) models.Event {
	if data == nil {
		data = make(map[string]string)
	}

	return models.Event{Name: name, Channel: channel, Data: data}
}

func newServiceEvent(name string, channel string, data map[string]string) models.Event {
	return newEvent("gymmer:"+name, channel, data)
}

func broadcastEvent(app *models.App, event models.Event) int {
	pushCount := 0

	log.Printf("Subscribed: %+v", app.Subscriptions[event.Channel])

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

func init() {
	Logger = logs.NewLogger(10000)
	Logger.SetLogger("console", "")

	go appDispatcher()
}

func findApp(appID string) *models.App {
	return store.Apps[appID]
}

func findOrAddApp(appID string) *models.App {
	if findApp(appID) == nil {
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
			data := map[string]string{
				"socket_id": connect_client.Uuid,
			}
			connect_client.Push(newServiceEvent("connection_established", "", data))

		case disconnect_client := <-wsDisconnect:
			var app *models.App = findOrAddApp(disconnect_client.AppID)

			if disconnect_client.Conn != nil {
				disconnect_client.Conn.Close()
				Logger.Warn("WebSocket closed: %s", disconnect_client.Uuid)
			}
			app.RemoveClient(disconnect_client)

		case message := <-wsMessage:
			client := message.Client
			app := getApp(client.AppID)
			eventName := message.Event.GetName()
			Logger.Info("App before: %+v", app)

			switch eventName {
			case "subscribe":
				app.SubscribeToChannel(client, message.Event.Channel)
				client.Push(newServiceEvent("subscription_success", message.Event.Channel, nil))
			case "unsubscribe":
				app.UnsubscribeToChannel(client, message.Event.Channel)
				client.Push(newServiceEvent("unsubscription_success", message.Event.Channel, nil))
			default:
				Logger.Info("Unknown event type: %+v", eventName)

			}

			Logger.Info("App after: %+v", app)
		}

	}
}
