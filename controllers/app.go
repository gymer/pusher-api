package controllers

import (
	"container/list"
	"log"

	"github.com/astaxie/beego"

	"github.com/gymer/pusher-api/models"
)

type DataStore struct {
	Apps map[string]*models.App
}

var (
	store        = DataStore{Apps: make(map[string]*models.App)}
	wsConnect    = make(chan models.WSClient)
	wsDisconnect = make(chan models.WSClient)
	publish      = make(chan models.Event, 10)
)

func Connect(client models.WSClient) {
	wsConnect <- client
}

func Disconnect(client models.WSClient) {
	wsDisconnect <- client
}

func init() {
	go storeDispatcher()
}

func findOrAddApp(appID string) *models.App {
	if store.Apps[appID] == nil {
		store.Apps[appID] = &models.App{ID: appID, Clients: list.New()}
	}

	return store.Apps[appID]
}

func getApp(appID string) *models.App {
	return store.Apps[appID]
}

func storeDispatcher() {
	for {
		select {
		case connect_client := <-wsConnect:
			var app *models.App = findOrAddApp(connect_client.AppID)
			log.Printf("%+v", app)
			app.AddClient(connect_client)
		case disconnect_client := <-wsDisconnect:
			var app *models.App = findOrAddApp(disconnect_client.AppID)

			for client := app.Clients.Front(); client != nil; client = client.Next() {
				if client.Value.(models.WSClient).Uuid == disconnect_client.Uuid {
					app.Clients.Remove(client)
					// Clone connection.
					ws := client.Value.(models.WSClient).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", disconnect_client.Uuid)
					}
					// publish <- newEvent(models.EVENT_LEAVE, disconnect_uuid, "")
					break
				}
			}
		}
	}
}
