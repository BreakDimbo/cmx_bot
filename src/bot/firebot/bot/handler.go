package bot

import (
	"bot/bredis"
	con "bot/firebot/const"
	gomastodon "bot/go-mastodon"
	"bot/log"
	"context"
	"fmt"
	"regexp"
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
		replyToID := notify.Status.InReplyToID
		firstContent := filter(toot.Content)

		if fromUser.Username == "xbot" || fromUser.Username == "zbot" {
			index := strings.Index(firstContent, "县民榜：\n")
			tootToPost = firstContent[index+12:]
			lastIndex := strings.Index(tootToPost, "条\n")
			tootToPost = fmt.Sprintf("%s:%s。", "咳咳...注意！昨天最活跃（话唠）县民是", tootToPost[:lastIndex+3])
		} else if strings.Contains(firstContent, "#树洞") || strings.Contains(filter(toot.SpoilerText), "#树洞") {
			tootToPost = firstContent
		} else {
			reg := regexp.MustCompile("^@[^ ]*")
			firstContent = reg.ReplaceAllString(firstContent, "")
			content := recurToot(replyToID)
			tootToPost = fmt.Sprintf("@%s:%s// %s", fromUser.Acct, firstContent, content)
			tootToPost = strings.TrimSuffix(tootToPost, "// ")
		}

		// id, _ := botClient.Post(tootToPost)
		id, _ := botClient.PostSensetiveWithPic(filter(toot.SpoilerText), tootToPost, toot.Sensitive, toot.MediaAttachments)

		err := bredis.Client.Set(string(toot.ID), string(id), con.TootIDRedisTimeout).Err()
		if err != nil {
			log.SLogger.Errorf("set id: %s to redis error: %s", id, err)
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

func recurToot(tootId interface{}) string {
	if tootId != nil {
		ctx := context.Background()
		originToot, err := botClient.Normal.GetStatus(ctx, tootId.(string))
		if err != nil {
			log.SLogger.Errorf("get status error: %s", err)
			return ""
		}
		polished := filter(originToot.Content)
		return fmt.Sprintf("@%s:%s// %s", originToot.Account.Acct, polished, recurToot(originToot.InReplyToID))
	}
	return ""
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	p.AllowElements("br")
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@firebot", "", -1)
	polished = strings.Replace(polished, "@fbot", "", -1)
	polished = strings.Replace(polished, "<br/>", "\n", -1)
	return
}
