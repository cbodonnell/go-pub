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
	conf config.Configuration
}

func NewRedisCache(_conf config.Configuration) Cache {
	return &RedisCache{
		conf: _conf,
	}
}

func (c *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.conf.Redis.Address,
		Password: c.conf.Redis.Password,
		DB:       c.conf.Redis.Db,
	})
}

func (c *RedisCache) Set(key string, value interface{}) error {
	client := c.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	client.Set(key, json, time.Duration(c.conf.Redis.RedisExpSeconds)*time.Second)
	log.Println(fmt.Sprintf("set cached %s", key))
	return nil
}

func (c *RedisCache) Get(key string, result interface{}) (interface{}, error) {
	client := c.getClient()
	value, err := client.Get(key).Result()
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
	client := c.getClient()
	_, err := client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
