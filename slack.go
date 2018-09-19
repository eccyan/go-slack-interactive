package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/myrefer/go-slack-interactive/commands"
	"github.com/myrefer/go-slack-interactive/slack"
	api "github.com/nlopes/slack"
)

type SlackListener struct {
	client *api.Client
	botID  string
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *api.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *api.MessageEvent) error {
	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	mux := slack.NewServeMessageMux()
	mux.Handle("hey", commands.NewHey("beer"))
	mux.Handle("ping", slack.MessageHandlerFunc(commands.Ping))
	mux.Handle("recoru", slack.MessageHandlerFunc(commands.Recoru))
	mux.ServeMessage(ev, s.client)

	return nil
}
