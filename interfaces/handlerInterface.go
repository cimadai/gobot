package interfaces

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.
type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type Postable interface {
	PostMessage(Message) error
}

type Handler interface {
	DoHandle(m Message, obj Postable) error
}


