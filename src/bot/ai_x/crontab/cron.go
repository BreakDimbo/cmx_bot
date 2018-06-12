package crontab

import (
	"bot/ai_x/config"
	"bot/ai_x/mastodon"
	"sync"

	"github.com/robfig/cron"
)

func Start() {
	var once sync.Once
	cronConfig := config.GetPostConfig()
	once.Do(func() {
		crontab := cron.New()
		crontab.AddFunc(cronConfig.DailyTime, func() {
			mastodon.DailyPost()
		})
		crontab.AddFunc(cronConfig.WeeklyTime, func() {
			mastodon.WeeklyPost()
		})
		crontab.AddFunc(cronConfig.CleanUnfollower, func() {
			mastodon.CleanUnfollower()
		})
		crontab.Start()
	})
}
