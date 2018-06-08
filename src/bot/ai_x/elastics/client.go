package elastics

import (
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

var once sync.Once
var err error
var Client *elastic.Client

const URL = "http://47.93.43.59:9201"

func InitOnce() {
	once.Do(func() {
		cfg := &config.Config{URL: URL, Username: "elastic", Password: "break12345"}
		Client, err = elastic.NewClientFromConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}
		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := Client.Ping(URL).Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := Client.ElasticsearchVersion(URL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Elasticsearch version %s\n", esversion)

		// Use the IndexExists service to check if a specified index exists.
		exists, err := Client.IndexExists("status").Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			// create mapping
			createIndex, err := Client.CreateIndex("status").Body(con.StatusMapping).Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			if !createIndex.Acknowledged {
				log.Println("createIndex not acknowledged")
			}
		}
	})
}
