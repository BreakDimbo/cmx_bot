package main

import (
	"bot/config"
	_ "bot/intelbot/elastics"
	"bot/wikibot/bot"
)

func main() {
	wikibot, err := bot.New(config.WikiBotClientInfo())
	if err != nil {
		panic(err)
	}
	wikibot.Launch()
}
