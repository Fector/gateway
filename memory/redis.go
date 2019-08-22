package memory

import (
	"encoding/json"
	"errors"
	"github.com/DeathHand/gateway/model"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"strings"
	"time"
)

type RedisMemory struct {
	Memory
	pool   *redis.Pool
	errors *chan error
}

func NewRedisMemory(network string, addr string, pool int) (*RedisMemory, error) {
	e := make(chan error, 1)
	return &RedisMemory{
		pool: &redis.Pool{
			MaxActive:   pool,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		errors: &e,
	}, nil
}

func (r *RedisMemory) Put(message *model.Message) (string, error) {
	m, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	id := u.String()
	conn := r.pool.Get()
	_, err = conn.Do("SETEX", id, m, message.Ttl)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *RedisMemory) Get(id string) (*model.Message, error) {
	conn := r.pool.Get()
	data, err := redis.Bytes(conn.Do("GET"))
	if err != nil {
		return nil, err
	}
	message := model.Message{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *RedisMemory) Delete(id string) error {
	conn := r.pool.Get()
	_, err := conn.Do("DELETE", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisMemory) Observe() {
	conn := r.pool.Get()
	expired, err := redis.String(
		conn.Do("SUBSCRIBE", "__keyevent@0__:expired"),
	)
	if err != nil {
		*r.errors <- err
		return
	}
	ev := strings.Split(expired, ":")
	if ev[1] == "" {
		*r.errors <- errors.New("Unknown key expired ")
	}
	data, err := redis.Bytes(conn.Do("GET", ev[1]))
	if err != nil {
		*r.errors <- err
		return
	}
	message := model.Message{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		*r.errors <- err
		return
	}
}
