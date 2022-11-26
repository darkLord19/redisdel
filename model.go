package main

type RedisServerConfigs struct {
	Address string
}

type RedisSentinelConfigs struct {
	MasterName string
	Password   string
	Addresses  []string
}

type RedisConfig struct {
	Username        string
	Password        string
	ServerConfigs   *RedisServerConfigs
	SentinelConfigs *RedisSentinelConfigs
}
