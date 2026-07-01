package sshreversetunnel

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func loadSSHCertFromBytes(certBytes []byte) (*ssh.Certificate, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH certificate: %v", err)
	}

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("parsed key is not an SSH certificate")
	}

	return cert, nil
}

func parsePrivateKey(keyBytes []byte) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return signer, nil
}
