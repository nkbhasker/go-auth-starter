package uid

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm/schema"
)

type IdSerializer interface {
	schema.SerializerInterface
}

type idSerializer struct{}

func NewIdSerializer() IdSerializer {
	return &idSerializer{}
}

func (idSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return errors.New("nil value")
	}
	dbId, ok := dbValue.(int64)
	if !ok {
		return errors.New("invalid type")
	}

	id := FromUid(KindEnum(field.Tag.Get("kind")), dbId)
	err := field.Set(ctx, dst, id)
	if err != nil {
		return err
	}

	return nil
}

func (idSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	v, ok := fieldValue.(Identifier)
	if !ok {
		return nil, errors.New("invalid identifier")
	}

	return v.Uid(), nil
}
