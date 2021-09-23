package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
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
	MSet(kv map[string]interface{}) error
	Get(key string) (string, error)
	MGet(keys []string) ([]interface{}, error)
	Del(keys ...string) error
	Expires(keys []string, expire time.Duration) error
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

	ctx, cancel := context.WithTimeout(context.Background(), cache.expires*time.Second)
	defer cancel()

	val, err := client.Get(ctx, key).Result()
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

func (cache *redisCache) MGet(keys []string) ([]interface{}, error) {
	client := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), cache.expires*time.Second)
	defer cancel()

	val, err := client.MGet(ctx, keys...).Result()
	// fmt.Printf("get key: %v\n", keys)
	// fmt.Println("get...")
	if err == redis.Nil {
		// Key does not exists
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return val, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), cache.expires*time.Second)
	defer cancel()

	// serialize value object to JSON
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", key)

	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *redisCache) MSet(kv map[string]interface{}) error {
	client := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), cache.expires*time.Second)
	defer cancel()

	pairs := []interface{}{}
	for k, v := range kv {

		str, ok := v.(string)
		// Check empty string if value string
		if ok && len(str) == 0 {
			pairs = append(pairs, k, "")
			continue
		}
		// If value is string, not pass it to json.Marshal
		if len(str) > 0 {
			pairs = append(pairs, k, str)
			continue
		}

		strb, err := json.Marshal(v)
		if err != nil {
			return err
		}

		pairs = append(pairs, k, strb)
	}

	// fmt.Printf("set...%v", pairs...)

	err := client.MSet(ctx, pairs...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cache *redisCache) Get(key string) (string, error) {
	client := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), cache.expires*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return resp, nil
}

// Del the cache by keys
func (cache *redisCache) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	c := cache.getClient()

	// Delete 10000 items per page
	pageLimit := 10000
	from := 0
	to := pageLimit

	for {
		// Lower bound
		if from >= len(keys) {
			break
		}
		// Upper bound
		if to > len(keys) {
			to = len(keys)
		}

		delKeys := keys[from:to]
		if len(delKeys) == 0 {
			break
		}

		_, err := c.Del(context.Background(), delKeys...).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			} else {
				return err
			}
		}
		from += pageLimit
		to += pageLimit
	}

	return nil
}

func (cache *redisCache) Expires(keys []string, expire time.Duration) error {
	c := cache.getClient()

	var lastErr error
	for _, key := range keys {
		err := c.Expire(context.Background(), key, expire).Err()
		if err != nil {
			if err == redis.Nil {
				// Key does not exists
				return nil
			} else {
				lastErr = err
			}
		}
	}
	return lastErr
}
