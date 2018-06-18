package mastodon

import (
	gomastodon "bot/go-mastodon"
	"context"
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func HandleNotification(e *gomastodon.NotificationEvent) {
	var tootToPost string
	switch e.Notification.Type {
	case "mention":
		n := e.Notification
		fromUser := n.Account
		toot := n.Status
		firstContent := filter(toot.Content)
		if fromUser.Username == "xbot" || fromUser.Username == "zbot" {
			index := strings.Index(firstContent, "县民榜：\n")
			tootToPost = firstContent[index+12:]
			lastIndex := strings.Index(tootToPost, "条\n")
			tootToPost = fmt.Sprintf("%s:%s。", "咳咳...注意！昨天最活跃（话唠）县民是", tootToPost[:lastIndex+3])
		} else {
			content := recurToot(n.Status.InReplyToID)
			tootToPost = fmt.Sprintf("@%s:%s// %s", fromUser.Acct, firstContent, content)
			tootToPost = strings.TrimSuffix(tootToPost, "// ")
		}

		post(tootToPost)

		fmt.Printf("[DEBUG] get toots fromuserid: %s, firsttoot: %s, tootToPost: %s\n", fromUser.Username, firstContent, tootToPost)
	}
}

func recurToot(tootId interface{}) string {
	if tootId != nil {
		ctx := context.Background()
		originToot, err := client.GetStatus(ctx, tootId.(string))
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
