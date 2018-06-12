package mastodon

import (
	"bot/ai_x/config"
	"context"
	"log"

	mastodon "bot/go-mastodon"
)

func post(toot string) {
	pc := config.GetPostConfig()
	_, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func DailyPost() {
	status := DailyAnalyze()
	post(status)
}

func WeeklyPost() {
	status := WeeklyAnalyze()
	post(status)
}
