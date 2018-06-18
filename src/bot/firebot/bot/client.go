package bot

import (
	"bot/client"
	"bot/config"
	"context"
	"fmt"
	"log"
	"sync"

	gomastodon "bot/go-mastodon"
)

var botClient *client.BotClient
var err error

func init() {
	var once sync.Once
	once.Do(func() {
		config := config.FireBotClientInfo()
		botClient, err = client.New(&config)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func Lauch() {
	ctx, cancel := context.WithCancel(context.Background())
	userCh, err := botClient.WS.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	for uq := range userCh {
		switch uq.(type) {
		case *gomastodon.NotificationEvent:
			e := uq.(*gomastodon.NotificationEvent)
			HandleNotification(e)
		default:
			fmt.Printf("other event: %s\n", uq)
		}
	}
}
