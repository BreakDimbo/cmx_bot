package elastics

import (
	"bot/config"
	zlog "bot/log"
	"context"
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
			zlog.SLogger.Errorf("new es client error: %s", err)
			return
		}
		// Ping the Elasticsearch server to get e.g. the version number
		info, code, err := Client.Ping(e.Url).Do(context.Background())
		if err != nil {
			zlog.SLogger.Errorf("ping es error: %s", err)
			return
		}
		zlog.SLogger.Infof("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

		// Getting the ES version number is quite common, so there's a shortcut
		esversion, err := Client.ElasticsearchVersion(e.Url)
		if err != nil {
			zlog.SLogger.Errorf("get es version error: %s", err)
			return
		}
		zlog.SLogger.Infof("Elasticsearch version %s", esversion)

		// Use the IndexExists service to check if a specified index exists.
		exists, err := Client.IndexExists("status").Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			// create mapping
			createIndex, err := Client.CreateIndex("status").Body(StatusMapping).Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			if !createIndex.Acknowledged {
				zlog.SLogger.Warn("createIndex not acknowledged")
			}
		}

		// Use the IndexExists service to check if a specified index exists.
		lexists, lerr := Client.IndexExists("local").Do(context.Background())
		if lerr != nil {
			log.Fatal(lerr)
		}
		if !lexists {
			// create mapping
			createIndex, err := Client.CreateIndex("local").Body(StatusMapping).Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			if !createIndex.Acknowledged {
				zlog.SLogger.Warn("createIndex not acknowledged")
			}
		}
	})
}
