package models

import "container/list"

type App struct {
	ID            string
	Clients       *list.List
	Subscriptions map[string]map[*WSClient]bool
}

func (a *App) AddClient(client *WSClient) {
	a.Clients.PushBack(client)
}

func (a *App) RemoveClient(client *WSClient) {
	a.UnsubscribeAll(client)

	for c := a.Clients.Front(); c != nil; c = c.Next() {
		if c.Value.(*WSClient).Uuid == client.Uuid {
			a.Clients.Remove(c)
			break
		}
	}
}

func (a *App) ChannelSubscribers(name string) map[*WSClient]bool {
	return a.Subscriptions[name]
}

func (a *App) SubscribeToChannel(client *WSClient, channel string) {
	if a.Subscriptions[channel] == nil {
		a.Subscriptions[channel] = make(map[*WSClient]bool)
	}
	a.Subscriptions[channel][client] = true
}

func (a *App) UnsubscribeToChannel(client *WSClient, channel string) {
	if a.Subscriptions[channel] != nil {
		delete(a.Subscriptions[channel], client)
	}
}

func (a *App) UnsubscribeAll(client *WSClient) {
	for _, subscriptions := range a.Subscriptions {
		delete(subscriptions, client)
	}
}
