package slack

import (
	api "github.com/nlopes/slack"
	"log"
	"strings"
)

// A MessageHandler responds to an Message
type MessageHandler interface {
	ServeMessage(*api.MessageEvent, *api.Client)
}

type ServeMessageMux struct {
	m map[string]messageMuxEntry
}

type messageMuxEntry struct {
	h       MessageHandler
	pattern string
}

func NewServeMessageMux() *ServeMessageMux { return new(ServeMessageMux) }

var DefaultServeMessageMux = &defaultServeMessageMux

var defaultServeMessageMux ServeMessageMux

func (mux *ServeMessageMux) match(text string) (h MessageHandler, pattern string) {
	// Parse message
	cmd := strings.Split(strings.TrimSpace(text), " ")[1:]
	if len(cmd) == 0 {
		return
	}

	v, ok := mux.m[cmd[0]]
	if ok {
		return v.h, v.pattern
	}

	return
}

func (mux *ServeMessageMux) MessageHandler(ev *api.MessageEvent, client *api.Client) (h MessageHandler, pattern string) {
	log.Printf("[INFO] message is %s", ev.Msg.Text)
	h, pattern = mux.match(ev.Msg.Text)

	if h == nil {
		log.Printf("[INFO] not found pattern for %s", ev.Msg.Text)
		h, pattern = CommandNotFoundHandler(), ""
	}

	return
}

func (mux *ServeMessageMux) Handle(pattern string, handler MessageHandler) {
	if pattern == "" {
		panic("slack: invalid pattern")
	}
	if handler == nil {
		panic("slack: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("slack: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]messageMuxEntry)
	}
	mux.m[pattern] = messageMuxEntry{h: handler, pattern: pattern}
}

func (mux *ServeMessageMux) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	h, _ := mux.MessageHandler(ev, client)
	h.ServeMessage(ev, client)
}

type MessageHandlerFunc func(ev *api.MessageEvent, client *api.Client)

// ServeMessage calls f(w, r).
func (f MessageHandlerFunc) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	f(ev, client)
}

func CommandNotFound(ev *api.MessageEvent, client *api.Client) {
	params := api.NewPostMessageParameters()
	message := `
そのコマンドはわからないまる :innocent:
したのコマンド一覧から選んで欲しいまる！

# コマンド一覧
- hey : ビールの願いを聞き受けるまる。
- ping : pong と返事するまる。
- assign [URL/文章] : 指定したURLや文章をエンジニアにアサインするまる。
	`
	if _, _, err := client.PostMessage(ev.Channel, message, params); err != nil {
		log.Printf("[ERROR] failed to post message: %s", err)
	}
}

func CommandNotFoundHandler() MessageHandler {
	return MessageHandlerFunc(CommandNotFound)
}
