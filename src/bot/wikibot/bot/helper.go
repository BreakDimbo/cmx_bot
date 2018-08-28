package bot

import (
	"bot/intelbot/elastics"
	zlog "bot/log"
	"context"
	"html"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/olivere/elastic"
)

func parseToot(toot string) (string, string) {
	toot = filter(toot)

	var kword, article string
	sharpreg := regexp.MustCompile(`^#\S*`)
	kword = sharpreg.FindString(toot)

	// replace all #xxx
	article = sharpreg.ReplaceAllString(toot, "")

	zlog.SLogger.Debugf("keyword: %s, article: %s", kword, article)

	return kword, article
}

type indexWiki struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Word      string    `json:"word"`
	Url       string    `json:"url"`
}

func (w *indexWiki) Store() error {
	ctx := context.Background()
	p, err := elastics.Client.Index().
		Index("wiki").
		Type("wiki").
		Id(w.ID).
		BodyJson(w).
		Do(ctx)
	if err != nil {
		return err
	}

	zlog.SLogger.Infof("indexed status %s to index %s, type %s", p.Id, p.Index, p.Type)
	return nil
}

func (w *indexWiki) QueryByWord() []string {
	urls := make([]string, 1)
	query := elastic.NewBoolQuery()
	termQuery := elastic.NewTermQuery("word", w.Word)
	query = query.Must(termQuery)

	result, err := elastics.Client.Search().
		Index("wiki").
		Query(query).
		Size(10000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		zlog.SLogger.Errorf("search from elastic error: %s", err)
		return nil
	}

	for _, item := range result.Each(reflect.TypeOf(w)) {
		i := item.(*indexWiki)
		urls = append(urls, i.Url)
	}

	zlog.SLogger.Debugf("query word: %s, result: %v", w.Word, result)

	return urls
}

func filter(raw string) (polished string) {
	p := bluemonday.StrictPolicy()
	p.AllowElements("br")
	polished = p.Sanitize(raw)
	polished = strings.Replace(polished, "@firebot", "", -1)
	polished = strings.Replace(polished, "@fbot", "", -1)
	polished = strings.Replace(polished, "@wbot", "", -1)
	polished = strings.Replace(polished, "<br/>", "\n", -1)
	polished = html.UnescapeString(polished)
	return
}
