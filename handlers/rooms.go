package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/guillermo-st/locus/chat"
)

type Rooms struct {
	rooms map[string]*chat.Room
	l     *log.Logger
}

type RoomCreation struct {
	Name string `json:"name"`
}

func NewRooms(r map[string]*chat.Room, l *log.Logger) *Rooms {
	return &Rooms{r, l}
}

func (r *Rooms) GetRooms(w http.ResponseWriter, rq *http.Request) {

	rl := make([]*chat.Room, 0, len(r.rooms))

	for _, value := range r.rooms {
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

	room, exists := r.rooms[params.Name]

	if !exists {
		room = chat.NewRoom(params.Name)
		r.rooms[params.Name] = room
		go room.SendMsgs()
	}
}
