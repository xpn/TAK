package sshserver

import (
	"fmt"
	"log"
	"os"

	gssh "github.com/gliderlabs/ssh"
	ssh "golang.org/x/crypto/ssh"
)

type SSHServer struct {
	certSigner ssh.Signer
}

func NewSSHServer(certPath, hostKeyPath string) *SSHServer {
	signer := &SSHServer{}

	certSigner, err := signer.hostCertSigner(hostKeyPath, certPath)
	if err != nil {
		log.Fatalf("host certificate setup failed: %v", err)
	}

	return &SSHServer{certSigner: certSigner}
}

func (*SSHServer) hostCertSigner(keyPath, certPath string) (ssh.Signer, error) {
	keyPEM, _ := os.ReadFile(keyPath)
	signer, err := ssh.ParsePrivateKey(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse host key: %w", err)
	}

	certBytes, _ := os.ReadFile(certPath)
	pub, _, _, _, err := ssh.ParseAuthorizedKey(certBytes)
	if err != nil {
		return nil, fmt.Errorf("parse host cert: %w", err)
	}
	cert, ok := pub.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("%s is not an SSH certificate", certPath)
	}
	if cert.CertType != ssh.HostCert {
		return nil, fmt.Errorf("cert is type %d, expected a host cert (sign with -h)", cert.CertType)
	}

	return ssh.NewCertSigner(cert, signer) // validates cert pubkey == signer pubkey
}

func (s *SSHServer) Start(callback func(s gssh.Session)) error {

	srv := &gssh.Server{
		Addr:    ":2223",
		Handler: callback,
		PublicKeyHandler: func(ctx gssh.Context, key gssh.PublicKey) bool {
			return true // accept-any for demo; validate real client keys here
		},
	}

	srv.AddHostKey(s.certSigner) // forwarded to ServerConfig.AddHostKey
	log.Fatal(srv.ListenAndServe())

	return nil
}
