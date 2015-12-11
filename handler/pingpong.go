package handler
import (
	"../interfaces"
	"strings"
	"errors"
)

// PingPongHandler is very simple ping-pong handler.
type PingPongHandler struct {}

// DoHandle handles a message.
func (h PingPongHandler) DoHandle(m interfaces.Message, obj interfaces.Postable) (err error) {
	if m.Type == "message" && strings.HasPrefix(m.Text, "ping") {
		m.Text = "pong"
		obj.PostMessage(m)
		err = nil
	} else {
		err = errors.New("Cannot parse.")
	}
	return
}
