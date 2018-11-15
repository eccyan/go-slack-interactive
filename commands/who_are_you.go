package commands

import (
	api "github.com/nlopes/slack"
	"log"
	"math/rand"
	"time"
)

func WhoAreYou(ev *api.MessageEvent, client *api.Client) {
	rand.Seed(time.Now().Unix())

	params := api.NewPostMessageParameters()
	params.LinkNames = 1
	params.EscapeText = false
	if _, _, err := client.PostMessage(ev.Channel, message(), params); err != nil {
		log.Printf("failed to post message: %s", err)
	}

}

func messages() []string {
	return []string{
		"リファ丸の名前はリファ丸まるというまる :kira: 運命石の扉（シュタインズゲート）の導きに従って、世界の支配構造を変革し、混沌を巻き起こすため研究を続ける狂気のマッドサイエンティストであるまる。",
		"性別はないまる。どっちもいけるまる",
		"リファ丸は自分のことをリファ丸って呼ぶまる",
		"ｱｻﾋｨｽｩﾊﾟｧﾄﾞｩﾙｧｧｧｧｲ",
		"リファ丸は Dr.スズキに製造された兵器まる :rocket:",
	}
}

func message() string {
	mem := messages()
	return mem[rand.Intn(len(mem))]
}
