package elastics

import (
	"bot/ai_x/config"
	con "bot/ai_x/const"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olivere/elastic"
)

var Client *elastic.Client

func init() {
	var once sync.Once
	var err error

	once.Do(func() {
		e := config.GetElastic()
		Client, err = elastic.NewClient(elastic.SetBasicAuth(e.Username, e.Password),
			elastic.SetURL(e.Url),
			elastic.SetSniff(false),
			elastic.SetRetrier(NewMyRetrier()),
		)
		if err != nil {
			fmt.Printf("[ERROR] new es client error: %s", err)
			return
		}
		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := Client.Ping(e.Url).Do(context.Background())
		if err != nil {
			fmt.Printf("[ERROR] ping es error: %s", err)
			return
		}
		fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := Client.ElasticsearchVersion(e.Url)
		if err != nil {
			fmt.Printf("[ERROR] get es version error: %s", err)
			return
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

		// Use the IndexExists service to check if a specified index exists.
		lexists, lerr := Client.IndexExists("local").Do(context.Background())
		if lerr != nil {
			log.Fatal(lerr)
		}
		if !lexists {
			// create mapping
			createIndex, err := Client.CreateIndex("local").Body(con.StatusMapping).Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			if !createIndex.Acknowledged {
				log.Println("createIndex not acknowledged")
			}
		}
	})
}
