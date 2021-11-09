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
			log.Printf("ERROR receiving message from client %v", ws)
			break
		}

		conn.chat.Lock()
		_, exists := conn.chat.Rooms[msg.Room]
		conn.chat.Unlock()

		if !exists {
			cl.SendJSON(&chat.Message{
				Action:     chat.ErrorAction,
				Room:       msg.Room,
				Username:   msg.Username,
				Content:    "Client trying to interact with non-existent room",
				HTTPstatus: http.StatusBadRequest,
			})
		} else {
			conn.processMessage(msg, cl)
		}
	}
}

func (conn *Connection) processMessage(msg chat.Message, cl *chat.Client) {
	hasJoined := cl.HasJoined(msg.Room)

	switch msg.Action {

	case chat.JoinAction:
		if !hasJoined {
			stream := chat.NewMsgStream()
			cl.RoomStreams[msg.Room] = stream
			conn.chat.Lock()
			go conn.chat.Rooms[msg.Room].ListenToStream(stream, cl)
			conn.chat.Unlock()
		}

	case chat.LeaveAction:
		if hasJoined {
			cl.RoomStreams[msg.Room].DoneSending <- true
			close(cl.RoomStreams[msg.Room].DoneSending)
			close(cl.RoomStreams[msg.Room].SentMsgs)
			delete(cl.RoomStreams, msg.Room)
		}

	case chat.TalkAction:
		if hasJoined {
			cl.RoomStreams[msg.Room].SentMsgs <- msg
		}
	}
}
