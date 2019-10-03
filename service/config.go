package service

import "github.com/DeathHand/gateway/model"

type Config struct {
	IngressSize   int
	EgressSize    int
	RedisMemory   RedisMemoryConfig
	MapMemory     MapMemoryConfig
	HttpCallback  HttpCallbackConfig
	RedisCallback RedisCallbackConfig
	Gateways      []model.Gateway
}
