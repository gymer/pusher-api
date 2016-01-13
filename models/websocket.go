package models

import "github.com/gorilla/websocket"

type Event struct {
	User    string
	Content string
}

type Message struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}

type WSClient struct {
	AppID string
	Uuid  string
	Conn  *websocket.Conn
}
