package main

import (
	"bot/intelbot/bot"
	"bot/intelbot/crontab"
	"bot/intelbot/monitor"
	log "bot/log"
)

func main() {
	defer log.Logger.Sync()
	bot.LoadStopWord()
	crontab.Start()
	go monitor.Launch()
	bot.Launch()
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
