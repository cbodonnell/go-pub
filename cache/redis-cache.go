package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cheebz/go-pub/config"
	"github.com/go-redis/redis/v7"
)

type RedisCache struct {
	conf   config.Configuration
	client *redis.Client
}

func NewRedisCache(_conf config.Configuration) Cache {
	return &RedisCache{
		conf: _conf,
		client: redis.NewClient(&redis.Options{
			Addr:     _conf.Redis.Address,
			Password: _conf.Redis.Password,
			DB:       _conf.Redis.Db,
		}),
	}
}

// func (c *RedisCache) getClient() *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr:     c.conf.Redis.Address,
// 		Password: c.conf.Redis.Password,
// 		DB:       c.conf.Redis.Db,
// 	})
// }

func (c *RedisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.client.Set(key, json, time.Duration(c.conf.Redis.RedisExpSeconds)*time.Second)
	log.Println(fmt.Sprintf("set cached %s", key))
	return nil
}

func (c *RedisCache) Get(key string, result interface{}) (interface{}, error) {
	value, err := c.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("got cached %s", key))
	return result, nil
}

func (c *RedisCache) Del(key string) error {
	_, err := c.client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) FlushDB() error {
	_, err := c.client.FlushDB().Result()
	if err != nil {
		return err
	}
	return nil
}
