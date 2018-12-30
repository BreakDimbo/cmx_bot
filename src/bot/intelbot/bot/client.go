package bot

import (
	"bot/client"
	"bot/config"
	con "bot/intelbot/const"
	"bot/intelbot/monitor"
	zlog "bot/log"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	gomastodon "bot/go-mastodon"
)

var (
	countHourly, countDaily int
	botClient               *client.Bot
	err                     error
)

var mu1 = &sync.Mutex{}
var mu2 = &sync.Mutex{}

func init() {
	var once sync.Once
	once.Do(func() {
		config := config.IntelBotClientInfo()
		botClient, err = client.New(&config)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func Launch() {
	ctx, cancel := context.WithCancel(context.Background())
	publicCh, err := botClient.WS.StreamingWSPublic(ctx, true)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	userCh, err := botClient.WS.StreamingWSUser(ctx)
	if err != nil {
		log.Fatal(err)
		cancel()
	}

	defer cancel()

	ntfCh := make(chan struct{})

	// if no message from server, restart the bot
	go func() {
		timer := time.NewTimer(65 * time.Minute)
		for {
			select {
			case <-ntfCh:
				timer.Reset(65 * time.Minute)
			case <-timer.C:
				panic("timeout for 65 minutes without message")
			}
		}
	}()

	// pusher toot data to monitor
	go func() {
		intervalHour := 60 * time.Minute
		intervalDaily := 24 * time.Hour

		hourTicker := time.NewTicker(intervalHour)
		dailyTicker := time.NewTicker(intervalDaily)

		for {
			select {
			case <-hourTicker.C:
				time := fmt.Sprintf("%d:%d", time.Now().Hour(), time.Now().Minute())
				newVisitsData := monitor.VisitsData{
					Count: countHourly,
					Time:  time,
				}
				monitor.Client.Trigger("tootCountHourly", "addNumber", newVisitsData)
				zlog.SLogger.Debugf("Trigger count %d", countHourly)
				SetCountHourly(0)
			case <-dailyTicker.C:
				time := fmt.Sprintf("%s", time.Now().Weekday().String())
				newVisitsData := monitor.VisitsData{
					Count: countDaily,
					Time:  time,
				}
				monitor.Client.Trigger("tootCountDaily", "addNumber", newVisitsData)
				SetCountDaily(0)
			}
		}
	}()

	for {
		select {
		case event := <-userCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				e := event.(*gomastodon.UpdateEvent)
				go HandleUpdate(e, con.ScopeTypeLocal)
			case *gomastodon.DeleteEvent:
				e := event.(*gomastodon.DeleteEvent)
				go HandleDelete(e, con.ScopeTypeLocal)
			case *gomastodon.NotificationEvent:
				e := event.(*gomastodon.NotificationEvent)
				go HandleNotification(e)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}

		case event := <-publicCh:
			switch event.(type) {
			case *gomastodon.UpdateEvent:
				AddCountDaily()
				AddCountHourly()
				e := event.(*gomastodon.UpdateEvent)
				ntfCh <- struct{}{}
				go HandleUpdate(e, con.ScopeTypePublic)
			case *gomastodon.DeleteEvent:
				RemoveCountDaily()
				e := event.(*gomastodon.DeleteEvent)
				go HandleDelete(e, con.ScopeTypePublic)
			default:
				zlog.SLogger.Infof("receive other event: %s", event)
				os.Exit(0)
			}
		}
	}
}

func AddCountHourly() {
	mu1.Lock()
	defer mu1.Unlock()
	countHourly++
}

func RemoveCountHourly() {
	mu1.Lock()
	defer mu1.Unlock()
	countHourly--
}

func SetCountHourly(i int) {
	mu1.Lock()
	defer mu1.Unlock()
	countHourly = i
}

func AddCountDaily() {
	mu2.Lock()
	defer mu2.Unlock()
	countDaily++
}

func RemoveCountDaily() {
	mu2.Lock()
	defer mu2.Unlock()
	countDaily--
}

func SetCountDaily(i int) {
	mu2.Lock()
	defer mu2.Unlock()
	countDaily = i
}
