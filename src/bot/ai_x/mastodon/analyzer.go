package mastodon

import (
	con "bot/ai_x/const"
	"bot/ai_x/elastics"
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/olivere/elastic"
	"github.com/yanyiwu/gojieba"
)

/*
word frequency
daily toots number
daily active users
daily most positive user
weekly toots number trend
weekly active user trend
todo daily popular toot
*/

type wordPair struct {
	key   string
	value int
}

func DoAnalyzeDaily() string {
	now := time.Now().Add(-8 * time.Hour)
	sTime := now.Add(-24 * time.Hour)
	totalToots := fetchDataByTime(sTime, now)
	calWordFrequency(10, totalToots)
	return ""
}

func fetchDataByTime(startTime time.Time, endTime time.Time) (sResult map[string]*indexStatus) {

	stStr := startTime.Format(con.RFC3339local)
	edStr := endTime.Format(con.RFC3339local)
	query := elastic.NewRangeQuery("created_at").
		Gte(stStr).
		Lte(edStr)

	searchResult, err := elastics.Client.Search().
		Index("status").
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.Fatalf("search error: %s", err)
		return nil
	}

	var toot *indexStatus
	sResult = make(map[string]*indexStatus)
	for _, item := range searchResult.Each(reflect.TypeOf(toot)) {
		t := item.(*indexStatus)
		sResult[t.ID] = t
	}

	fmt.Printf("[DEBUG] fetch data by time result: %s\n", sResult)
	return
}

func calWordFrequency(limit int, totalToots map[string]*indexStatus) (wFreMap map[string]int) {
	x := gojieba.NewJieba()
	defer x.Free()
	use_hmm := true
	wFreMap = make(map[string]int)

	for _, s := range totalToots {
		words := x.Cut(s.Content, use_hmm)
		for _, w := range words {
			if len(w) <= con.SingleChineseByte {
				continue
			}
			wFreMap[w] += 1
		}
	}

	fmt.Printf("[DEBUG] calculate word frequency result: %s\n", wFreMap)
	return
}

func generateWordCloud() (medisId string) {
	// upload media
	return ""
}
