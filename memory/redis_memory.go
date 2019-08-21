package memory

import (
	"github.com/DeathHand/gateway/model"
	"github.com/mediocregopher/radix/v3"
)

type RedisMemory struct {
	Memory
	pool *radix.Pool
}

func NewRedisMemory(network string, addr string, pool int) (*RedisMemory, error) {
	p, err := radix.NewPool(network, addr, pool)

	if err != nil {
		return nil, err
	}

	return &RedisMemory{pool: p}, nil
}

func (r *RedisMemory) Put(message *model.Message) (string, error) {
	return "", nil
}

func (r *RedisMemory) Get(id string) (*model.Message, error) {
	return &model.Message{}, nil
}

func (r *RedisMemory) Delete(id string) error {
	return nil
}

func (r *RedisMemory) Observe() error {
	return nil
}
