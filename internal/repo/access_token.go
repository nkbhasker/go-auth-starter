package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/nkbhasker/go-auth-starter/internal/misc"
	"github.com/nkbhasker/go-auth-starter/internal/storage"
)

const accessTokenKey = "atk"

type AccessTokenRepo interface {
	Create(sub string, name string) (string, error)
}

type accessTokenRepo struct {
	cacheStore storage.CacheStore
	jwtHelper  misc.JwtHelper
	ttl        time.Duration
}

func (r *accessTokenRepo) Create(sub string, name string) (string, error) {
	id, accessToken, err := r.jwtHelper.NewAccessToken(sub, name, r.ttl)
	if err != nil {
		return "", err
	}
	err = r.addToken(id, sub)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func NewAccessToeknRepo(cacheStore storage.CacheStore, jwtHelper misc.JwtHelper, expiresInMinutes int) AccessTokenRepo {
	return &accessTokenRepo{
		cacheStore: cacheStore,
		jwtHelper:  jwtHelper,
		ttl:        time.Duration(expiresInMinutes * int(time.Minute)),
	}
}

func (r *accessTokenRepo) addToken(jti, sub string) error {
	key := fmt.Sprintf("%s_%s_%s", accessTokenKey, sub, jti)
	return r.cacheStore.WithTTL(r.ttl).Set(context.Background(), key, "1")
}
