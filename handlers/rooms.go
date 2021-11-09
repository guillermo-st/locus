package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/guillermo-st/locus/chat"
)

type Rooms struct {
	c *chat.Chat
	l *log.Logger
}

type RoomCreation struct {
	Name string `json:"name"`
}

func NewRooms(c *chat.Chat, l *log.Logger) *Rooms {
	return &Rooms{c, l}
}

func (r *Rooms) GetRooms(w http.ResponseWriter, rq *http.Request) {

	rl := make([]*chat.Room, 0, len(r.c.Rooms))

	r.c.Lock()
	defer r.c.Unlock()

	for _, value := range r.c.Rooms {
		rl = append(rl, value)
	}

	roomsJson, _ := json.Marshal(rl)
	w.Write(roomsJson)
}

func (r *Rooms) CreateRoom(w http.ResponseWriter, rq *http.Request) {
	var params RoomCreation

	err := json.NewDecoder(rq.Body).Decode(&params)
	if err != nil {
		r.l.Println("Invalid parameters for Room creation provided")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	r.c.Lock()
	defer r.c.Unlock()

	room, exists := r.c.Rooms[params.Name]
	if !exists {
		room = chat.NewRoom(params.Name)
		r.c.Rooms[params.Name] = room
		go room.SendMsgs()
	}
}
