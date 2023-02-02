package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var host string
var port string
var username string
var password string

func getRedisConfig(cfgFile string) *RedisConfig {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalln("Failed to read redis del config")
	}
	var config RedisConfig
	json.Unmarshal(data, &config)
	return &config
}

func getMatchingKeys(pattern string, matchedKeys chan []string, client *redis.Client) {
	var cursor uint64
	var keys []string
	hasNextPage := true

	for hasNextPage {
		var matchedKeysSoFar []string
		var err error
		matchedKeysSoFar, cursor, err = client.Scan(context.TODO(), cursor, pattern, 1000).Result()
		if err != nil {
			log.Printf("Failed to get matching keys for pattern: %s with err:%v", pattern, err)
			close(matchedKeys)
			return
		}
		keys = append(keys, matchedKeysSoFar...)
		hasNextPage = cursor != 0
	}

	matchedKeys <- keys
	close(matchedKeys)
}

func main() {
	app := &cli.App{
		Name:    "redisdel",
		Usage:   "Scan for redis keys matching a given pattern and delete them",
		Version: "0.1.0",
		Action:  appAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "redis host",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "port",
				Usage:       "redis port",
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "username",
				Usage:       "redis username",
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "password",
				Usage:       "redis password",
				Destination: &password,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
