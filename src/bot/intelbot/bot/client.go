package bot

import (
	"bot/client"
	"bot/config"
	con "bot/intelbot/const"
	zlog "bot/log"
	"context"
	"log"
	"os"
	"sync"
	"time"

	gomastodon "bot/go-mastodon"
)

var botClient *client.Bot
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

func Launch() {
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

	ntfCh := make(chan struct{})

	go func() {
		timer := time.NewTimer(65 * time.Minute)

		for {
			select {
			case <-ntfCh:
				timer.Reset(65 * time.Minute)
			case <-timer.C:
				panic("timeout for 65 minutes without message")
			}
		}
	}()

	for {
		select {
		case event := <-userCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				e := event.(*gomastodon.UpdateEvent)
				go HandleUpdate(e, con.ScopeTypeLocal)
			case *gomastodon.DeleteEvent:
				e := event.(*gomastodon.DeleteEvent)
				go HandleDelete(e, con.ScopeTypeLocal)
			case *gomastodon.NotificationEvent:
				e := event.(*gomastodon.NotificationEvent)
				go HandleNotification(e)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}

		case event := <-publicCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				e := event.(*gomastodon.UpdateEvent)
				ntfCh <- struct{}{}
				go HandleUpdate(e, con.ScopeTypePublic)
			case *gomastodon.DeleteEvent:
				e := event.(*gomastodon.DeleteEvent)
				go HandleDelete(e, con.ScopeTypePublic)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}
		}
	}
}
