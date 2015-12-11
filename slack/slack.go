package slack

/*
// original is https://github.com/rapidloop/mybot/blob/master/slack.go

Copyright (c) 2015 daisuke shimada.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"

	"../handler"
	"../interfaces"
)

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.
type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	URL   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	ID string `json:"id"`
}

// Slack facade object for using slack api.
type Slack struct {
	counter uint64
	webSocketConnection *websocket.Conn
	handlers []interfaces.Handler
}

// Connect starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func (slack *Slack) Connect(token string) (id string, err error) {
	wsurl, id, err := slack.start(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	slack.webSocketConnection = ws
	slack.handlers = handler.LoadHandlers()
	return
}

// GetMessage await for a message from slack.
func (slack *Slack) GetMessage() (m interfaces.Message, err error) {
	err = websocket.JSON.Receive(slack.webSocketConnection, &m)
	return
}

// PostMessage posts a message to slack.
func (slack *Slack) PostMessage(m interfaces.Message) error {
	m.ID = atomic.AddUint64(&slack.counter, 1)
	return websocket.JSON.Send(slack.webSocketConnection, m)
}

// DoHandle handles a message on each Slack#handlers.
func (slack *Slack) DoHandle(m interfaces.Message) error {
	// ignore error
	// TODO: Fix me
	for _, handler := range slack.handlers {
		handler.DoHandle(m, slack)
	}
	return nil
}

// start does a rtm.start, and returns a websocket URL and user ID.
// The websocket URL can be used to initiate an RTM session.
func (slack *Slack) start(token string) (wsURL string, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsURL = respObj.URL
	id = respObj.Self.ID
	return
}
