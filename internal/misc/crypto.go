package misc

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func Base64ToPrivateKey(str string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func PublicKeyToBytes(key *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	bytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return bytes, nil
}
