package models

import "container/list"

type App struct {
	ID      string
	Clients *list.List
}

func (a *App) AddClient(client WSClient) {
	a.Clients.PushBack(client)
}
