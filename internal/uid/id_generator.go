package uid

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sony/sonyflake"
)

var startTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type KindEnum string
type FieldNameEnum string

const (
	KindUser    KindEnum      = "usr"
	FieldNameID FieldNameEnum = "ID"
)

type IdGenerator interface {
	Next(kind KindEnum) (Identifier, error)
	NextFromFieldTag(value interface{}, fieldName FieldNameEnum) (Identifier, error)
}

type idGenerator struct {
	*sonyflake.Sonyflake
}

func NewIdGenerator() IdGenerator {
	return &idGenerator{
		sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
		}),
	}
}

func (i *idGenerator) Next(kind KindEnum) (Identifier, error) {
	uid, err := i.NextID()
	if err != nil {
		return nil, err
	}

	return NewIdentifier(
		kind,
		uid,
		timestamp(uid),
	), nil
}

func (i *idGenerator) NextFromFieldTag(value interface{}, fieldName FieldNameEnum) (Identifier, error) {
	kind, err := kind(value, fieldName)
	if err != nil {
		return nil, err
	}

	return i.Next(kind)
}

func Timestamp(uid int64) time.Time {
	return timestamp(uint64(uid))
}

func timestamp(uid uint64) time.Time {
	elapsedTime := sonyflake.ElapsedTime(uid)
	timestamp := startTime.Add(elapsedTime)

	return timestamp
}

func kind(value interface{}, fieldName FieldNameEnum) (KindEnum, error) {
	k := reflect.Indirect(reflect.ValueOf(value)).Kind()
	if k != reflect.Struct {
		return "", fmt.Errorf("invalid type")
	}
	field, ok := reflect.TypeOf(value).FieldByName(string(fieldName))
	if !ok {
		return "", fmt.Errorf("invalid field")
	}
	idKind, ok := field.Tag.Lookup("kind")
	if !ok {
		return "", fmt.Errorf("invalid kind")
	}

	return KindEnum(idKind), nil
}
