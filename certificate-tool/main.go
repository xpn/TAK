package main

import (
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/gravitational/teleport/api/utils/keys"
	"golang.org/x/crypto/ssh"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
)

func generateWindowsCert() {
	cert, key, err := cert.GenerateWindowsCSR("Administrator", "TELEPORT-WINCLI")

	if err != nil {
		panic(err)
	}

	// Write the CSR and key to files
	os.WriteFile("/tmp/administrator.csr", cert, 0644)
	os.WriteFile("/tmp/administrator.key", key, 0600)
}

func generateSSHKeys() {
	privateKey, err := cert.GenerateSSHKey()
	if err != nil {
		panic(err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})

	sshPub, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		panic(err)
	}

	tlsPub, err := keys.MarshalPublicKey(privateKey.Public())
	if err != nil {
		panic(err)
	}

	publicSSH := ssh.MarshalAuthorizedKey(sshPub)

	// Write the certificate and key to files
	os.WriteFile("/tmp/n.pub", publicSSH, 0644)
	os.WriteFile("/tmp/n.key", privPEM, 0600)
	os.WriteFile("/tmp/n_tls.pub", tlsPub, 0644)
}

func generateDatabaseCert() {
	csr, key, err := cert.GenerateDatabaseCSR("alice")
	if err != nil {
		panic(err)
	}

	// Write the CSR and key to files
	os.WriteFile("/tmp/dbuser.csr", csr, 0644)
	os.WriteFile("/tmp/dbuser.key", key, 0600)
}

func generateAppCert() {
}

func generateAppCertFromAuthServer() {
	pub, priv, err := cert.GenerateAppPublicPrivateKey()
	if err != nil {
		panic(err)
	}

	// Convert the public key to PEM format
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub,
	})

	// Write the public and private keys to files
	os.WriteFile("/tmp/app.pub", pubPEM, 0644)
	os.WriteFile("/tmp/app.key", priv, 0600)
}

func main() {

	switch os.Args[1] {
	case "ssh":
		generateSSHKeys()
	case "windows":
		generateWindowsCert()
	case "database":
		generateDatabaseCert()
	case "app":
		generateAppCert()
	case "user-app":
		generateAppCertFromAuthServer()
	default:
		panic("unknown command: " + os.Args[1])
	}
}
