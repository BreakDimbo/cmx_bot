package client

import (
	"bot/config"
	gomastodon "bot/go-mastodon"
	zlog "bot/log"
	"context"
	"log"
	"strconv"
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

func (bc *BotClient) Post(toot string) (gomastodon.ID, error) {
	pc := config.GetPostConfig()
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *BotClient) PostSpoiler(spolier string, toot string) (gomastodon.ID, error) {
	pc := config.GetPostConfig()
	status, err := bc.Normal.PostStatus(context.Background(), &gomastodon.Toot{
		Status:      toot,
		Visibility:  pc.Scope,
		SpoilerText: spolier,
	})
	if err != nil {
		zlog.SLogger.Errorf("post toot: %s error: %s", toot, err)
		return "", err
	}
	return status.ID, nil
}

func (bc *BotClient) DeleteToot(id string) error {
	ctx := context.Background()
	fbotTootID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zlog.SLogger.Errorf("parse id: %s error: %s", id, err)
		return err
	}
	return bc.Normal.DeleteStatus(ctx, int64(fbotTootID))
}
