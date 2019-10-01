package memory

import (
	"encoding/json"
	"errors"
	"github.com/DeathHand/gateway/model"
	"github.com/gomodule/redigo/redis"
	"strings"
	"time"
)

type RedisMemory struct {
	Memory
	pool    *redis.Pool
	notify  *chan model.Message
	errChan *chan error
}

func NewRedisMemory(network string, addr string, pool int, errors *chan error) (*RedisMemory, error) {
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
		errChan: errors,
	}, nil
}

func (r *RedisMemory) Put(message *model.Message) error {
	m, err := json.Marshal(message)
	if err != nil {
		return err
	}
	conn := r.pool.Get()
	_, err = conn.Do("SETEX", message.Uuid, m, message.Ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisMemory) Get(uuid string) (*model.Message, error) {
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

func (r *RedisMemory) Delete(uuid string) error {
	conn := r.pool.Get()
	_, err := conn.Do("DELETE", uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisMemory) Run() {
	conn := r.pool.Get()
	for {
		expired, err := redis.String(
			conn.Do("SUBSCRIBE", "__keyevent@0__:expired"),
		)
		if err != nil {
			*r.errChan <- err
			return
		}
		ev := strings.Split(expired, ":")
		if ev[1] == "" {
			*r.errChan <- errors.New("Unknown key expired ")
		}
		data, err := redis.Bytes(conn.Do("GET", ev[1]))
		if err != nil {
			*r.errChan <- err
			return
		}
		message := model.Message{}
		err = json.Unmarshal(data, &message)
		if err != nil {
			*r.errChan <- err
			return
		}
		*r.notify <- message
	}
}

func (r *RedisMemory) Notify() *chan model.Message {
	return r.notify
}
