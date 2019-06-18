package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type PubSub struct {
	Total         int                  `json:"total"`
	Clients       []*Client            `json:"clients"`
	Subscriptions map[string][]*Client `json:"subscriptions"`
}

type Client struct {
	Authorized bool            `json:"authorized"`
	Closed     bool            `json:"-"`
	Id         string          `json:"id"`
	Topics     []string        `json:"topics"`
	Connection *websocket.Conn `json:"-"`
	Mutex      sync.Mutex      `json:"-"`
}

type Subscription struct {
	Topic  string  `json:"topic"`
	Client *Client `json:"client"`
}

type PubSubMessage struct {
	Command string            `json:"command,omitonempty"`
	Room    string            `json:"room,omitonempty"`
	Message string            `json:"data,omitonempty"`
	Params  map[string]string `json:"params,omitonempty"`
}
