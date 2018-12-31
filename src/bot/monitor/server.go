package main

import (
	"bot/config"
	con "bot/intelbot/const"
	"bot/intelbot/elastics"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

func Listen() {
	http.HandleFunc("/tootCountHourly", handlerHourly)
	http.HandleFunc("/tootCountDaily", handlerDaily)
	log.Fatal(http.ListenAndServe("127.0.0.1:8085", nil))
}

type TootCountData struct {
	Count int64
	Time  string
}

func handlerHourly(w http.ResponseWriter, r *http.Request) {
	tcDataSet := make([]*TootCountData, 12)
	cfg := config.IntelBotClientInfo()
	endTime := time.Now().Add(13 * time.Hour)
	endTimeForSearch := endTime.Add(cfg.Timezone * time.Hour).Add(-13 * time.Hour)
	for i := 0; i < 12; i++ {
		endTime := endTime.Add(time.Duration(int64(-i)) * time.Hour)
		endTimeForSearch := endTimeForSearch.Add(time.Duration(int64(-i)) * time.Hour)
		startTimeForSearch := endTimeForSearch.Add(-1 * time.Hour)
		tootCount := countTootByTimeRange(startTimeForSearch, endTimeForSearch)
		hourStr := endTime.Hour()
		minuteStr := endTime.Minute()
		timeStr := fmt.Sprintf("%d:%d", hourStr, minuteStr)
		tcData := &TootCountData{tootCount, timeStr}
		tcDataSet = append(tcDataSet, tcData)
	}

	fmt.Printf("send tcDataSet: %v\n", tcDataSet)

	tcByte, err := json.Marshal(tcDataSet)
	if err != nil {
		fmt.Printf("json error: %s\n", err)
	}

	enableCors(&w)

	_, err = w.Write(tcByte)
	if err != nil {
		fmt.Printf("send response error: %s\n", err)
	}
}

func handlerDaily(w http.ResponseWriter, r *http.Request) {
	tcDataSet := make([]*TootCountData, 7)
	cfg := config.IntelBotClientInfo()
	endTime := time.Now().Add(13 * time.Hour)
	endTimeForSearch := endTime.Add(cfg.Timezone * time.Hour).Add(-13 * time.Hour)
	for i := 0; i < 8; i++ {
		endTime := endTime.Add(time.Duration(int64(-i*24)) * time.Hour)
		endTimeForSearch := endTimeForSearch.Add(time.Duration(int64(-i*24)) * time.Hour)
		startTimeForSearch := endTimeForSearch.Add(-24 * time.Hour)
		tootCount := countTootByTimeRange(startTimeForSearch, endTimeForSearch)
		timeStr := endTime.Weekday().String()
		tcData := &TootCountData{tootCount, timeStr}
		tcDataSet = append(tcDataSet, tcData)
	}

	fmt.Printf("send tcDataSet: %v\n", tcDataSet)

	tcByte, err := json.Marshal(tcDataSet)
	if err != nil {
		fmt.Printf("json error: %s\n", err)
	}

	enableCors(&w)

	_, err = w.Write(tcByte)
	if err != nil {
		fmt.Printf("send response error: %s\n", err)
	}
}

func countTootByTimeRange(startTime, endTime time.Time) int64 {
	stStr := startTime.Format(con.RFC3339local)
	edStr := endTime.Format(con.RFC3339local)

	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewRangeQuery("created_at").Gte(stStr).Lte(edStr))
	cfg := config.IntelBotClientInfo()
	termQuery := elastic.NewTermQuery("server", cfg.Sever)
	query = query.Must(termQuery)

	searchResult, err := elastics.Client.Search().
		Index("status").
		Query(query).
		Size(10000000).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search from elastic error: %s\n", err)
	}
	return searchResult.TotalHits()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
