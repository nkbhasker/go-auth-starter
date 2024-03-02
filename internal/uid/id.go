package uid

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var _ Identifier = (*identifier)(nil)

type Identifier interface {
	driver.Valuer
	schema.GormDataTypeInterface
	migrator.GormDataTypeInterface
	Kind() KindEnum
	Uid() int64
	Timestamp() time.Time
	String() string
}

type identifier struct {
	kind      KindEnum
	uid       int64
	timestamp time.Time
}

func NewIdentifier(kind KindEnum, uid uint64, timestamp time.Time) Identifier {
	return &identifier{
		kind:      kind,
		uid:       int64(uid),
		timestamp: timestamp,
	}
}

func (id identifier) Kind() KindEnum {
	return id.kind
}

func (id identifier) Uid() int64 {
	return id.uid
}

func (id identifier) Timestamp() time.Time {
	return id.timestamp
}

func (id identifier) String() string {
	return fmt.Sprintf("%s_%d", id.kind, id.uid)
}

func (id identifier) Value() (driver.Value, error) {
	return id.uid, nil
}

func (id *identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *identifier) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	_id, err := FromIdString(string(data))
	if err != nil {
		return err
	}
	*id = *_id

	return nil
}

func (identifier) GormDataType() string {
	return "bigint"
}

func (identifier) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "bigint"
}

func FromIdString(str string) (*identifier, error) {
	idParts := strings.Split(str, "_")
	if len(idParts) < 2 {
		return nil, errors.New("invalid str")
	}

	kind := idParts[0]
	uid, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	return FromUid(KindEnum(kind), uid), nil
}

func FromUid(kind KindEnum, uid int64) *identifier {
	timestamp := Timestamp(uid)
	return &identifier{
		uid:       uid,
		kind:      kind,
		timestamp: timestamp,
	}
}
