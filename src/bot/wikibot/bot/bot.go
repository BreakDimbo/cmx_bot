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
	"regexp"
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
		toot := &gomastodon.Toot{Status: content}
		status, err := b.client.RawPost(toot)
		if err != nil {
			zlog.SLogger.Errorf("post to mastodon error: %v", err)
		}

		// record toot id
		wiki.Url = status.URL
		err = wiki.Store()
		if err != nil {
			zlog.SLogger.Errorf("update wiki url error: %v", err)
		}

	} else if strings.Contains(ntf.Status.Content, "?") { // query with ?xxx
		reg := regexp.MustCompile(`^?\S*`)
		kword := "#" + reg.FindString(filter(ntf.Status.Content))[1:]

		// query from elastic
		wiki := indexWiki{Word: kword}
		urls := wiki.QueryByWord()
		urlsStr := strings.Join(urls, "\n")
		tootContent := fmt.Sprintf("@%s\n%s", ntf.Account.Username, urlsStr)

		// reply to user
		toot := &gomastodon.Toot{
			Status:      tootContent,
			InReplyToID: ntf.Status.ID,
			Visibility:  "direct",
		}

		_, err := b.client.RawPost(toot)
		if err != nil {
			zlog.SLogger.Errorf("post to user %s error: %v", ntf.Account.Username, err)
		}
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
