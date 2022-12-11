package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cheebz/go-pub/pkg/config"
	"github.com/go-redis/redis/v7"
)

type RedisCache struct {
	conf   config.Configuration
	client *redis.Client
}

func NewRedisCache(_conf config.Configuration) Cache {
	return &RedisCache{
		conf:   _conf,
		client: createClient(_conf.Redis),
	}
}

func createClient(conf config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Db,
	})
}

func (c *RedisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.client.Set(key, json, time.Duration(c.conf.Redis.RedisExpSeconds)*time.Second)
	log.Println(fmt.Sprintf("set cached %s", key))
	return nil
}

func (c *RedisCache) Get(key string, pointer interface{}) (interface{}, error) {
	value, err := c.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(value), pointer)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("got cached %s", key))
	return pointer, nil
}

func (c *RedisCache) Del(patterns ...string) error {
	for _, pattern := range patterns {
		keys, err := c.client.Keys(pattern).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			_, err = c.client.Del(keys...).Result()
			if err != nil {
				return err
			}
		}
	}
	log.Println(fmt.Sprintf("deleted cached %s", strings.Join(patterns, " ")))
	return nil
}

func (c *RedisCache) FlushDB() error {
	_, err := c.client.FlushDB().Result()
	if err != nil {
		return err
	}
	return nil
}
