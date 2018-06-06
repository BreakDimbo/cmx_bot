package mastodon

import (
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"
	"sync"

	gomastodon "bot/go-mastodon"
)

var once sync.Once
var wsClient *gomastodon.WSClient
var client *gomastodon.Client

func InitOnce() {
	once.Do(func() {
		c := gomastodon.NewClient(&gomastodon.Config{
			Server:       con.Server,
			ClientID:     con.ClientId,
			ClientSecret: con.ClientSecret,
		})
		err := c.Authenticate(context.Background(), con.ClientEmail, con.ClientPassword)
		if err != nil {
			log.Fatal(err)
		}

		client = c
		wsClient = c.NewWSClient()
	})
}

func Lauch() {
	go RunPoster(client)

	ctx, cancel := context.WithCancel(context.Background())
	q, err := wsClient.StreamingWSPublic(ctx, true)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	for event := range q {
		switch event.(type) {
		case *gomastodon.UpdateEvent:
			e := event.(*gomastodon.UpdateEvent)
			HandleUpdate(e)
		case *gomastodon.DeleteEvent:
			e := event.(*gomastodon.DeleteEvent)
			HandleDelete(e)
		case *gomastodon.NotificationEvent:
			e := event.(*gomastodon.NotificationEvent)
			fmt.Println(e.Notification.Type)
		default:
			fmt.Println(event)
		}
	}
}
