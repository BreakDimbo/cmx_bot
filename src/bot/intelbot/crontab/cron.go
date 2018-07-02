package crontab

import (
	"bot/config"
	"bot/intelbot/bot"
	"sync"

	"github.com/robfig/cron"
)

func Start() {
	var once sync.Once
	cronConfig := config.GetPostConfig()
	once.Do(func() {
		crontab := cron.New()
		crontab.AddFunc(cronConfig.DailyTime, func() {
			bot.DailyPost()
		})
		crontab.AddFunc(cronConfig.WeeklyTime, func() {
			bot.WeeklyPost()
		})
		crontab.AddFunc(cronConfig.MonthlyTime, func() {
			bot.MonthlyPost()
		})
		crontab.AddFunc(cronConfig.CleanUnfollower, func() {
			bot.CleanUnfollower()
		})
		crontab.Start()
	})
}
