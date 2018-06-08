package mastodon

import (
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

var sResult map[string]*indexStatus

func DoAnalyzeDaily() string {
	now := time.Now().Add(-8 * time.Hour)
	sTime := now.Add(-24 * time.Hour)
	fetchDataByTime(sTime, now)
	calWordFrequency(10)
	return ""
}

func fetchDataByTime(startTime time.Time, endTime time.Time) {
	RFC3339local := "2006-01-02T15:04:05Z"
	stStr := startTime.Format(RFC3339local)
	edStr := endTime.Format(RFC3339local)
	query := elastic.NewRangeQuery("created_at").
		Gte(stStr).
		Lte(edStr)

	searchResult, err := elastics.Client.Search().
		Index("status").
		Query(query).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		log.Fatalf("search error: %s", err)
		return
	}

	var toot *indexStatus
	sResult = make(map[string]*indexStatus)
	for _, item := range searchResult.Each(reflect.TypeOf(toot)) {
		t := item.(*indexStatus)
		sResult[t.ID] = t
	}

	fmt.Printf("[DEBUG] fetch data by time result: %s\n", sResult)
}

func calWordFrequency(limit int) (wFreMap map[string]int) {
	x := gojieba.NewJieba()
	defer x.Free()

	s := "我是天才"
	words = x.Cut(s, use_hmm)
	fmt.Printf("analyze result: %s\n", words)
	return nil
}

func generateWordCloud() (medisId string) {
	// upload media
	return ""
}
