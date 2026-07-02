package cert

import (
	"crypto/x509"

	"xpnsec.com/certificate-tool/v2/pkg/crypto"
)

func GenerateAppPublicPrivateKey() ([]byte, []byte, error) {
	key, err := crypto.GenerateRSAPrivateKey(2048)
	if err != nil {
		return nil, nil, err
	}

	keyDER := x509.MarshalPKCS1PrivateKey(key)
	pubDER, err := x509.MarshalPKIXPublicKey(key.Public())

	return pubDER, keyDER, err
}
