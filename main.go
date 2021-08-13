package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func getKeysMatchingPattern(pattern string, matchedKeys chan []string) {
	ctx := context.Background()
	var cursor uint64
	var keys []string
	for {
		var matchedKeysSoFar []string
		var err error
		matchedKeysSoFar, cursor, err = redisClient.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			log.Println("Failed to get matching keys for pattern", pattern)
			return
		}
		if cursor == 0 {
			break
		}
		keys = append(keys, matchedKeysSoFar...)
	}
	matchedKeys <- keys
}

func main() {
	argsWithoutProg := os.Args[1:]
	lenghtOfArgsWithoutProg := len(argsWithoutProg)
	if lenghtOfArgsWithoutProg == 0 {
		log.Println("Please provide search patterns")
		os.Exit(1)
	}
	matchedKeysChans := make([]chan []string, lenghtOfArgsWithoutProg)
	for i := range matchedKeysChans {
		matchedKeysChans[i] = make(chan []string)
	}
	for i, pattern := range argsWithoutProg {
		go getKeysMatchingPattern(pattern, matchedKeysChans[i])
	}
	for i := range matchedKeysChans {
		keys := <-matchedKeysChans[i]
		log.Println(len(keys))
	}
}
