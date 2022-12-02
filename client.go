package main

import (
	"github.com/go-redis/redis/v8"
)

func newRedisClient(config RedisConfig) *redis.Client {
	if config.ServerConfigs != nil {
		return redis.NewClient(&redis.Options{
			Addr: config.ServerConfigs.Address,
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
