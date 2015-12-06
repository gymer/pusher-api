package models

import "github.com/gorilla/websocket"

type Event struct {
	User    string
	Content string
}

type Message struct {
	event   string
	message string
}

type WSClient struct {
	AppID string
	Uuid  string
	Conn  *websocket.Conn
}
