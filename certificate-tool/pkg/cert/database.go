package cert

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"

	"xpnsec.com/certificate-tool/v2/pkg/crypto"
)

func GenerateDatabaseCSR(username string) ([]byte, []byte, error) {
	key, err := crypto.GenerateRSAPrivateKey(2048)
	if err != nil {
		return nil, nil, err
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: username},
	}

	keyDER := x509.MarshalPKCS1PrivateKey(key)

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, key)
	if err != nil {
		return nil, nil, err
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	return csrPEM, keyDER, nil
}
