package chat

import (
	"encoding/json"
	"sync"
)

type Room struct {
	Name    string `json:"name"`
	msgs    chan Message
	clients map[*Client]bool
	sync.Mutex
}

type Chat struct {
	sync.Mutex
	Rooms map[string]*Room
}

func NewChat() *Chat {
	return &Chat{
		Rooms: make(map[string]*Room),
	}
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		msgs:    make(chan Message),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) MarshalJSON() ([]byte, error) {
	type SimplifiedRoom Room
	return json.Marshal(&struct {
		*SimplifiedRoom
		UserCount int `json:"count"`
	}{
		SimplifiedRoom: (*SimplifiedRoom)(r),
		UserCount:      len(r.clients),
	})

}

func (room *Room) SendMsgs() {
	for {
		msg := <-room.msgs

		room.Lock()

		for client := range room.clients {

			err := client.SendJSON(&msg)
			if err != nil {
				room.unregisterUser(client)
			}
		}
		room.Unlock()
	}
}

func (room *Room) registerUser(c *Client) {
	room.Lock()
	defer room.Unlock()

	room.clients[c] = true
}

func (room *Room) unregisterUser(c *Client) {
	room.Lock()
	defer room.Unlock()

	room.clients[c] = false
	delete(room.clients, c)
}

func (room *Room) ListenToStream(s *MsgStream, c *Client) {
	room.registerUser(c)

	for {
		select {
		case msg := <-s.SentMsgs:
			room.msgs <- msg

		case <-s.DoneSending:
			room.unregisterUser(c)
			return
		}
	}
}
