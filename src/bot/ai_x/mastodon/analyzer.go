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

	gomastodon "bot/go-mastodon"

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

func DailyAnalyze() string {
	now := time.Now().Add(4 * time.Hour)
	sTime := now.Add(-20 * time.Hour)
	totalToots := fetchDataByTime(sTime, now)
	wfMap := calWordFrequency(totalToots)
	wpairs := extractKeyWord(20, wfMap)
	tootsCount := len(totalToots)
	tpMap := tootsByPerson(totalToots)
	activePersonNum := len(tpMap)
	id, num := mostActivePerson(tpMap)
	account, err := client.GetAccount(context.Background(), gomastodon.ID(id))
	if err != nil {
		log.Fatalf("get account with id: %s error", id, err)
	}

	tootToPost := fmt.Sprintf("1.昨日本县关键词前五名：%s | %s | %s | %s | %s\n 2.昨日本县嘟嘟数：%d\n 3.昨日本县冒泡人数：%d\n 4.昨日最活跃县民：%s, 共嘟嘟了%d条\n",
		wpairs[0].key, wpairs[1].key, wpairs[2].key, wpairs[3].key, wpairs[4].key, tootsCount,
		activePersonNum, account.Username, num)
	return tootToPost
}

func WeeklyAnalyze() string {
	now := time.Now().Add(4 * time.Hour)
	sTime := now.Add(-164 * time.Hour)
	totalToots := fetchDataByTime(sTime, now)
	wfMap := calWordFrequency(totalToots)
	wpairs := extractKeyWord(20, wfMap)
	tootsCount := len(totalToots)
	tpMap := tootsByPerson(totalToots)
	activePersonNum := len(tpMap)
	id, num := mostActivePerson(tpMap)
	account, err := client.GetAccount(context.Background(), gomastodon.ID(id))
	if err != nil {
		log.Fatalf("get account with id: %s error", id, err)
	}

	tootToPost := fmt.Sprintf("1.上周本县关键词前五名：%s | %s | %s | %s | %s\n 2.上周本县嘟嘟数：%d\n 3.上周本县冒泡人数：%d\n 4.上周最活跃县民：%s, 共嘟嘟了%d条\n",
		wpairs[0].key, wpairs[1].key, wpairs[2].key, wpairs[3].key, wpairs[4].key, tootsCount,
		activePersonNum, account.Username, num)
	return tootToPost
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
		Size(10000).
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
			if stopwords[w] {
				continue
			}

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

func tootsByPerson(totalToots map[string]*indexStatus) (tootsNumPersonMap map[string]int) {
	tootsNumPersonMap = make(map[string]int)
	for _, k := range totalToots {
		tootsNumPersonMap[k.AccountId] += 1
	}
	return
}

func mostActivePerson(tpMap map[string]int) (id string, tootNum int) {
	for k, v := range tpMap {
		if v >= tootNum {
			tootNum = v
			id = k
		}
	}
	return
}
