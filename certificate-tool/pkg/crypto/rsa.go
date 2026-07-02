package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateRSAPrivateKey(length int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, err
	}
	return key, nil
}
