package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	username string
	content  string
}

type Room struct {
	Name     string `json:"name"`
	upgrader websocket.Upgrader
	msgs     chan message
	users    map[*websocket.Conn]bool
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

func NewRoom(name string) *Room {
	return &Room{name,
		websocket.Upgrader{},
		make(chan message),
		make(map[*websocket.Conn]bool),
	}
}

func (room *Room) SendMsgs() {

	for {
		msg := <-room.msgs

		for user := range room.users {

			err := user.WriteJSON(msg)
			if err != nil {
				log.Printf("ERROR sending message to user %v: %v", user, err)
				user.Close()
				delete(room.users, user)
			}
		}
	}
}

func (room *Room) RegisterUser(w http.ResponseWriter, r *http.Request) {

	ws, err := room.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	room.users[ws] = true

	for {
		var msg message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ERROR receiving message from user %v", ws)
			//TODO: write error to client
		}

		room.msgs <- msg
	}

}
