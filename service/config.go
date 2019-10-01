package service

type Config struct {
	IngressSize   int
	EgressSize    int
	RedisMemory   *RedisMemoryConfig
	MapMemory     *MapMemoryConfig
	HttpCallback  *HttpCallbackConfig
	RedisCallback *RedisCallbackConfig
}
