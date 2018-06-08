package main

import (
	"bot/ai_x/elastics"
	"bot/ai_x/mastodon"
)

func main() {
	elastics.InitOnce()
	mastodon.InitOnce()
	mastodon.LoadStopWord()
	mastodon.Lauch()
}

/* only used once
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
*/
