package mastodon

import (
	gomastodon "bot/go-mastodon"
	"context"
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func HandleNotification(e *gomastodon.NotificationEvent) {
	switch e.Notification.Type {
	case "mention":
		n := e.Notification
		fromUser := n.Account
		toot := n.Status
		firstContent := filter(toot.Content)
		content := recurToot(n.Status.InReplyToID)
		tootToPost := fmt.Sprintf("@%s:%s// %s", fromUser.Username, firstContent, content)
		tootToPost = strings.TrimSuffix(tootToPost, "// ")
		post(tootToPost)

		fmt.Printf("[DEBUG] get toots fromuserid: %s, toot: %s, originToot: %s\n", fromUser.Username, filter(toot.Content), content)
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
		return fmt.Sprintf("@%s:%s// %s", originToot.Account.Username, polished, recurToot(originToot.InReplyToID))
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
