package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

var subsMutex sync.RWMutex
var clientsMutex sync.RWMutex

func New() PubSub {
	newPS := PubSub{}

	newPS.Clients = []*Client{}
	newPS.Subscriptions = make(map[string][]*Client)

	return newPS
}

func (ps *PubSub) AddClient(client *Client) *PubSub {
	ps.Clients = append(ps.Clients, client)
	ps.Total = len(ps.Clients)
	return ps
}

func (ps *PubSub) RemoveClient(client *Client) *PubSub {
	for _, topic := range client.Topics {
		ps.Unsubscribe(client, topic)
	}

	clientsMutex.Lock()
	if len(ps.Clients) <= 1 {
		ps.Clients = []*Client{}
		ps.Total = len(ps.Clients)
		return ps
	}

	for i, subClient := range ps.Clients {
		if client.Id == subClient.Id {
			if i+1 >= len(ps.Clients) {
				ps.Clients = ps.Clients[:i]
			} else {
				ps.Clients = append(ps.Clients[:i], ps.Clients[i+1:]...)
			}
			ps.Total = len(ps.Clients)
		}
	}
	clientsMutex.Unlock()
	return ps
}

func (ps *PubSub) Subscribe(client *Client, topic string) Subscription {
	subsMutex.Lock()
	if ps.Subscriptions == nil {
		ps.Subscriptions = make(map[string][]*Client)
	}
	if ps.Subscriptions[topic] != nil {
		for _, c := range ps.Subscriptions[topic] {
			if c == client {
				return Subscription{Topic: topic, Client: client}
			}
		}
	}
	ps.Subscriptions[topic] = append(ps.Subscriptions[topic], client)
	subsMutex.Unlock()
	client.Topics = append(client.Topics, topic)

	return Subscription{Topic: topic, Client: client}
}

func (ps *PubSub) Unsubscribe(client *Client, topic string) Subscription {
	clientsMutex.Lock()
	for i, cTopic := range client.Topics {
		if cTopic == topic {
			client.Topics = append(client.Topics[:i], client.Topics[i+1:]...)
		}
	}
	clientsMutex.Unlock()

	subsMutex.Lock()
	if ps.Subscriptions == nil {
		ps.Subscriptions = make(map[string][]*Client)
	}

	for i, subClient := range ps.Subscriptions[topic] {
		if client == subClient {
			ps.Subscriptions[topic] = append(ps.Subscriptions[topic][:i], ps.Subscriptions[topic][i+1:]...)
		}
	}

	if len(ps.Subscriptions[topic]) == 0 {
		delete(ps.Subscriptions, topic)
	}
	subsMutex.Unlock()

	return Subscription{}
}

func (client *Client) SendMessage(msg []byte) {
	client.Mutex.Lock()
	client.Connection.WriteMessage(websocket.TextMessage, msg)
	client.Mutex.Unlock()
}
