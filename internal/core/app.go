package core

import (
	"runtime"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nkbhasker/go-auth-starter/internal/comm"
	"github.com/nkbhasker/go-auth-starter/internal/health"
	"github.com/nkbhasker/go-auth-starter/internal/repo"
	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
)

type App interface {
	health.Checker
	Version() string
	DBStore() storage.DBStore
	CacheStore() storage.CacheStore
	Repo() repo.Repo
	IdGenerator() uid.IdGenerator
	Emailer() comm.Emailer
	Validate() *validator.Validate
}

type app struct {
	startAt     time.Time
	version     string
	dbStore     storage.DBStore
	cacheStore  storage.CacheStore
	repos       repo.Repo
	idGenerator uid.IdGenerator
	emailer     comm.Emailer
	validate    *validator.Validate
}

type AppOption struct {
	Version     string
	DBStore     storage.DBStore
	CacheStore  storage.CacheStore
	Repos       repo.Repo
	IdGenerator uid.IdGenerator
	Emailer     comm.Emailer
	Validate    *validator.Validate
}

func NewApp(options AppOption) App {
	return &app{
		startAt:     time.Now(),
		version:     options.Version,
		dbStore:     options.DBStore,
		cacheStore:  options.CacheStore,
		repos:       options.Repos,
		idGenerator: options.IdGenerator,
		emailer:     options.Emailer,
		validate:    options.Validate,
	}
}

func (a *app) Version() string {
	return a.version
}

func (a *app) CacheStore() storage.CacheStore {
	return a.cacheStore
}

func (a *app) DBStore() storage.DBStore {
	return a.dbStore
}

func (a *app) Repo() repo.Repo {
	return a.repos
}

func (a *app) IdGenerator() uid.IdGenerator {
	return a.idGenerator
}

func (a *app) Validate() *validator.Validate {
	return a.validate
}

func (a *app) Emailer() comm.Emailer {
	return a.emailer
}

func (a *app) Check() *health.Health {
	h := health.NewHealth()
	h.SetStatus(health.HealthStatusUp)
	h.SetInfo("version", a.Version())
	h.SetInfo("uptime", time.Since(a.startAt).String())
	h.SetInfo("cpus", runtime.NumCPU())
	h.SetInfo("os", runtime.GOOS)

	return h
}
