package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	username string
	content  string
}

type Room struct {
	Name     string
	upgrader websocket.Upgrader
	msgs     chan message
	users    map[*websocket.Conn]bool
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

//the Room type satisfies the Handler interface declared in Go's http package from the stdlib.
func (room *Room) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
