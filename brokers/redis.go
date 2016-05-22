package brokers

import (
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	config "github.com/shidenggui/gotq/config"
)

type RedisBroker struct {
	Pool *redis.Pool
}

func (r *RedisBroker) Delay(jsonByte []byte, queue string) error {
	c := r.Pool.Get()
	defer c.Close()
	_, err := c.Do("LPUSH", queue, jsonByte)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisBroker) Expire(key string, expireTime int64) error {
	c := r.Pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", key, expireTime)
	if err != nil {
		return err
	}
	return nil
}

func NewRedisBroker(cfg *config.BrokerCfg) *RedisBroker {
	return &RedisBroker{
		Pool: NewRedisPool(cfg.Host, cfg.Port, cfg.Password, cfg.DB),
	}
}

func (r *RedisBroker) Request(queue string, blockTime int64) ([]byte, error) {
	c := r.Pool.Get()
	defer c.Close()

	taskSlice, err := c.Do("BLPOP", queue, blockTime)
	if err != nil {
		return nil, err
	}

	taskPairs, err := redis.ByteSlices(taskSlice, nil)
	if err != nil {
		return nil, err
	}

	if len(taskPairs) != 2 {
		return nil, err
	}

	return taskPairs[1], nil
}

func (r *RedisBroker) Receive(queue string) ([]byte, error) {
	c := r.Pool.Get()
	defer c.Close()

	taskSlice, err := c.Do("BLPOP", queue, "0")
	if err != nil {
		return nil, err
	}

	taskPairs, err := redis.ByteSlices(taskSlice, nil)
	if err != nil {
		return nil, err
	}

	if len(taskPairs) != 2 {
		return nil, err
	}

	return taskPairs[1], nil
}

func NewRedisPool(host string, port int64, password string, DB int64) *redis.Pool {
	server := host + ":" + strconv.Itoa(int(port))
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			_, err = c.Do("SELECT", DB)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}
