package bot

import (
	"bot/client"
	"bot/config"
	gomastodon "bot/go-mastodon"
	"bot/intelbot/elastics"
	zlog "bot/log"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	NtfMention = "mention"
)

type WikiBot struct {
	client *client.Bot
}

func New(config config.MastodonClientInfo) (*WikiBot, error) {
	wikibot := new(WikiBot)
	client, err := client.New(&config)
	if err != nil {
		return nil, err
	}

	wikibot.client = client
	return wikibot, nil
}

func (b *WikiBot) Launch() {
	ctx, cancel := context.WithCancel(context.Background())
	userCh, err := b.client.WS.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer func() {
		cancel()
		close(userCh)
	}()

	for evt := range userCh {
		switch e := evt.(type) {
		case *gomastodon.NotificationEvent:
			b.handleNotification(e)
		case *gomastodon.DeleteEvent:
			b.handleDelete(e)
		default:
			// zlog.SLogger.Infof("receive other event: %s", e)
		}
	}
}

func (b *WikiBot) handleNotification(e *gomastodon.NotificationEvent) {
	ntf := e.Notification
	if ntf.Type != NtfMention {
		return
	}

	if strings.Contains(ntf.Status.Content, "#") { // add with #xxx
		kword, article := parseToot(ntf.Status.Content)
		fromUser := ntf.Account.Username
		addTime := time.Now().Format("Jan 2 15:04:05")
		content := fmt.Sprintf("@%s-%s-添加条目：%s\n解释：%s", fromUser, addTime, kword, article)
		// add to elastic
		wiki := indexWiki{
			ID:        ntf.ID,
			CreatedAt: time.Now(),
			Word:      kword,
			Content:   content,
		}
		err := wiki.Store()
		if err != nil {
			zlog.SLogger.Errorf("store wiki content to elastic error: %v", err)
		}

		// post to tl
		/*
			name-time-添加条目：XXX
			解释：xxxxxxx
		*/
		_, err = b.client.Post(content)
		if err != nil {
			zlog.SLogger.Errorf("post to mastodon error: %v", err)
		}

	} else if strings.Contains(ntf.Status.Content, "?") { // query with ?xxx
		// query from elastic

		// report to user
	}
}

func (b *WikiBot) handleDelete(e *gomastodon.DeleteEvent) {
	ctx := context.Background()
	_, err := elastics.Client.Delete().Index("wiki").Type("wiki").Id(e.ID).Do(ctx)
	if err != nil {
		return
	}
	zlog.SLogger.Infof("delete from es ok with id: %s", e.ID)
}
