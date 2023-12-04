package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func appAction(cCtx *cli.Context) error {
	argsWithoutProg := cCtx.Args().Slice()
	lenghtOfArgsWithoutProg := len(argsWithoutProg)
	if lenghtOfArgsWithoutProg == 0 {
		fmt.Println("No key patterns provided")
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	fileConfig := getRedisConfig("redisdel.conf")
	redisServerConfig := RedisServerConfigs{
		Address: addr,
	}
	if redisServerConfig.Address == ":" {
		redisServerConfig.Address = fileConfig.ServerConfigs.Address
	}

	redisConfig := RedisConfig{
		Username:        username,
		Password:        password,
		ServerConfigs:   &redisServerConfig,
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
		fmt.Println("Failed to initialise redis client")
		os.Exit(1)
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
		fmt.Printf("Found %d keys for pattern %s", len(keys), argsWithoutProg[i])
		chunkedKeys := lo.Chunk(keys, 1000)
		var wg sync.WaitGroup
		for batch, chunk := range chunkedKeys {
			wg.Add(1)
			go func(chunk []string, batch int) {
				defer wg.Done()
				client.Del(context.Background(), chunk...)
				fmt.Printf("Deleted batch %d of %d for %s", batch, len(chunkedKeys), argsWithoutProg[i])
			}(chunk, batch)
		}
		wg.Wait()
	}
	return nil
}
