package handlers

import (
	"log"
	"net/http"

	"github.com/guillermo-st/locus/chat"

	"github.com/gorilla/websocket"
)

type Connection struct {
	rooms    map[string]*chat.Room
	upgrader websocket.Upgrader
	l        *log.Logger
}

func NewConnection(rooms map[string]*chat.Room, l *log.Logger) *Connection {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Connection{rooms, upgrader, l}
}

func CheckChatOrigin(r *http.Request) bool {
	//return r.Header.Get("Origin") == "http://127.0.0.1:8001"
	return true
}

func (conn *Connection) ListenToWs(w http.ResponseWriter, r *http.Request) {

	conn.upgrader.CheckOrigin = CheckChatOrigin

	ws, err := conn.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	c := chat.NewClient(ws)

	for {
		var msg chat.Message

		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("ERROR receiving message from user %v", ws)
			break
		}

		roomName := msg.Room
		_, exists := conn.rooms[roomName]

		if exists == false {
			c.SendJSON(&chat.Message{
				Action:     chat.ErrorAction,
				Room:       roomName,
				Username:   msg.Username,
				Content:    "User trying to join non-existent room",
				HTTPstatus: http.StatusBadRequest,
			})
		}

		if exists {
			switch msg.Action {

			case chat.JoinAction:
				if !c.HasJoined(roomName) {
					stream := chat.NewMsgStream()
					c.RoomStreams[roomName] = stream
					go conn.rooms[roomName].ListenToStream(stream, c)
				}

			case chat.LeaveAction:
				if c.HasJoined(roomName) {
					c.RoomStreams[roomName].DoneSending <- true
					close(c.RoomStreams[roomName].DoneSending)
					close(c.RoomStreams[roomName].SentMsgs)
					delete(c.RoomStreams, roomName)
				}

			case chat.TalkAction:
				if c.HasJoined(roomName) {
					c.RoomStreams[roomName].SentMsgs <- msg
				}
			}
		}
	}
}
