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
	pc := config.GetPostConfig()
	_, err := c.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RunPoster(c *mastodon.Client) {
	var once sync.Once
	cronConfig := config.GetPostConfig()
	once.Do(func() {
		crontab := cron.New()
		crontab.AddFunc(cronConfig.DailyTime, func() {
			status := DailyAnalyze()
			post(c, status)
		})
		crontab.AddFunc(cronConfig.WeeklyTime, func() {
			status := WeeklyAnalyze()
			post(c, status)
		})
		crontab.Start()
	})
}
