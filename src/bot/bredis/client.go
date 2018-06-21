package bredis

import (
	"log"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
