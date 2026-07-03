package authserver

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authserver "xpnsec.com/certificate-tool/v2/pkg/grpc"
)

type AuthServerClient struct {
	client *authserver.AuthServiceClient
}

func NewClient(certPath, keyPath, target string) (*AuthServerClient, error) {
	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		NextProtos:         []string{"teleport-auth@6578616d706c652e636f6d.teleport.cluster.local", "h2"},
		Certificates:       []tls.Certificate{clientCertificate},
		InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
	)

	client := authserver.NewAuthServiceClient(conn)

	return &AuthServerClient{
		client: &client,
	}, nil

}

type HostCerts struct {
	SSH []byte
	TLS []byte
}

func (c *AuthServerClient) GenerateSSHHostCertificate(ctx context.Context, nodeName string, publicKey []byte, publicTLSKey []byte) (*HostCerts, error) {
	certs, err := (*c.client).GenerateHostCerts(ctx, &authserver.HostCertsRequest{
		NodeName:     nodeName,
		HostID:       "b14086e9-0294-408d-9b76-0405f2409929",
		Role:         "Node",
		PublicSSHKey: publicKey,
		PublicTLSKey: publicTLSKey,
		AdditionalPrincipals: []string{
			nodeName,
			"b14086e9-0294-408d-9b76-0405f2409929.example.com",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate host certificate: %w", err)
	}

	sshCert := certs.GetSSH()
	tlsCert := certs.GetTLS()

	return &HostCerts{
		SSH: sshCert,
		TLS: tlsCert,
	}, nil
}
