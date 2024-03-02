package enum

import "fmt"

type IdentityProviderEnum string

const (
	IdentityProviderLocal  IdentityProviderEnum = "LOCAL"
	IdentityProviderGoogle IdentityProviderEnum = "GOOGLE"
	IdentityProviderApple  IdentityProviderEnum = "APPLE"
)

func (e *IdentityProviderEnum) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = IdentityProviderEnum(str)

	return nil
}

func (e IdentityProviderEnum) Value() (interface{}, error) {
	return string(e), nil
}
