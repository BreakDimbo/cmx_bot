package main

import (
	"bot/ai_x/crontab"
	"bot/ai_x/mastodon"
)

func main() {
	mastodon.LoadStopWord()
	crontab.Start()
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
