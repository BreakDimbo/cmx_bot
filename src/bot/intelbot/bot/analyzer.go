package bot

import (
	"bot/config"
	con "bot/intelbot/const"
	"bot/intelbot/elastics"
	"bot/log"
	"context"
	"fmt"
	"math/rand"
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
	count int
}

func DailyAnalyze() (string, string) {
	return analyze(con.AnalyzeIntervalDaily)
}

func WeeklyAnalyze() string {
	t, _ := analyze(con.AnalyzeIntervalWeekly)
	return t
}

func MonthlyAnalyze() string {
	t, _ := analyze(con.AnalyzeIntervalMonthly)
	return t
}

func analyze(interval string) (toot string, hideToot string) {
	// TODO: refactor
	var startTime time.Time
	var localHuaLao string
	var accNameTootsCounts []kvPair

	config := config.IntelBotClientInfo()
	toTime := time.Now().Add(config.Timezone * time.Hour)
	switch interval {
	case con.AnalyzeIntervalDaily:
		startTime = toTime.Add((-24 + config.Timezone) * time.Hour)
	case con.AnalyzeIntervalWeekly:
		startTime = toTime.Add((7*-24 + config.Timezone) * time.Hour)
	case con.AnalyzeIntervalMonthly:
		startTime = toTime.Add((30*-24 + config.Timezone) * time.Hour)
	}

	publicToots := fetchDataByTime(startTime, toTime, con.ScopeTypePublic)
	localToots := fetchDataByTime(startTime, toTime, con.ScopeTypeLocal)

	wfMap := calWordFrequency(publicToots)
	wordcounts := topN(20, wfMap)
	publicTootCount := len(publicToots)
	publicTootsCbyP := tootsCountByPerson(publicToots)
	localTootsCbyP := tootsCountByPerson(localToots)
	activePersonCount := len(publicTootsCbyP)

	accIDTootsCounts := topN(3, publicTootsCbyP)
	for _, v := range accIDTootsCounts {
		account, err := botClient.Normal.GetAccount(context.Background(), gomastodon.ID(v.key))
		if err != nil {
			log.SLogger.Errorf("get account with id: %s error: %s", v.key, err)
			tpaccount := kvPair{key: "无", count: v.count}
			accNameTootsCounts = append(accNameTootsCounts, tpaccount)
			continue
		}
		name := fmt.Sprintf("%s·%s", account.DisplayName, account.Username)
		accNameTootsCounts = append(accNameTootsCounts, kvPair{key: name, count: v.count})
	}

	id, count := mostActivePerson(localTootsCbyP)
	laccount, lerr := botClient.Normal.GetAccount(context.Background(), gomastodon.ID(id))
	if lerr != nil {
		log.SLogger.Errorf("get account with id: %s error: %s", id, lerr)
		localHuaLao = "无"
	} else {
		localHuaLao = fmt.Sprintf("%s·%s", laccount.DisplayName, laccount.Username)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	emoji := con.Emoji[r1.Intn(len(con.Emoji))]

	var shiningToot *gomastodon.Status
	// if interval == con.AnalyzeIntervalDaily {
	// 	shiningToot = findMostShiningToot(publicToots)
	// }

	toot = parseToToot(interval, wordcounts, publicTootCount,
		activePersonCount, accNameTootsCounts, emoji, localHuaLao,
		count, config.FbotName, shiningToot)

	return toot, hideToot
}

func parseToToot(interval string, wordcounts []kvPair, publicTootCount int,
	activePersonCount int, accNameTootsCounts []kvPair, emoji string, localHuaLao string,
	huaLaoCount int, firebot string, shiningToot *gomastodon.Status) (toot string) {

	var intervalStr string
	switch interval {
	case con.AnalyzeIntervalDaily:
		intervalStr = "昨日"
	case con.AnalyzeIntervalWeekly:
		intervalStr = "上周"
	case con.AnalyzeIntervalMonthly:
		toot = fmt.Sprintf("%d年%d月本县最强话唠是：%s,共嘟嘟%d条", time.Now().Year(), int(time.Now().Add(-30*24*time.Hour).Month()),
			accNameTootsCounts[0].key, accNameTootsCounts[0].count)
		return
	}
	//TODO: use loop
	keyWordsStr := fmt.Sprintf("1.%s本县关键词前五名：%s(%d) | %s(%d) | %s(%d) | %s(%d) | %s(%d)\n",
		intervalStr,
		wordcounts[0].key, wordcounts[0].count,
		wordcounts[1].key, wordcounts[1].count,
		wordcounts[2].key, wordcounts[2].count,
		wordcounts[3].key, wordcounts[3].count,
		wordcounts[4].key, wordcounts[4].count)
	tootCountStr := fmt.Sprintf("2.%s本县嘟嘟数：%d\n", intervalStr, publicTootCount)
	activePersonCountStr := fmt.Sprintf("3.%s本县冒泡人数：%d\n", intervalStr, activePersonCount)
	mostActiveRankStr := fmt.Sprintf("4.%s最活跃县民榜：\n%s %s,嘟嘟%d条\n%s %s,嘟嘟%d条\n%s %s,嘟嘟%d条\n",
		intervalStr,
		emoji, accNameTootsCounts[0].key, accNameTootsCounts[0].count,
		emoji, accNameTootsCounts[1].key, accNameTootsCounts[1].count,
		emoji, accNameTootsCounts[2].key, accNameTootsCounts[2].count)
	secretaryHuaLaoStr := fmt.Sprintf("5.%s局长眼中话唠：\n%s %s,嘟嘟%d条\n",
		intervalStr, emoji, localHuaLao, huaLaoCount)
	secretaryCooperateStr := fmt.Sprintf("6.局长联动：本县入住传火局局长 @%s，扫黄局局长 @hbot，草莓百科 @wbot \n", firebot)
	toot = keyWordsStr + tootCountStr + activePersonCountStr + mostActiveRankStr + secretaryHuaLaoStr + secretaryCooperateStr
	// if interval == con.AnalyzeIntervalDaily {
	// 	toot = toot + fmt.Sprintf("8.昨日最✨嘟嘟来自：%s·%s，转嘟%d次，收藏%d次, 🔗:%s \n", shiningToot.Account.DisplayName,
	// 		shiningToot.Account.Username, shiningToot.ReblogsCount, shiningToot.FavouritesCount, shiningToot.URL)
	// }
	return
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

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewRangeQuery("created_at").Gte(stStr).Lte(edStr))
	cfg := config.IntelBotClientInfo()
	termQuery := elastic.NewTermQuery("server", cfg.Sever)
	query = query.Must(termQuery)

	searchResult, err := elastics.Client.Search().
		Index(index).
		Query(query).
		Size(10000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		log.SLogger.Errorf("search from elastic error: %s", err)
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

func findMostShiningToot(toots map[string]*indexStatus) (stoot *gomastodon.Status) {
	ctx := context.Background()
	shingNum := int64(0)

	log.SLogger.Infof("totoal toots num to cal shinging: %d", len(toots))
	for id, v := range toots {
		toot, err := botClient.Normal.GetStatus(ctx, id)
		if err != nil {
			log.SLogger.Errorf("get toot status error: %s", err)
			continue
		}
		(*v).ReblogsCount = toot.ReblogsCount
		(*v).FavouritesCount = toot.FavouritesCount

		elastics.Client.Update().Index(con.ScopeTypePublic).Type("status").Id(id).Doc(*v).Do(ctx)
		if err != nil {
			log.SLogger.Errorf("update favourite count to es error: %s", err)
		}

		n := toot.FavouritesCount + toot.ReblogsCount
		if n > shingNum {
			shingNum = n
			stoot = toot
		}
		time.Sleep(1 * time.Second)
		log.SLogger.Infof("over toot: %d", id)
	}
	return
}

func calWordFrequency(totalToots map[string]*indexStatus) (wFreMap map[string]int) {
	x := gojieba.NewJieba()
	addWord(x)
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

func addWord(x *gojieba.Jieba) {
	x.AddWord("夜光内裤")
	x.AddWord("炼金术士")
	x.AddWord("鲁便器")
	x.AddWord("邦站")
	x.AddWord("草莓县")
	x.AddWord("炭烧鸡")
	x.AddWord("吃烧鸡")
}

func topN(top int, m map[string]int) (pair []kvPair) {
	mlen := len(m)
	pair = make([]kvPair, mlen)
	for k, v := range m {
		pair = append(pair, kvPair{key: k, count: v})
	}
	sort.Slice(pair, func(i, j int) bool {
		return pair[i].count > pair[j].count
	})
	if top > mlen {
		pair = pair[:mlen]
	} else {
		pair = pair[:top]
	}
	log.SLogger.Infof("top %d results: %s", top, pair)
	return
}

func tootsCountByPerson(totalToots map[string]*indexStatus) (tootsNumPersonMap map[string]int) {
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
