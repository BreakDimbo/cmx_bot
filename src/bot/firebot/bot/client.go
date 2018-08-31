package bot

import (
	"bot/client"
	"bot/config"
	"context"
	"log"
	"sync"

	gomastodon "bot/go-mastodon"
)

var botClient *client.Bot
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
			go HandleNotification(e)
		case *gomastodon.DeleteEvent:
			e := uq.(*gomastodon.DeleteEvent)
			go HandleDelete(e)
		default:
			// zlog.SLogger.Infof("receive other event: %s", uq)
		}
	}
}
