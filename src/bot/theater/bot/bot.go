package bot

import (
	"bot/client"
	"bot/config"
	"bot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"context"
	"html"
	"strings"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

type Actor struct {
	Name      string
	LineCh    chan string
	BlockCh   chan string
	UnBlockCh chan string
	client    *client.Bot
}

func New(name string) *Actor {
	cfg, err := config.ActorBotClientInfo(name)
	if err != nil {
		panic(err)
	}
	c, err := client.New(&cfg)
	if err != nil {
		panic(err)
	}

	return &Actor{
		Name:   name,
		LineCh: make(chan string),
		client: c,
	}
}

func (a *Actor) Act(wg *sync.WaitGroup) {
	defer wg.Done()
	isContine := true

	for isContine {
		select {
		case line, ok := <-a.LineCh:
			if !ok {
				isContine = false
				break
			}
			_, err := a.client.PostSpoiler(line, "#来自草莓县石头门bot剧组")
			if err != nil {
				log.SLogger.Errorf("%s post line [%s] to mastodon error: %v", a.Name, line, err)
			}
		case accountID := <-a.BlockCh:
			a.client.BlockAccount(accountID)
			log.SLogger.Infof("block user %s ok", accountID)
		case accountID := <-a.UnBlockCh:
			a.client.UnBlockAccount(accountID)
		}
	}
}

func (a *Actor) ListenAudiences(actors map[string]*Actor) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userCh, err := a.client.WS.StreamingWSUser(ctx)
	if err != nil {
		log.SLogger.Errorf("new user ws connction error: %s", err)
		return
	}
	defer close(userCh)

	for ntf := range userCh {
		switch ntf.(type) {
		case *gomastodon.NotificationEvent:
			n := ntf.(*gomastodon.NotificationEvent)
			handleNotification(n, actors)
		default:
			// zlog.SLogger.Infof("receive other event: %s", uq)
		}
	}
}

func handleNotification(ntf *gomastodon.NotificationEvent, actors map[string]*Actor) {
	n := ntf.Notification
	content := filter(n.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	if strings.Contains(content, "EL_PSY_CONGROO") {
		for _, actor := range actors {
			if actor.Name == cons.Okabe {
				continue
			}
			actor.BlockCh <- string(n.Account.ID)
			log.SLogger.Infof("start to block %s", n.Account.ID)
		}
	} else if strings.Contains(content, "Love_You") {
		for _, actor := range actors {
			if actor.Name == cons.Okabe {
				continue
			}
			actor.UnBlockCh <- string(n.Account.ID)
		}
	}
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@rintarou", "", -1)
	polished = html.UnescapeString(polished)
	return
}
