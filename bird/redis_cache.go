package bird

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisCache(config CacheConfig) (*RedisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisServer,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	cache := &RedisCache{
		client: client,
	}

	return cache, nil
}

func (self *RedisCache) Get(key string) (Parsed, error) {
	key = self.keyPrefix + key //"B" + IPVersion + "_" + key
	data, err := self.client.Get(key).Result()
	if err != nil {
		return NilParse, err
	}

	parsed := Parsed{}
	err = json.Unmarshal([]byte(data), &parsed)

	ttl, correct := parsed["ttl"].(time.Time)
	if !correct {
		return NilParse, errors.New("Invalid TTL value for key" + key)
	}

	if ttl.Before(time.Now()) {
		return NilParse, err // TTL expired
	} else {
		return parsed, err // cache hit
	}
}

func (self *RedisCache) Set(key string, parsed Parsed, ttl int) error {
	switch {
	case ttl == 0:
		return nil // do not cache

	case ttl > 0:
		key = self.keyPrefix + key //TODO "B" + IPVersion + "_" + key
		payload, err := json.Marshal(parsed)
		if err != nil {
			return err
		}

		_, err = self.client.Set(key, payload, time.Duration(ttl)*time.Minute).Result()
		return err

	default: // ttl negative - invalid
		return errors.New("Negative TTL value for key" + key)
	}
}

func (self *RedisCache) Expire() int {
	log.Printf("Cannot expire entries in RedisCache backend, redis does this automatically")
	return 0
}
