package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

var topicsMutex sync.RWMutex
var clientsMutex sync.RWMutex

func New() PubSub {
	newPS := PubSub{}

	newPS.Clients = []*Client{}
	newPS.Topics = make(map[string][]*Client)

	return newPS
}

func (ps *PubSub) AddClient(client *Client) *PubSub {
	ps.Clients = append(ps.Clients, client)
	ps.TotalClients = len(ps.Clients)
	ps.TotalTopics = len(ps.Topics)
	return ps
}

func (ps *PubSub) RemoveClient(client *Client) *PubSub {
	for _, topic := range client.Topics {
		ps.Unsubscribe(client, topic)
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if len(ps.Clients) <= 1 {
		ps.Clients = []*Client{}
		ps.TotalClients = len(ps.Clients)
		ps.TotalTopics = len(ps.Topics)
		return ps
	}

	for i, subClient := range ps.Clients {
		if client.Id == subClient.Id {
			if i+1 >= len(ps.Clients) {
				ps.Clients = ps.Clients[:i]
			} else {
				ps.Clients = append(ps.Clients[:i], ps.Clients[i+1:]...)
			}
			ps.TotalClients = len(ps.Clients)
			ps.TotalTopics = len(ps.Topics)
		}
	}
	return ps
}

func (ps *PubSub) Subscribe(client *Client, topic string) *Client {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()
	if ps.Topics == nil {
		ps.Topics = make(map[string][]*Client)
	}
	if ps.Topics[topic] != nil {
		for _, c := range ps.Topics[topic] {
			if c.Id == client.Id {
				return c
			}
		}
	}
	ps.Topics[topic] = append(ps.Topics[topic], client)
	client.Topics = append(client.Topics, topic)

	ps.TotalTopics = len(ps.Topics)
	return client
}

func (ps *PubSub) Unsubscribe(client *Client, topic string) {

	if len(client.Topics) == 1 {
		if client.Topics[0] == topic {
			client.Topics = []string{}
		}
	} else {
		for i, cTopic := range client.Topics {
			if cTopic == topic {
				if i+1 == len(client.Topics) {
					client.Topics = client.Topics[0 : i-1]
				} else {
					client.Topics = append(client.Topics[:i], client.Topics[i+1:]...)
				}
			}
		}
	}

	topicsMutex.Lock()
	defer topicsMutex.Unlock()
	if ps.Topics == nil {
		ps.Topics = make(map[string][]*Client)
	}

	for i, subClient := range ps.Topics[topic] {
		if client == subClient {
			ps.Topics[topic] = append(ps.Topics[topic][:i], ps.Topics[topic][i+1:]...)
		}
	}

	if len(ps.Topics[topic]) == 0 {
		delete(ps.Topics, topic)
	}

	ps.TotalTopics = len(ps.Topics)
	return
}

func (client *Client) SendMessage(msg []byte) {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	client.Connection.WriteMessage(websocket.TextMessage, msg)
}
