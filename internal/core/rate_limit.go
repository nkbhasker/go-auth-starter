package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/redis/go-redis/v9"
)

type RateLimiterKindEnum string

const (
	RateLimiterKindOtpVerify   RateLimiterKindEnum = "OTP_VERIFY"
	RateLimiterKindOtpGenerate RateLimiterKindEnum = "OTP_GENERATE"
)

type RateLimiter interface {
	Evaluate(identifier string) (bool, error)
	Reset(identifier string) error
}

type rateLimiter struct {
	cacheStore storage.CacheStore
	kind       RateLimiterKindEnum
	limit      int
	window     time.Duration
}

func NewRateLimiter(
	cacheStore storage.CacheStore,
	kind RateLimiterKindEnum,
	limit int,
	timeWindowInSeconds int,
) RateLimiter {
	return &rateLimiter{
		cacheStore: cacheStore,
		kind:       kind,
		limit:      limit,
		window:     time.Duration(timeWindowInSeconds * int(time.Second)),
	}
}

func (r *rateLimiter) Evaluate(identifier string) (bool, error) {
	ctx := context.Background()
	timestamp := time.Now()
	key := strings.ToLower(fmt.Sprintf(`%s_%s`, r.kind, identifier))
	pipe := r.cacheStore.DB().TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0.0", strconv.FormatFloat(float64(timestamp.Add(-r.window).UnixMilli()), 'f', -1, 64))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(timestamp.UnixMilli()), Member: float64(timestamp.UnixMilli())})
	results := pipe.ZRange(ctx, key, 0, -1)
	pipe.Expire(ctx, key, r.window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	if len(results.Val()) > r.limit {
		return false, nil
	}

	return true, nil
}

func (r *rateLimiter) Reset(identifier string) error {
	ctx := context.Background()
	key := strings.ToLower(fmt.Sprintf(`%s_%s`, r.kind, identifier))
	err := r.cacheStore.DB().Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
