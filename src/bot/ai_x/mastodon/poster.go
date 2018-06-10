package mastodon

import (
	"bot/ai_x/config"
	"context"
	"log"
	"sync"

	mastodon "bot/go-mastodon"

	"github.com/robfig/cron"
)

func post(c *mastodon.Client, toot string) {
	_, err := c.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: "private",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RunPoster(c *mastodon.Client) {
	var once sync.Once
	cronConfig := config.GetPostCron()
	once.Do(func() {
		crontab := cron.New()
		crontab.AddFunc(cronConfig.ConTime, func() {
			status := DailyAnalyze()
			post(c, status)
		})
		crontab.Start()
	})
}
