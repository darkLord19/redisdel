package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisConfig(cfgFile string) *RedisConfig {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalln("Failed to read redis del config")
	}
	var config RedisConfig
	json.Unmarshal(data, &config)
	return &config
}

func newRedisClient(cfgFile string, Addr string) *redis.Client {
	config := getRedisConfig(cfgFile)
	if config.ServerConfigs != nil {
		return redis.NewClient(&redis.Options{
			Addr: Addr,
			Password: config.Password,
			Username: config.Username,
		})
	} else if config.SentinelConfigs != nil {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       config.SentinelConfigs.MasterName,
			SentinelAddrs:    config.SentinelConfigs.Addresses,
			SentinelPassword: config.SentinelConfigs.Password,
			Username:         config.Username,
			Password:         config.Password,
		})
	}
	return nil
}
