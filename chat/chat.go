package chat

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	TalkAction  = "talk"
	JoinAction  = "join"
	LeaveAction = "leave"
)

type Message struct {
	Action   string
	Room     string
	username string
	content  string
}

type MsgStream struct {
	SentMsgs    chan Message
	DoneSending chan bool
}

type Client struct {
	RoomStreams map[string]*MsgStream
}

type Room struct {
	Name  string `json:"name"`
	msgs  chan Message
	users map[*websocket.Conn]bool
	sync.Mutex
}

func NewMsgStream() *MsgStream {
	return &MsgStream{
		make(chan Message),
		make(chan bool),
	}
}

func NewClient() *Client {
	return &Client{make(map[string]*MsgStream)}
}

func NewRoom(name string) *Room {
	return &Room{
		Name:  name,
		msgs:  make(chan Message),
		users: make(map[*websocket.Conn]bool),
	}
}

func (r *Room) MarshalJSON() ([]byte, error) {
	type SimplifiedRoom Room
	return json.Marshal(&struct {
		*SimplifiedRoom
		UserCount int `json:"count"`
	}{
		SimplifiedRoom: (*SimplifiedRoom)(r),
		UserCount:      len(r.users),
	})

}

func (room *Room) SendMsgs() {

	for {
		msg := <-room.msgs

		for user := range room.users {

			err := user.WriteJSON(msg)
			if err != nil {
				log.Printf("ERROR sending message to user %v: %v", user, err)
				room.unregisterUser(user)
			}
		}
	}
}

func (room *Room) registerUser(ws *websocket.Conn) {
	room.Lock()
	defer room.Unlock()

	room.users[ws] = true
}

func (room *Room) unregisterUser(ws *websocket.Conn) {
	room.Lock()
	defer room.Unlock()

	room.users[ws] = false
	delete(room.users, ws)
}

func (room *Room) ListenToStream(s *MsgStream, ws *websocket.Conn) {
	room.registerUser(ws)

	for {
		select {
		case msg := <-s.SentMsgs:
			room.msgs <- msg

		case <-s.DoneSending:
			room.unregisterUser(ws)
			return
		}
	}
}
