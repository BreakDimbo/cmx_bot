package elastics

import (
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olivere/elastic"
)

var once sync.Once
var err error
var Client *elastic.Client

func InitOnce() {
	once.Do(func() {
		Client, err = elastic.NewClient()
		if err != nil {
			log.Fatal(err)
		}
		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := Client.Ping("http://127.0.0.1:9200").Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := Client.ElasticsearchVersion("http://127.0.0.1:9200")
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
