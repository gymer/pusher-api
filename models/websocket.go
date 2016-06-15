package models

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Event struct {
	AppId   string
	Name    string                 `json:"event"`
	Channel string                 `json:"channel"`
	Data    map[string]interface{} `json:"data"`
}

func (e *Event) GetName() string {
	return strings.Split(e.Name, "gymer:")[1]
}

type WSClient struct {
	AppID string
	Uuid  string
	Conn  *websocket.Conn
}

func (w *WSClient) Push(event Event) error {
	b, err := json.Marshal(event)

	if err != nil {
		log.Printf("error: %+v", err)
	}

	if w.Conn != nil {
		if w.Conn.WriteMessage(websocket.TextMessage, b) != nil {
			err = errors.New("Client disconnected")
		}
	}

	return err
}

type WSMessage struct {
	Client *WSClient
	Event  Event
}
