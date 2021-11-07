package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

type Chat struct {
	Rooms map[string]*Room
	l     *log.Logger
}

func NewChat(l *log.Logger) *Chat {

	return &Chat{make(map[string]*Room),
		l,
	}
}

//the Chat type satisfies the Handler interface declared in Go's http package from the stdlib.
func (chat *Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomName := mux.Vars(r)["roomName"]

	room, exists := chat.Rooms[roomName]

	if !exists {
		room = NewRoom(roomName)
		go room.SendMsgs()
	}
	room.RegisterUser(w, r)
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
