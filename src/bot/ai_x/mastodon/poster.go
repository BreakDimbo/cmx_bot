package mastodon

import (
	"context"
	"log"

	mastodon "bot/go-mastodon"

	"github.com/robfig/cron"
)

func post(c *mastodon.Client, toot string) {
	_, err := c.PostStatus(context.Background(), &mastodon.Toot{
		Status: toot,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RunPoster(c *mastodon.Client) {
	crontab := cron.New()
	crontab.AddFunc("0 0 8 * * *", func() {
		status := DoAnalyze()
		post(c, status)
	})
}
