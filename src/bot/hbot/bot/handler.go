package bot

import (
	"bot/bredis"
	gomastodon "bot/go-mastodon"
	con "bot/hbot/const"
	"bot/log"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func HandleNotification(e *gomastodon.NotificationEvent) {
	var tootToPost string
	notify := e.Notification

	switch notify.Type {
	case con.NotificationTypeMention:
		fromUser := notify.Account
		toot := notify.Status
		firstContent := filter(toot.Content)

		if fromUser.Username == "xbot" || fromUser.Username == "zbot" {
			return
		}

		tootToPost = firstContent

		// id, _ := botClient.Post(tootToPost)
		spoilerText := filter(toot.SpoilerText)
		if spoilerText == "" {
			spoilerText = "学习资料"
		}
		status, _ := botClient.PostSensetiveWithPic(spoilerText, tootToPost, toot.Sensitive, toot.MediaAttachments)

		err := bredis.Client.Set(string(toot.ID), string(status.ID), con.TootIDRedisTimeout).Err()
		if err != nil {
			log.SLogger.Errorf("set id: %s to redis error: %s", status.ID, err)
		}

		log.SLogger.Infof("get toots from user id: %s, first toot: %s, tootToPost: %s\n", fromUser.Username, firstContent, tootToPost)
	}
}

func HandleDelete(e *gomastodon.DeleteEvent) {
	id := e.ID
	fbotTootID, err := bredis.Client.Get(string(id)).Result()
	if err != nil {
		// log.SLogger.Errorf("get toot: %s from redis error: %s", id, err)
		return
	}
	log.SLogger.Infof("start to delete toot with id: %s", fbotTootID)
	botClient.DeleteToot(fbotTootID)
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	p.AllowElements("br")
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@firebot", "", -1)
	polished = strings.Replace(polished, "@fbot", "", -1)
	polished = strings.Replace(polished, "@hbot", "", -1)
	polished = strings.Replace(polished, "<br/>", "\n", -1)
	return
}
