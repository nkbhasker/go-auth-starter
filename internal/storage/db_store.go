package storage

import (
	"github.com/nkbhasker/go-auth-starter/internal/health"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBStore interface {
	health.Checker
	DB() *gorm.DB
	CloseDB() error
	WithTx(tx *gorm.DB) DBStore
}

type dbStore struct {
	db *gorm.DB
}

func init() {
	registerSerializers()
}

func InitDBStore(postgresUrl string) (DBStore, error) {
	db, err := gorm.Open(postgres.Open(postgresUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &dbStore{db: db}, nil
}

func (s dbStore) DB() *gorm.DB {
	return s.db
}

func (s *dbStore) CloseDB() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s dbStore) WithTx(tx *gorm.DB) DBStore {
	return &dbStore{
		db: tx,
	}
}

func (s *dbStore) Check() *health.Health {
	var version string
	h := health.NewHealth()
	if err := s.db.Raw("SELECT VERSION();").Scan(&version).Error; err != nil {
		h.SetStatus(health.HealthStatusDown)
		h.SetInfo("error", err.Error())
	} else {
		h.SetStatus(health.HealthStatusUp)
		h.SetInfo("version", version)
	}

	return h
}

func registerSerializers() {
	schema.RegisterSerializer("id", uid.NewIdSerializer())
}
