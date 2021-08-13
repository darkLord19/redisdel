package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

type RedisConfig struct {
	Address  string
	Username string
	Password string
	DB       uint8
}

func init() {
	config := getRedisConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		Username: config.Username,
		DB:       int(config.DB),
	})
}

func getRedisConfig() *RedisConfig {
	fileName := "redisdel.conf"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("Failed to read redis del config")
		os.Exit(1)
	}
	var config RedisConfig
	json.Unmarshal(data, &config)
	return &config
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
			log.Println("Failed to get matching keys for pattern", pattern, err)
			close(matchedKeys)
			return
		}
		if cursor == 0 {
			break
		}
		keys = append(keys, matchedKeysSoFar...)
	}
	matchedKeys <- keys
	close(matchedKeys)
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
