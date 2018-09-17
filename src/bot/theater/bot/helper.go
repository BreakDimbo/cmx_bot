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
