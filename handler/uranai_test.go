package handler_test

import (
	. "."
	"testing"
	"../interfaces"
	"fmt"
)

type testPost struct {
	f func(interfaces.Message) error
}
func (p testPost) PostMessage (m interfaces.Message) error {
	return p.f(m)
}

func newPostable(f func(interfaces.Message) error) interfaces.Postable {
	t := testPost{f}
	return t
}

func TestUranai(t *testing.T) {
	var uranai UranaiHandler

	var m interfaces.Message
	m = interfaces.Message{0, "message", "general", "hello"}

	var err error
	nullPostable := newPostable( func(m interfaces.Message) error {
		t.Errorf("This must not be called.")
		return nil
	})
	// invalid command
	err = uranai.DoHandle(m, nullPostable)
	if (err == nil) {
		t.Errorf("UranaiHandler must not process a message %s", m.Text)
	}

	// invalid constellation
	m = interfaces.Message{0, "message", "general", "uranai:abc"}
	err = uranai.DoHandle(m, nullPostable)
	if (err == nil) {
		t.Errorf("UranaiHandler must not process a message %s", m.Text)
	}

	// invalid format
	m = interfaces.Message{0, "message", "general", "uranai:otome"}
	err = uranai.DoHandle(m, nullPostable)
	if (err == nil) {
		t.Errorf("UranaiHandler must not process a message %s", m.Text)
	}

	validPostable := newPostable(func(m interfaces.Message) error {
		return nil
	})

	for key, _ := range Constellations {
		// valid format
		m = interfaces.Message{0, "message", "general", fmt.Sprintf("uranai: %s", key)}
		err = uranai.DoHandle(m, validPostable)
		if (err != nil) {
			t.Errorf("UranaiHandler must not process a message %s", m.Text)
		}
	}
}