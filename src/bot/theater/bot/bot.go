package bot

import (
	"bot/bredis"
	"bot/client"
	"bot/config"
	"bot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	LoveYouKey     = "LoveKurisu"
	LoveYouTimeout = 6 * time.Hour
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
		Name:      name,
		LineCh:    make(chan string),
		BlockCh:   make(chan string),
		UnBlockCh: make(chan string),
		client:    c,
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
			a.handleNotification(n, actors)
		default:
			// log.SLogger.Infof("receive other event: %s", uq)
		}
	}
}

func (a *Actor) handleNotification(ntf *gomastodon.NotificationEvent, actors map[string]*Actor) {
	n := ntf.Notification
	if n.Status == nil {
		return
	}
	content := filter(n.Status.Content)
	log.SLogger.Infof("get notification: %s", content)

	switch a.Name {
	case cons.Okabe:
		if strings.Contains(content, "EL_PSY_CONGROO") {
			switch a.Name {
			case cons.Okabe:
				for _, actor := range actors {
					if actor.Name == cons.Okabe {
						continue
					}
					actor.BlockCh <- string(n.Account.ID)
					log.SLogger.Infof("start to block %s", n.Account.ID)
				}
			}
		} else if strings.Contains(content, "Love_You") {
			for _, actor := range actors {
				if actor.Name == cons.Okabe {
					continue
				}
				actor.UnBlockCh <- string(n.Account.ID)
			}
		}
	case cons.Kurisu:
		if isLoveYou(content) {
			// if the toot is for kurisu and on public then kurisu will reply he(she) on public line
			if n.Status.Visibility == "public" {

				key := fmt.Sprintf("%s:%s", LoveYouKey, n.Account.Username)
				// if loved already, toot hentai and return
				res, err := bredis.Client.Get(key).Result()
				if err == nil && res != "" {
					toot := fmt.Sprintf("@%s %s", n.Account.Username, "够了！变态！")
					_, err = a.client.Post(toot)
					if err != nil {
						log.SLogger.Errorf("kurisu reply to error %v", err)
					}
					return
				}

				err = bredis.Client.Set(key, n.Account.Username, LoveYouTimeout).Err()
				if err != nil {
					log.SLogger.Errorf("set key to redis error: %v", err)
				}
				reply := selectReply(cons.Kurisu)
				toot := fmt.Sprintf("@%s %s", n.Account.Username, reply)
				_, err = a.client.Post(toot)
				if err != nil {
					log.SLogger.Errorf("kurisu reply to error %v", err)
				}
			}
		}
	case cons.Itaru:
		if strings.Contains(content, "#菜谱") {
			food := strings.Trim(content, "@itaru")
			i := strings.Index(food, "#菜谱")
			food = food[i+7:]
			key := fmt.Sprintf("%s:%s", FoodKey, food)
			err := bredis.Client.Set(key, "true", 1024*24*time.Hour).Err()
			if err != nil {
				log.SLogger.Errorf("save %s to redis error: %v", key, err)
				return
			}

			toot := fmt.Sprintf("@%s %s", n.Account.Username, "乙！")
			script := fmt.Sprintf("诶嘿嘿，%s 怎么样？", food)
			iteraSlice = append(iteraSlice, script)
			_, err = a.client.Post(toot)
			if err != nil {
				log.SLogger.Errorf("kurisu reply to error %v", err)
			}
		} else if strings.Contains(content, "桶子") && n.Status.Visibility == "public" {
			reply := selectReply(cons.Itaru)
			toot := fmt.Sprintf("@%s %s", n.Account.Username, reply)
			_, err := a.client.Post(toot)
			if err != nil {
				log.SLogger.Errorf("kurisu reply to error %v", err)
			}
		}
	}
}
