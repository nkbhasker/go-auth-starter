package enum

import "fmt"

type GenderEnum string

const (
	GenderMale   GenderEnum = "MALE"
	GenderFemale GenderEnum = "FEMALE"
)

func (e *GenderEnum) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = GenderEnum(str)

	return nil
}

func (e GenderEnum) Value() (interface{}, error) {
	return string(e), nil
}
