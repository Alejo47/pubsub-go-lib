package pubsub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type PubSub struct {
	TotalClients int                  `json:"totalClients"`
	TotalTopics  int                  `json:"totalTopics"`
	Clients      []*Client            `json:"clients"`
	Topics       map[string][]*Client `json:"topics"`
}

type Client struct {
	Authorized bool            `json:"authorized"`
	Closed     bool            `json:"-"`
	Id         string          `json:"id"`
	Topics     []string        `json:"topics"`
	Connection *websocket.Conn `json:"-"`
	Mutex      sync.Mutex      `json:"-"`
}

type PubSubMessage struct {
	Command string                 `json:"command,omitonempty"`
	Room    string                 `json:"room,omitonempty"`
	Message string                 `json:"data,omitonempty"`
	Params  map[string]interface{} `json:"params,omitonempty"`
}
