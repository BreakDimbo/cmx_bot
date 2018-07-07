package bot

import (
	"bot/client"
	"bot/config"
	"bot/intelbot/const"
	zlog "bot/log"
	"context"
	"log"
	"os"
	"sync"

	gomastodon "bot/go-mastodon"
)

var botClient *client.BotClient
var err error

func init() {
	var once sync.Once
	once.Do(func() {
		config := config.IntelBotClientInfo()
		botClient, err = client.New(&config)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func Lauch() {
	ctx, cancel := context.WithCancel(context.Background())
	publicCh, err := botClient.WS.StreamingWSPublic(ctx, true)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	userCh, err := botClient.WS.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	for {
		select {
		case event := <-userCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				e := event.(*gomastodon.UpdateEvent)
				HandleUpdate(e, con.ScopeTypeLocal)
			case *gomastodon.DeleteEvent:
				e := event.(*gomastodon.DeleteEvent)
				HandleDelete(e, con.ScopeTypeLocal)
			case *gomastodon.NotificationEvent:
				e := event.(*gomastodon.NotificationEvent)
				HandleNotification(e)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}

		case event := <-publicCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				e := event.(*gomastodon.UpdateEvent)
				HandleUpdate(e, con.ScopeTypePublic)
			case *gomastodon.DeleteEvent:
				e := event.(*gomastodon.DeleteEvent)
				HandleDelete(e, con.ScopeTypePublic)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}
		}
	}
}
