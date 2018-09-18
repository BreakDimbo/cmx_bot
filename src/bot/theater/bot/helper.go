package bot

import (
	"html"
	"math/rand"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func selectReply(name string) string {
	index := rand.Intn(len(replySlice))
	return replySlice[index]
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@rintarou", "", -1)
	polished = html.UnescapeString(polished)
	return
}

func isLoveYou(content string) bool {
	return strings.Contains(content, "Love_You") || strings.Contains(content, "love you") ||
		strings.Contains(content, "Love You") || strings.Contains(content, "爱你") || strings.Contains(content, "喜欢你") ||
		strings.Contains(content, "吃了你") || strings.Contains(content, "好き")
}
