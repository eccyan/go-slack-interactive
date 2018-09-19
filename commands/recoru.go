package commands

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/kelseyhightower/envconfig"
	api "github.com/nlopes/slack"
)

type envConfig struct {
	// authId=$RECORU_AUTH_ID&password=$RECORU_PASSWORD
	RecoruWorkPlaceID string `envconfig:"RECORU_WORK_PLACE_ID" required:"true"`
	RecoruAuthID      string `envconfig:"RECORU_AUTH_ID" required:"true"`
	RecoruPassword    string `envconfig:"RECORU_PASSWORD" required:"true"`
}

func Recoru(ev *api.MessageEvent, client *api.Client) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return
	}
	values := url.Values{}
	values.Add("contractId", env.RecoruWorkPlaceID)
	values.Add("authId", env.RecoruAuthID)
	values.Add("password", env.RecoruPassword)

	resp, err := login(env)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Println("[ERROR] Failed to read the body", err)
	}

	log.Printf("[INFO] Body: %s", squish(string(body[:])))
}

func login(env envConfig) (*http.Response, error) {
	values := url.Values{}
	values.Add("contractId", env.RecoruWorkPlaceID)
	values.Add("authId", env.RecoruAuthID)
	values.Add("password", env.RecoruPassword)

	resp, err := http.PostForm("https://app.recoru.in/ap/login", values)
	if err != nil {
		log.Printf("[ERROR] Failed to post var: %s", err)
		return nil, err
	}

	log.Printf("[INFO] Status: %s", resp.Status)

	return resp, nil
}

func squish(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}