package cache

import (
	"app/models"
	"encoding/json"
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
	Get(key string) []models.Product
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

	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) []models.Product {
	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	product := []models.Product{}
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		panic(err)
	}

	return product
}
