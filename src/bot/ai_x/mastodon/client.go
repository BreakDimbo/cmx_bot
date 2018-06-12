package mastodon

import (
	"bot/ai_x/config"
	"bot/ai_x/const"
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
	ctx, cancel := context.WithCancel(context.Background())
	q, err := wsClient.StreamingWSPublic(ctx, true)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

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
				e := uq.(*gomastodon.UpdateEvent)
				HandleUpdate(e, con.ScopeTypeLocal)
			case *gomastodon.DeleteEvent:
				e := uq.(*gomastodon.DeleteEvent)
				HandleDelete(e)
			case *gomastodon.NotificationEvent:
				e := uq.(*gomastodon.NotificationEvent)
				HandleNotification(e)
			default:
				fmt.Printf("other event: %s\n", uq)
			}

		case pq := <-q:
			switch pq.(type) {
			case *gomastodon.UpdateEvent:
				e := pq.(*gomastodon.UpdateEvent)
				HandleUpdate(e, con.ScopeTypePublic)
			case *gomastodon.DeleteEvent:
				e := pq.(*gomastodon.DeleteEvent)
				HandleDelete(e)
			default:
				fmt.Printf("other event: %s\n", pq)
			}
		}
	}
}
