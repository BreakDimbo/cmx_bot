package mastodon

import (
	"bot/config"
	mastodon "bot/go-mastodon"
	"context"
	"fmt"
)

func post(toot string) {
	pc := config.GetPostConfig()
	_, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: pc.Scope,
	})
	if err != nil {
		fmt.Printf("[ERROR]: post error: %s", err)
		return
	}
}
