package main

import (
	"bot/bredis"
	"bot/config"
	"bot/const"
	"bot/log"
	"bot/theater/bot"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	ActInterval = 15 * time.Minute
	NightStart  = 12
	NightEnd    = 21
	Timeout     = 365 * 24 * 2 * time.Hour
)

func sendLine(actors map[string]*bot.Actor) {
	filename := config.ScriptFilePath()
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		for _, actor := range actors {
			close(actor.LineCh)
		}
	}()

	defer wg.Done()

	var id int
	var prevEp string

	input := bufio.NewScanner(f)
	for input.Scan() {
		id++
		content := input.Text()
		ep, name, line, err := parseText(content)
		if err != nil {
			log.SLogger.Errorf("parse text:[%s] error: %v", content, err)
			continue
		}

		if ep != prevEp {
			id = 1
			prevEp = ep
		}

		acted, err := checkActed(ep, id)
		if acted || err != nil {
			continue
		}

		for checkNight() {
			time.Sleep(5 * time.Minute)
		}

		actor, ok := actors[name]
		if !ok {
			log.SLogger.Errorf("not find actor by name: %s on line id: %d", name, id)
			continue
		}

		select {
		case actor.LineCh <- line:
			log.SLogger.Infof("acts ep %s id %d", ep, id)
		default:
			log.SLogger.Errorf("actor %s LineCh blocked with line id: %d", actor.Name, id)
		}

		time.Sleep(ActInterval)
	}
}

/*
line example:

ep/id/name/line
*/
func parseText(content string) (string, string, string, error) {
	s := strings.Split(content, "/")
	if len(s) < 3 {
		return "", "", "", fmt.Errorf("split content [%s] len less 3 error", content)
	}

	ep, name, line := s[0], s[1], s[2]
	return ep, name, line, nil
}

func checkActed(ep string, id int) (bool, error) {
	key := fmt.Sprintf("%s:%s", cons.Stein, ep)
	value, err := bredis.Client.Get(key).Result()
	if err == nil {
		valueInt, _ := strconv.Atoi(value)

		if id <= valueInt {
			return true, nil
		}

		err := bredis.Client.Set(key, id, Timeout).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %d from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil

	} else if err == redis.Nil {
		err := bredis.Client.Set(key, id, Timeout).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %d from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil
	}

	log.SLogger.Errorf("get ep %s with id %d from redis error: %v", ep, id, err)
	return false, err
}

func checkNight() bool {
	now := time.Now()

	if now.Hour() >= NightStart && now.Hour() < NightEnd {
		return true
	}
	return false
}
