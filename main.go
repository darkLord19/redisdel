package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
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

func appAction(cCtx *cli.Context) error {
	argsWithoutProg := cCtx.Args().Slice()
	lenghtOfArgsWithoutProg := len(argsWithoutProg)
	if lenghtOfArgsWithoutProg == 0 {
		log.Fatalln("No key patterns provided")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	fileConfig := getRedisConfig("redisdel.conf")
	redisServerConfig := RedisServerConfigs{
		Address: addr,
	}
	if redisServerConfig.Address == "" {
		redisServerConfig.Address = fileConfig.ServerConfigs.Address
	}

	redisConfig := RedisConfig{
		Username: username,
		Password: password,
		ServerConfigs: &redisServerConfig,
		SentinelConfigs: fileConfig.SentinelConfigs,
	}
	if redisConfig.Username == "" {
		redisConfig.Username = fileConfig.Username
	}
	if redisConfig.Password == "" {
		redisConfig.Password = fileConfig.Password
	}

	client := newRedisClient(redisConfig)
	if client == nil {
		log.Fatalln("Failed to initialise redis client")
	}

	matchedKeysChans := make([]chan []string, lenghtOfArgsWithoutProg)
	for i := range matchedKeysChans {
		matchedKeysChans[i] = make(chan []string)
	}

	for i, pattern := range argsWithoutProg {
		go getMatchingKeys(pattern, matchedKeysChans[i], client)
	}

	for i := range matchedKeysChans {
		keys := <-matchedKeysChans[i]
		log.Printf("Found %d keys for pattern %s", len(keys), argsWithoutProg[i])
		chunkedKeys := lo.Chunk(keys, 1000)
		var wg sync.WaitGroup
		for batch, chunk := range chunkedKeys {
			wg.Add(1)
			go func(chunk []string, batch int) {
				defer wg.Done()
				client.Del(context.Background(), chunk...)
				log.Printf("Deleted batch %d of %d for %s", batch, len(chunkedKeys), argsWithoutProg[i])
			}(chunk, batch)
		}
		wg.Wait()
	}
	return nil
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
		Name:  "redisdel",
		Usage: "Scan for redis keys matching a given pattern and delete them",
		Version: "0.1.0",
		Action: appAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
					Name:  "host",
					Usage: "redis host",
					Destination: &host,
					Required: true,
			},
			&cli.StringFlag{
				Name:  "port",
				Usage: "redis port",
				Destination: &port,
				Required: true,
			},
			&cli.StringFlag{
				Name:  "username",
				Usage: "redis username",
				Destination: &username,
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "redis password",
				Destination: &password,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
