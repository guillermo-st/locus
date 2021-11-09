package chat

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sync.Mutex
	RoomStreams map[string]*MsgStream
	ws          *websocket.Conn
}

func NewClient(w *websocket.Conn) *Client {
	return &Client{
		RoomStreams: make(map[string]*MsgStream),
		ws:          w,
	}
}

func (cl *Client) SendJSON(m *Message) error {
	cl.Lock()
	defer cl.Unlock()

	err := cl.ws.WriteJSON(m)
	if err != nil {
		log.Printf("ERROR sending message to user %v: %v", cl.ws, err)
	}
	return nil
}

func (cl *Client) HasJoined(room string) bool {
	_, hasJoined := cl.RoomStreams[room]
	return hasJoined
}
