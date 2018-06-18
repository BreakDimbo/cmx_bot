package bot

import (
	con "bot/firebot/const"
	gomastodon "bot/go-mastodon"
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
		reg := regexp.MustCompile("^@(.*)[[:space:]]")
		firstContent = reg.ReplaceAllString(firstContent, "")
		if fromUser.Username == "xbot" || fromUser.Username == "zbot" {
			index := strings.Index(firstContent, "县民榜：\n")
			tootToPost = firstContent[index+12:]
			lastIndex := strings.Index(tootToPost, "条\n")
			tootToPost = fmt.Sprintf("%s:%s。", "咳咳...注意！昨天最活跃（话唠）县民是", tootToPost[:lastIndex+3])
		} else {
			content := recurToot(replyToID)
			tootToPost = fmt.Sprintf("@%s:%s// %s", fromUser.Acct, firstContent, content)
			tootToPost = strings.TrimSuffix(tootToPost, "// ")
		}

		botClient.Post(tootToPost)

		fmt.Printf("[DEBUG] get toots from user id: %s, firsttoot: %s, tootToPost: %s\n", fromUser.Username, firstContent, tootToPost)
	}
}

func recurToot(tootId interface{}) string {
	if tootId != nil {
		ctx := context.Background()
		originToot, err := botClient.Normal.GetStatus(ctx, tootId.(string))
		if err != nil {
			fmt.Printf("[Error] get status error: %s\n", err)
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
