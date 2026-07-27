package cert

import (
	"crypto/x509"
	"encoding/pem"

	"xpnsec.com/certificate-tool/v2/pkg/crypto"
)

func pemConvert(derBytes []byte, blockType string) []byte {
	pemBlock := &pem.Block{
		Type:    blockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	out := pem.EncodeToMemory(pemBlock)
	return out
}

func GenerateAppPublicPrivateKey() ([]byte, []byte, error) {
	key, err := crypto.GenerateRSAPrivateKey(2048)
	if err != nil {
		return nil, nil, err
	}

	keyDER := x509.MarshalPKCS1PrivateKey(key)
	pubDER, err := x509.MarshalPKIXPublicKey(key.Public())

	// Convert to PEM
	pubPEM := pemConvert(pubDER, "PUBLIC KEY")

	return pubPEM, keyDER, err
}
