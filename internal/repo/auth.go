package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
	"github.com/redis/go-redis/v9"
)

type AuthKeyEnum string

const (
	AuthKeyOTP AuthKeyEnum = "OTP"
)

type AuthRepo interface {
	SaveOTP(ctx context.Context, key string, otp string) error
	GetOTP(ctx context.Context, key string) (string, error)
}

type authRepo struct {
	dbStore      storage.DBStore
	cacheStore   storage.CacheStore
	idGenerator  uid.IdGenerator
	otpExpiresIn time.Duration
}

func NewAuthRepo(
	dbStore storage.DBStore,
	cacheStore storage.CacheStore,
	idGenerator uid.IdGenerator,
	otpExpiryInMintues int,
) AuthRepo {
	return &authRepo{
		dbStore:      dbStore,
		cacheStore:   cacheStore,
		idGenerator:  idGenerator,
		otpExpiresIn: time.Duration(otpExpiryInMintues * int(time.Minute)),
	}
}

func (r *authRepo) SaveOTP(ctx context.Context, key string, otp string) error {
	return r.cacheStore.DB().Set(ctx, strings.ToLower(key), otp, r.otpExpiresIn).Err()
}

func (r *authRepo) GetOTP(ctx context.Context, key string) (string, error) {
	result := r.cacheStore.DB().Get(ctx, strings.ToLower(key))
	if result.Err() == redis.Nil {
		return "", fmt.Errorf("otp expired")
	}
	if result.Err() != nil {
		return "", result.Err()
	}

	return result.Val(), nil
}
