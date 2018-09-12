package main

import (
	"bot/config"
	"bot/log"
	"bot/theater/bot"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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
		id, name, line, err := parseText(content)
		if err != nil {
			log.SLogger.Errorf("parse text:[%s] error: %v", content, err)
			continue
		}

		actor, ok := actors[name]
		if !ok {
			log.SLogger.Errorf("not find actor by name: %s on line id: %s", name, id)
			continue
		}

		select {
		case actor.LineCh <- line:
		default:
			log.SLogger.Errorf("actor %s LineCh blocked with line id: %s", actor.Name, id)
		}

		time.Sleep(30 * time.Minute)
	}
}

/*
line example:

id/name/line
*/
func parseText(content string) (string, string, string, error) {
	s := strings.Split(content, "/")
	if len(s) < 3 {
		return "", "", "", fmt.Errorf("split content [%s] len less 3 error", content)
	}

	id, name, line := s[0], s[1], s[2]
	return id, name, line, nil
}
