package mastodon

import (
	con "bot/ai_x/const"
	"bot/ai_x/elastics"
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"
	"unicode"

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
	wfMap := calWordFrequency(totalToots)
	extractKeyWord(3, wfMap)
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

func calWordFrequency(totalToots map[string]*indexStatus) (wFreMap map[string]int) {
	x := gojieba.NewJieba()
	defer x.Free()
	use_hmm := true
	wFreMap = make(map[string]int)

	for _, s := range totalToots {
		words := x.Cut(s.Content, use_hmm)
		for _, w := range words {
			if len(w) <= con.SingleChineseByte {
				continue
			} else if stopwords[w] {
				continue
			} else {
				hasAlphabet := false
				for _, r := range w {
					if !unicode.Is(unicode.Scripts["Han"], r) {
						hasAlphabet = true
						break
					}
				}
				if hasAlphabet {
					continue
				}
			}
			wFreMap[w] += 1
		}
	}

	fmt.Printf("[DEBUG] calculate word frequency result: %s\n", wFreMap)
	return
}

func extractKeyWord(top int, wfMap map[string]int) (keywords []wordPair) {
	wfMapLen := len(wfMap)
	keywords = make([]wordPair, wfMapLen)
	for k, v := range wfMap {
		keywords = append(keywords, wordPair{key: k, value: v})
	}
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].value > keywords[j].value
	})
	if top > wfMapLen {
		keywords = keywords[:wfMapLen]
	} else {
		keywords = keywords[:top]
	}
	fmt.Printf("[DEBUG] keywords result: %s\n", keywords)
	return
}

func generateWordCloud() (medisId string) {
	// upload media
	return ""
}
