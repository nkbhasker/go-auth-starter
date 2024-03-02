package repo

import (
	"fmt"

	"github.com/nkbhasker/go-auth-starter/internal/model"
	"github.com/nkbhasker/go-auth-starter/internal/storage"
	"github.com/nkbhasker/go-auth-starter/internal/uid"
	"gorm.io/gorm"
)

var ErrUserNotFound = fmt.Errorf("user not found")

type UserRepo interface {
	New(options model.User) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Get(id uid.Identifier) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	WithTx(tx *gorm.DB) UserRepo
}

type userRepo struct {
	dbStore     storage.DBStore
	idGenerator uid.IdGenerator
}

func NewUserRepo(dbStore storage.DBStore, idGenerator uid.IdGenerator) UserRepo {
	return &userRepo{
		dbStore:     dbStore,
		idGenerator: idGenerator,
	}
}

func (r userRepo) WithTx(tx *gorm.DB) UserRepo {
	return NewUserRepo(r.dbStore.WithTx(tx), r.idGenerator)
}

func (r userRepo) New(options model.User) (*model.User, error) {
	if options.ID == nil {
		id, err := r.idGenerator.NextFromFieldTag(options, uid.FieldNameID)
		if err != nil {
			return nil, err
		}
		options.ID = id
	}

	return &options, nil
}

func (r userRepo) Create(user *model.User) error {
	return r.dbStore.DB().Save(user).Error
}

func (r userRepo) Get(id uid.Identifier) (*model.User, error) {
	user := &model.User{}
	err := r.dbStore.DB().Find(user, id).Error
	if err != nil {
		return nil, nil
	}
	if user.ID == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r userRepo) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.dbStore.DB().Where(`"email"= ?`, email).Find(user).Error
	if err != nil {
		return nil, nil
	}
	if user.ID == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r userRepo) Update(user *model.User) error {
	return r.dbStore.DB().Model(user).Updates(user).Error
}
