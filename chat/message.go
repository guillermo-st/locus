package chat

const (
	TalkAction  = "talk"
	JoinAction  = "join"
	LeaveAction = "leave"
	ErrorAction = "error"
)

type Message struct {
	Action     string
	Room       string
	Username   string
	Content    string
	HTTPstatus int
}

type MsgStream struct {
	SentMsgs    chan Message
	DoneSending chan bool
}

func NewMsgStream() *MsgStream {
	return &MsgStream{
		make(chan Message),
		make(chan bool),
	}
}
