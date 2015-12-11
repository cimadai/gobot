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
	response, err := h.process(m)
	if err != nil {
		return
	}

	m.Text = response
	obj.PostMessage(m)
	return
}

func (h PingPongHandler) process(m interfaces.Message) (response string, err error) {
	if m.Type == "message" && strings.HasPrefix(m.Text, "ping") {
		response = "pong"
		err = nil
	} else {
		err = errors.New("Cannot parse.")
	}
	return
}
