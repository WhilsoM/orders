package transport

import "github.com/coder/websocket"

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Client struct {
	ID   int
	Role string
	Conn *websocket.Conn
}

type Chat struct {
	ID      int
	User    *Client
	Support *Client
}
