package service

type MapMemoryConfig struct {
	Path string
}

type RedisMemoryConfig struct {
	Type    string
	Network string
	Address string
	Pool    int
}
