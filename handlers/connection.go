package handlers

import (
	"log"
	"net/http"

	"github.com/guillermo-st/locus/chat"

	"github.com/gorilla/websocket"
)

type Connection struct {
	chat     *chat.Chat
	upgrader websocket.Upgrader
	l        *log.Logger
}

func NewConnection(c *chat.Chat, l *log.Logger) *Connection {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Connection{c, upgrader, l}
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

	cl := chat.NewClient(ws)

	for {
		var msg chat.Message

		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("ERROR receiving message from user %v", ws)
			break
		}

		roomName := msg.Room

		conn.chat.Lock()
		_, exists := conn.chat.Rooms[roomName]
		conn.chat.Unlock()

		if exists == false {
			cl.SendJSON(&chat.Message{
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
				if !cl.HasJoined(roomName) {
					stream := chat.NewMsgStream()
					cl.RoomStreams[roomName] = stream
					conn.chat.Lock()
					go conn.chat.Rooms[roomName].ListenToStream(stream, cl)
					conn.chat.Unlock()
				}

			case chat.LeaveAction:
				if cl.HasJoined(roomName) {
					cl.RoomStreams[roomName].DoneSending <- true
					close(cl.RoomStreams[roomName].DoneSending)
					close(cl.RoomStreams[roomName].SentMsgs)
					delete(cl.RoomStreams, roomName)
				}

			case chat.TalkAction:
				if cl.HasJoined(roomName) {
					cl.RoomStreams[roomName].SentMsgs <- msg
				}
			}
		}
	}
}
