package repo

import (
	"github.com/nkbhasker/go-auth-starter/internal/misc"
	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
)

type Repo interface {
	UserRepo() UserRepo
	AuthRepo() AuthRepo
	AccessTokenRepo() AccessTokenRepo
}

type repo struct {
	userRepo        UserRepo
	authRepo        AuthRepo
	accessTokenRepo AccessTokenRepo
}

type RepoOptions struct {
	DBStore                    storage.DBStore
	CacheStore                 storage.CacheStore
	IdGenerator                uid.IdGenerator
	JwtHelper                  misc.JwtHelper
	AccessTokenExpiryInMinutes int
	OtpExpiryInMinutes         int
}

func NewRepo(options RepoOptions) Repo {
	return &repo{
		userRepo:        NewUserRepo(options.DBStore, options.IdGenerator),
		authRepo:        NewAuthRepo(options.DBStore, options.CacheStore, options.IdGenerator, options.OtpExpiryInMinutes),
		accessTokenRepo: NewAccessToeknRepo(options.CacheStore, options.JwtHelper, options.AccessTokenExpiryInMinutes),
	}
}

func (r repo) UserRepo() UserRepo {
	return r.userRepo
}

func (r repo) AuthRepo() AuthRepo {
	return r.authRepo
}

func (r repo) AccessTokenRepo() AccessTokenRepo {
	return r.accessTokenRepo
}
