package bot

import (
	"bot/client"
	"bot/config"
	"bot/log"
	"sync"
)

type Actor struct {
	Name   string
	LineCh chan string
	client *client.Bot
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

	for line := range a.LineCh {
		_, err := a.client.PostSpoiler("来自草莓县石头门bot剧组", line)
		if err != nil {
			log.SLogger.Errorf("%s post line [%s] to mastodon error: %v", a.Name, line, err)
		}
	}
}
