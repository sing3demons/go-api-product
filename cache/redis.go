package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type ProductCache interface {
	Set(key string, value interface{})
	Get(key string) (string, error)
	GetPage(key string) *pagingResult
}

func NewRedisCache(host string, db int, exp time.Duration) ProductCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) GetPage(key string) *pagingResult {
	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	page := pagingResult{}
	err = json.Unmarshal([]byte(val), &page)
	if err != nil {
		panic(err)
	}

	return &page
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value interface{}) {
	client := cache.getClient()

	// serialize value object to JSON
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", key)

	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) (string, error) {
	client := cache.getClient()

	resp, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return resp, nil
}
