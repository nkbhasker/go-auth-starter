package storage

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/nkbhasker/go-auth-starter/internal/health"
	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
	health.Checker
	DB() *redis.Client
	CloseDB() error
	WithTTL(ttl time.Duration) CacheStore
	WithKeepTTL() CacheStore
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) error
}

type cacheStore struct {
	db  *redis.Client
	ttl time.Duration
}

func InitCacheStore(redisUrl string) (CacheStore, error) {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &cacheStore{db: client}, nil
}

func (s cacheStore) DB() *redis.Client {
	return s.db
}

func (s *cacheStore) CloseDB() error {
	return s.db.Close()
}

func (s *cacheStore) Set(ctx context.Context, key string, value interface{}) error {
	v, err := toValue(value)
	if err != nil {
		return err
	}

	return s.db.Set(ctx, key, v, s.ttl).Err()
}

func (s *cacheStore) Get(ctx context.Context, key string) error {
	return s.db.Get(ctx, key).Err()
}

func (s cacheStore) WithTTL(ttl time.Duration) CacheStore {
	return &cacheStore{
		db:  s.db,
		ttl: ttl,
	}
}

func (s cacheStore) WithKeepTTL() CacheStore {
	return &cacheStore{
		db:  s.db,
		ttl: redis.KeepTTL,
	}
}

func (s *cacheStore) Check() *health.Health {
	h := health.NewHealth()
	res := s.db.InfoMap(context.Background(), "server")
	if res.Err() != nil {
		h.SetStatus(health.HealthStatusDown)
		h.SetInfo("error", res.Err().Error())
	} else {
		h.SetStatus(health.HealthStatusUp)
		h.SetInfo("version", res.Val()["Server"]["redis_version"])
	}

	return h
}

func toValue(value interface{}) (interface{}, error) {
	k := reflect.Indirect(reflect.ValueOf(value)).Kind()
	if k == reflect.Struct {
		return json.Marshal(value)
	}

	return value, nil
}
