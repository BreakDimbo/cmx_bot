package main

import (
	"bot/theater/bot"
	"sync"
)

var wg sync.WaitGroup

func main() {
	actors := make(map[string]*bot.Actor)
	actorsName := []string{
		"okabe",
		"mayuri",
		"itaru",
		"kurisu",
		"moeka",
		"ruka",
		"nyannyan",
		"suzuha",
		"maho",
		"kagari",
		"yuki",
		"tennouji",
		"nae",
	}
	for _, name := range actorsName {
		actor := bot.New(name)
		actors[name] = actor
		wg.Add(1)
		go actor.Act(&wg)
	}

	wg.Add(1)
	go sendLine(actors)

	wg.Wait()
}
