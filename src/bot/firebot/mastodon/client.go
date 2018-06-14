package mastodon

import (
	"bot/config"
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
		mci := config.GetFBotMClientInfo()
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
	ctx, cancel := context.WithCancel(context.Background())
	userq, err := wsClient.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	for {
		select {
		case uq := <-userq:
			switch uq.(type) {
			case *gomastodon.UpdateEvent:
			case *gomastodon.DeleteEvent:
			case *gomastodon.NotificationEvent:
				e := uq.(*gomastodon.NotificationEvent)
				HandleNotification(e)
			default:
				fmt.Printf("other event: %s\n", uq)
			}
		}
	}
}
