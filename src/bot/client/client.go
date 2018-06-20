package client

import (
	"bot/config"
	gomastodon "bot/go-mastodon"
	zlog "bot/log"
	"context"
	"log"
)

type BotClient struct {
	Normal *gomastodon.Client
	WS     *gomastodon.WSClient
}

// New new bot client which should be called only when init
func New(config *config.MastodonClientInfo) (*BotClient, error) {
	c := gomastodon.NewClient(&gomastodon.Config{
		Server:       config.Sever,
		ClientID:     config.ID,
		ClientSecret: config.Secret,
	})
	err := c.Authenticate(context.Background(), config.Email, config.Password)
	if err != nil {
		log.Fatalf("[Fatal]: authenticate error of mastodon client: %s\n", err)
		return nil, err
	}
	bc := &BotClient{Normal: c, WS: c.NewWSClient()}
	return bc, nil
}

func (bc *BotClient) Post(toot string) error {
	pc := config.GetPostConfig()
	_, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return err
	}
	return nil
}
