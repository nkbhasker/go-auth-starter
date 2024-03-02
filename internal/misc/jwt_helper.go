package misc

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtHelper interface {
	NewAccessToken(sub string, name string, expiresIn time.Duration) (string, string, error)
	VerifyAccessToken(accessToken string) (*claims, error)
}

type jwtHelper struct {
	issuer string
	key    *rsa.PrivateKey
	kid    string
}

type claims struct {
	jwt.RegisteredClaims
	Name string `json:"name"`
}

func NewJwtHelper(issuer string, base64Str string) (JwtHelper, error) {
	key, err := Base64ToPrivateKey(base64Str)
	if err != nil {
		return nil, err
	}
	bytes, err := PublicKeyToBytes(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	_, err = h.Write(bytes)
	if err != nil {
		return nil, err
	}
	kid := hex.EncodeToString(h.Sum(nil))

	return &jwtHelper{
		issuer: issuer,
		key:    key,
		kid:    kid,
	}, nil
}

func (j *jwtHelper) NewAccessToken(sub string, name string, expiresIn time.Duration) (string, string, error) {
	registeredClaims := jwt.RegisteredClaims{
		ID:       uuid.NewString(),
		Issuer:   j.issuer,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  sub,
	}
	if expiresIn != 0 {
		registeredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims{
		Name:             name,
		RegisteredClaims: registeredClaims,
	})
	token.Header["kid"] = j.kid

	// Create the JWT string.
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", "", err
	}

	return registeredClaims.ID, tokenString, nil
}

func (j *jwtHelper) VerifyAccessToken(accessToken string) (*claims, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Name {
			return nil, fmt.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodRS256)
		}
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected kid")
		}
		if kid != j.kid {
			return nil, fmt.Errorf("unexpectd kid")
		}

		// Return public key pointer expected by rsa verify
		return &j.key.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("unexpected access token")
	}

	return claims, nil
}
