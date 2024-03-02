package misc

import (
	"crypto/rand"
)

const (
	length = 6
	digits = "0123456789"
)

func GenerateOtp() (string, error) {
	max := len(digits) - 1
	str := make([]byte, length)
	randBytes := make([]byte, length)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	for i, rb := range randBytes {
		str[i] = digits[int(rb)%max]
	}

	return string(str), nil
}

func ValidateOtp(want, got string) bool {
	if want == "" || got == "" {
		return false
	}

	return want == got
}
