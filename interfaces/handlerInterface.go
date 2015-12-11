package interfaces

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// This can post a message to slack.
type Postable interface {
	PostMessage(Message) error
}

// This can handle a message from slack.
type Handler interface {
	DoHandle(m Message, obj Postable) error
}


