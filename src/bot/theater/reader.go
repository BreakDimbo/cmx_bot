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

	input := bufio.NewScanner(f)
	for input.Scan() {
		content := input.Text()
		ep, id, name, line, err := parseText(content)
		if err != nil {
			log.SLogger.Errorf("parse text:[%s] error: %v", content, err)
			continue
		}

		acted, err := checkActed(ep, id)
		if acted || err != nil {
			continue
		}

		actor, ok := actors[name]
		if !ok {
			log.SLogger.Errorf("not find actor by name: %s on line id: %s", name, id)
			continue
		}

		select {
		case actor.LineCh <- line:
			log.SLogger.Infof("acts ep %s id %s", ep, id)
		default:
			log.SLogger.Errorf("actor %s LineCh blocked with line id: %s", actor.Name, id)
		}

		for checkNight() {
			time.Sleep(5 * time.Minute)
		}

		time.Sleep(20 * time.Minute)
	}
}

/*
line example:

ep/id/name/line
*/
func parseText(content string) (string, string, string, string, error) {
	s := strings.Split(content, "/")
	if len(s) < 4 {
		return "", "", "", "", fmt.Errorf("split content [%s] len less 4 error", content)
	}

	ep, id, name, line := s[0], s[1], s[2], s[3]
	return ep, id, name, line, nil
}

func checkActed(ep string, id string) (bool, error) {
	key := fmt.Sprintf("%s:%s", cons.Stein, ep)
	value, err := bredis.Client.Get(key).Result()
	if err == nil {
		valueInt, _ := strconv.Atoi(value)
		idInt, _ := strconv.Atoi(id)

		if idInt <= valueInt {
			return true, nil
		}

		err := bredis.Client.Set(key, id, 7*24*time.Hour).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %s from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil

	} else if err == redis.Nil {
		err := bredis.Client.Set(key, id, 7*24*time.Hour).Err()
		if err != nil {
			log.SLogger.Errorf("set ep %s with id %s from redis error: %v", ep, id, err)
			return false, err
		}
		return false, nil
	}

	log.SLogger.Errorf("get ep %s with id %s from redis error: %v", ep, id, err)
	return false, err
}

func checkNight() bool {
	now := time.Now()
	start := 12
	end := 21
	if now.Hour() > start && now.Hour() < end {
		return true
	}
	return false
}
