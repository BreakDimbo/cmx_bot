package bot

import (
	"bot/bredis"
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
	if ntf.Account.Username == "xbot" {
		return
	}

	if strings.Contains(ntf.Status.Content, "#") { // add with #xxx
		b.addWiki(ntf)
	} else if strings.Contains(ntf.Status.Content, "?") { // query with ?xxx
		b.queryWiki(ntf)
	}
}

// the logic of this place is disgusting
func (b *WikiBot) handleDelete(e *gomastodon.DeleteEvent) {
	botTootID, err := bredis.Client.Get(string(e.ID)).Result()
	if err == nil {
		// the e.ID is a user toot ID
		// delete the relative botTootID
		err := b.client.DeleteToot(botTootID)
		if err != nil {
			zlog.SLogger.Errorf("delete toot from mastodon: %s error: %v", botTootID, err)
		}

		ctx := context.Background()
		_, err = elastics.Client.Delete().Index("wiki").Type("wiki").Id(e.ID).Do(ctx)
		if err != nil {
			zlog.SLogger.Errorf("delete from es error with id: %s, error: %v", e.ID, err)
		}
		zlog.SLogger.Debugf("delete from ok with id: %s", e.ID)
	}

}

func (b *WikiBot) addWiki(ntf *gomastodon.Notification) {
	kword, article := parseToot(ntf.Status.Content)
	fromUser := ntf.Account.Username
	addTime := time.Now().Format("Jan 2 15:04:05")
	content := fmt.Sprintf("@%s-%s-添加条目：%s\n解释：%s", fromUser, addTime, kword, article)

	// add to elastic
	wiki := indexWiki{
		ID:        string(ntf.Status.ID),
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

	status, err := b.client.PostSensetiveWithPic("", content, false, ntf.Status.MediaAttachments)
	if err != nil {
		zlog.SLogger.Errorf("post to mastodon error: %v", err)
	}

	// record toot id
	wiki.Url = status.URL
	err = wiki.Store()
	if err != nil {
		zlog.SLogger.Errorf("update wiki url error: %v", err)
	}

	// Store originID to newID map in redis
	ntfID := ntf.Status.ID
	statusID := status.ID
	err = bredis.Client.Set(string(ntfID), string(statusID), 30*7*24*time.Hour).Err()
	if err != nil {
		zlog.SLogger.Errorf("store origin id: %s with new status id: %s error: %v", ntfID, statusID, err)
	}
}

func (b *WikiBot) queryWiki(ntf *gomastodon.Notification) {
	reg := regexp.MustCompile(`\?\S*`)
	queryStr := reg.FindString(filter(ntf.Status.Content))
	kword := "#" + strings.Trim(queryStr, "?")

	zlog.SLogger.Debugf("query string: %s, keyword: %s", queryStr, kword)

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
