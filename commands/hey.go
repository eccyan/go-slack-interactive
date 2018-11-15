package commands

import (
	"encoding/json"
	"fmt"
	"github.com/myrefer/go-slack-interactive/slack"
	api "github.com/nlopes/slack"
	"log"
	"net/http"
	"strings"
)

const (
	actionSelectName = "select"
	actionStartName  = "start"
	actionCancelName = "cancel"
)

type Hey struct {
	callbackID string
	mux        *slack.ServeInteractiveActionMux
}

func NewHey(callbackID string) *Hey {
	return &Hey{callbackID: callbackID, mux: NewHeyServeInteractiveActionMux(callbackID)}
}

func (hey *Hey) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	attachment := api.Attachment{
		CallbackID: hey.callbackID,
		Text:       "どの :beer: が欲しいまるか？  ",
		Color:      "#78c62c",
		Actions: []api.AttachmentAction{
			{
				Name: actionSelectName,
				Type: "select",
				Options: []api.AttachmentActionOption{
					{
						Text:  "ｱｻﾋｨｽｩﾊﾟｧﾄﾞｩﾙｧｧｧｧｲ",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "キリンラガービール",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "サッポロブラックラベル",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "モルツ",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona エール",
						Value: "Yona Yona Ale",
					},
					{
						Text:  "プレミアムモルツ",
						Value: "Plemiun Molt",
					},
					{
						Text:  "ヱビス",
						Value: "Yebisu",
					},
				},
			},

			{
				Name:  actionCancelName,
				Text:  "やっぱりやめる",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	params := api.PostMessageParameters{
		Attachments: []api.Attachment{
			attachment,
		},
	}

	if _, _, err := client.PostMessage(ev.Channel, "", params); err != nil {
		log.Printf("failed to post message: %s", err)
	}
}

// interactive actions
func NewHeyServeInteractiveActionMux(callbackID string) *slack.ServeInteractiveActionMux {
	mux := slack.NewServeInteractiveActionMux(callbackID)
	mux.Handle("select", slack.InteractiveActionHandlerFunc(actionSelect))
	mux.Handle("start", slack.InteractiveActionHandlerFunc(actionStart))
	mux.Handle("cancel", slack.InteractiveActionHandlerFunc(actionCancel))
	return mux
}

func actionSelect(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	action := callback.Actions[0]
	value := action.SelectedOptions[0].Value

	// Overwrite original drop down message.
	originalMessage := callback.OriginalMessage
	originalMessage.Attachments[0].Text = fmt.Sprintf("%s を頼んじゃうまるよ？", strings.Title(value))
	originalMessage.Attachments[0].Actions = []api.AttachmentAction{
		{
			Name:  actionStartName,
			Text:  "頼んで",
			Type:  "button",
			Value: "start",
			Style: "primary",
		},
		{
			Name:  actionCancelName,
			Text:  "やっぱダメ",
			Type:  "button",
			Style: "danger",
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&originalMessage)
}

func actionStart(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	title := "あなたの願いを受けたわまったまるよ :kira:"
	responseMessage(w, callback.OriginalMessage, title, "")
}

func actionCancel(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	title := fmt.Sprintf(" @%s 願いはキャンセルしたまる :cry:", callback.User.Name)
	responseMessage(w, callback.OriginalMessage, title, "")
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original api.Message, title, value string) {
	original.Attachments[0].Actions = []api.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []api.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
