package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

var host string
var port string

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
		Action: func(cCtx *cli.Context) error {
			argsWithoutProg := cCtx.Args().Slice()
			lenghtOfArgsWithoutProg := len(argsWithoutProg)
			if lenghtOfArgsWithoutProg == 0 {
				log.Fatalln("No key patterns provided")
			}

			client := newRedisClient("redisdel.conf")
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
				log.Printf("Found %s keys for pattern %d", argsWithoutProg[i], len(keys))
				chunkedKeys := lo.Chunk(keys, 1000)
				var wg sync.WaitGroup
				for batch, chunk := range chunkedKeys {
					wg.Add(1)
					go func(chunk []string, batch int) {
						defer wg.Done()
						client.Del(context.Background(), chunk...)
						log.Printf("Deleted batch %d of %d for %s", batch+1, len(chunkedKeys), argsWithoutProg[i])
					}(chunk, batch)
				}
				wg.Wait()
			}
			return nil
	},
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
	},
}

if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
