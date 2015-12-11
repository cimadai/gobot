package main

import (
	"fmt"
	"log"
	"os"
	"./slack"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: gobot slack-bot-token\n")
		os.Exit(1)
	}

	var slack slack.Slack
	// start a websocket-based Real Time API session
	id, err := slack.Connect(os.Args[1])
	if err != nil {
		// cannot connect slack api server.
		log.Fatal(err)
		os.Exit(1)
	}

	commandLineResp := fmt.Sprintf("gobot (%s) ready, ^C exits", id)
	fmt.Println(commandLineResp)

	for {
		// read each incoming message
		m, err := slack.GetMessage()
		if err != nil {
			log.Fatal(err)
		}

		slack.DoHandle(m)
	}
}

