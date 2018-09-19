package main

import (
	"bot/const"
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
		"tennouji",
		"nae",
		"nakabachi",
		// "maho",
		// "kagari",
		// "yuki",
	}
	for _, name := range actorsName {
		actor := bot.New(name)
		actors[name] = actor
		wg.Add(1)
		go actor.Act(&wg)
		if name == cons.Okabe || name == cons.Kurisu || name == cons.Itaru {
			go actor.ListenAudiences(actors)
		}
	}

	wg.Add(1)
	go sendLine(actors)

	wg.Wait()
}
