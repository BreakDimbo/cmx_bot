package mastodon

import (
	con "bot/ai_x/const"
	"bot/ai_x/elastics"
	"bot/config"
	"context"
	"fmt"
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

type kvPair struct {
	key   string
	value int
}

func DailyAnalyze() string {
	cf := config.GetMastodonClientInfo()
	now := time.Now().Add(cf.Timezone * time.Hour)
	sTime := now.Add((-24 + cf.Timezone) * time.Hour)
	totalToots := fetchDataByTime(sTime, now, con.ScopeTypePublic)
	localToots := fetchDataByTime(sTime, now, con.ScopeTypeLocal)
	wfMap := calWordFrequency(totalToots)
	wpairs := topN(20, wfMap)
	tootsCount := len(totalToots)
	tpMap := tootsByPerson(totalToots)
	ltpMap := tootsByPerson(localToots)
	activePersonNum := len(tpMap)
	tpSlice := topN(3, tpMap)
	var topAccounts []*kvPair
	for _, v := range tpSlice {
		account, err := client.GetAccount(context.Background(), gomastodon.ID(v.key))
		if err != nil {
			fmt.Printf("[ERROR] get account with id: %s error: %s\n", v.key, err)
			tpaccount := &kvPair{key: "无", value: v.value}
			topAccounts = append(topAccounts, tpaccount)
			continue
		}
		name := fmt.Sprintf("%s·%s", account.DisplayName, account.Username)
		topAccounts = append(topAccounts, &kvPair{key: name, value: v.value})
	}

	var hualao string
	lid, lnum := mostActivePerson(ltpMap)
	laccount, lerr := client.GetAccount(context.Background(), gomastodon.ID(lid))
	if lerr != nil {
		fmt.Printf("[ERROR] get account with id: %s error: %s\n", lid, lerr)
		hualao = "无"
	} else {
		hualao = fmt.Sprintf("%s·%s", laccount.DisplayName, laccount.Username)
	}

	tootToPost := fmt.Sprintf("1.昨日本县关键词前五名：%s(%d) | %s(%d) | %s(%d) | %s(%d) | %s(%d)\n2.昨日本县嘟嘟数：%d\n3.昨日本县冒泡人数：%d\n4.昨日最活跃县民榜：\n(^з^)-☆ %s,嘟嘟%d条\n(^з^)-☆ %s,嘟嘟%d条\n(^з^)-☆ %s,嘟嘟%d条\n5.昨日局长眼中话唠：\n(^з^)-☆ %s,嘟嘟%d条\n6.局长联动：本县入住传火局局长 @%s\n",
		wpairs[0].key, wpairs[0].value, wpairs[1].key, wpairs[1].value, wpairs[2].key, wpairs[2].value, wpairs[3].key, wpairs[3].value, wpairs[4].key, wpairs[4].value, tootsCount,
		activePersonNum, topAccounts[0].key, topAccounts[0].value, topAccounts[1].key, topAccounts[1].value, topAccounts[2].key, topAccounts[2].value, hualao, lnum, cf.Fbot)
	return tootToPost
}

func WeeklyAnalyze() string {
	now := time.Now().Add(4 * time.Hour)
	sTime := now.Add(-164 * time.Hour)
	totalToots := fetchDataByTime(sTime, now, con.ScopeTypePublic)
	localToots := fetchDataByTime(sTime, now, con.ScopeTypeLocal)
	wfMap := calWordFrequency(totalToots)
	wpairs := topN(20, wfMap)
	tootsCount := len(totalToots)
	tpMap := tootsByPerson(totalToots)
	ltpMap := tootsByPerson(localToots)
	activePersonNum := len(tpMap)
	id, num := mostActivePerson(tpMap)
	account, err := client.GetAccount(context.Background(), gomastodon.ID(id))
	if err != nil {
		fmt.Printf("[ERROR] get account with id: %s error: %s\n", id, err)
	}

	var hualao string
	lid, lnum := mostActivePerson(ltpMap)
	laccount, lerr := client.GetAccount(context.Background(), gomastodon.ID(lid))
	if lerr != nil {
		fmt.Printf("[ERROR] get account with id: %s error: %s\n", lid, lerr)
		hualao = "无"
	} else {
		hualao = fmt.Sprintf("%s@%s", laccount.DisplayName, laccount.Username)
	}

	tootToPost := fmt.Sprintf("1.上周本县关键词前五名：%s(%d) | %s(%d) | %s(%d) | %s(%d) | %s(%d)\n 2.上周本县嘟嘟数：%d\n 3.上周本县冒泡人数：%d\n 4.上周最活跃县民：%s@%s, 共嘟嘟了%d条\n 5.上周局长眼中话唠：%s, 共嘟嘟了%d条\n",
		wpairs[0].key, wpairs[0].value, wpairs[1].key, wpairs[1].value, wpairs[2].key, wpairs[2].value, wpairs[3].key, wpairs[3].value, wpairs[4].key, wpairs[4].value, tootsCount,
		activePersonNum, account.DisplayName, account.Username, num, hualao, lnum)
	return tootToPost
}

func fetchDataByTime(startTime time.Time, endTime time.Time, scope string) (sResult map[string]*indexStatus) {
	stStr := startTime.Format(con.RFC3339local)
	edStr := endTime.Format(con.RFC3339local)
	var index string
	switch scope {
	case con.ScopeTypeLocal:
		index = "local"
	case con.ScopeTypePublic:
		index = "status"
	}
	query := elastic.NewRangeQuery("created_at").
		Gte(stStr).
		Lte(edStr)

	searchResult, err := elastics.Client.Search().
		Index(index).
		Query(query).
		Size(10000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		fmt.Printf("[ERROR]:search error: %s\n", err)
		return nil
	}

	var toot *indexStatus
	sResult = make(map[string]*indexStatus)
	for _, item := range searchResult.Each(reflect.TypeOf(toot)) {
		t := item.(*indexStatus)
		sResult[t.ID] = t
	}
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
	return
}

func topN(top int, m map[string]int) (pair []kvPair) {
	mlen := len(m)
	pair = make([]kvPair, mlen)
	for k, v := range m {
		pair = append(pair, kvPair{key: k, value: v})
	}
	sort.Slice(pair, func(i, j int) bool {
		return pair[i].value > pair[j].value
	})
	if top > mlen {
		pair = pair[:mlen]
	} else {
		pair = pair[:top]
	}
	fmt.Printf("[DEBUG] top n results: %s\n", pair)
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
