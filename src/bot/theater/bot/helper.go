package bot

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

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
		strings.Contains(content, "吃了你") || strings.Contains(content, "好き") || strings.Contains(content, "吃掉你") ||
		strings.Contains(content, "梦到你") || strings.Contains(content, "吻你") || strings.Contains(content, "可爱")
}
