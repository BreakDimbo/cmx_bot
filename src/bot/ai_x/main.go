package main

import (
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"

	mastodon "github.com/mattn/go-mastodon"
)

func main() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       con.Server,
		ClientID:     con.ClientId,
		ClientSecret: con.ClientSecret,
	})

	wsClient := c.NewWSClient()
	ctx, cancel := context.WithCancel(context.Background())
	q, err := wsClient.StreamingWSPublic(ctx, true)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	err = c.Authenticate(context.Background(), con.ClientEmail, con.ClientPassword)
	if err != nil {
		log.Fatal(err)
	}
	for event := range q {
		switch event.(type) {
		case *mastodon.UpdateEvent:
			updateEvent := event.(*mastodon.UpdateEvent)
			fmt.Println(updateEvent.Status.Content)
		}
	}
}

// only used once
func registerApp() {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     con.Server,
		ClientName: "reader",
		Scopes:     "read write follow",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("client-id    : %s\n", app.ClientID)
	fmt.Printf("client-secret: %s\n", app.ClientSecret)
}
