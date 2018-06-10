package elastics

import (
	"bot/ai_x/config"
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olivere/elastic"
	econfig "github.com/olivere/elastic/config"
)

var Client *elastic.Client

func init() {
	var once sync.Once
	var err error

	once.Do(func() {
		t := false
		e := config.GetElastic()
		cfg := &econfig.Config{URL: e.Url, Username: e.Username, Password: e.Password, Sniff: &t}
		Client, err = elastic.NewClientFromConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}
		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := Client.Ping(e.Url).Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := Client.ElasticsearchVersion(e.Url)
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
