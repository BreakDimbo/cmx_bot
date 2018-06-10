package mastodon

import (
	"bot/ai_x/config"
	"context"
	"fmt"
	"log"
	"sync"

	gomastodon "bot/go-mastodon"
)

var wsClient *gomastodon.WSClient
var client *gomastodon.Client

func init() {
	var once sync.Once
	once.Do(func() {
		mci := config.GetMastodonClientInfo()
		c := gomastodon.NewClient(&gomastodon.Config{
			Server:       mci.Sever,
			ClientID:     mci.ID,
			ClientSecret: mci.Secret,
		})
		err := c.Authenticate(context.Background(), mci.Email, mci.Password)
		if err != nil {
			log.Fatal(err)
		}

		client = c
		wsClient = c.NewWSClient()
	})
}

func Lauch() {
	RunPoster(client)

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
