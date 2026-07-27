package authserver

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"

	authserver "xpnsec.com/certificate-tool/v2/pkg/grpc"
)

type AuthServerClient struct {
	client      *authserver.AuthServiceClient
	clusterName string
}

type HostCerts struct {
	SSH []byte
	TLS []byte
}

type WindowsHostCerts struct {
	Cert []byte
}

type DatabaseCerts struct {
	Cert []byte
}

type AppCerts struct {
	Cert []byte
}

func NewClient(certPath, keyPath, target, clusterName string) (*AuthServerClient, error) {
	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Hex encode the cluster name to use as the SNI
	clusterNameHex := hex.EncodeToString([]byte(clusterName))

	creds := credentials.NewTLS(&tls.Config{
		NextProtos:         []string{"teleport-auth@" + clusterNameHex + ".teleport.cluster.local", "h2"},
		Certificates:       []tls.Certificate{clientCertificate},
		InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
	)

	client := authserver.NewAuthServiceClient(conn)

	return &AuthServerClient{
		client:      &client,
		clusterName: clusterName,
	}, nil

}

func (c *AuthServerClient) GenerateSSHHostCertificate(ctx context.Context, nodeName string, nodeId string, publicKey []byte, publicTLSKey []byte) (*HostCerts, error) {
	certs, err := (*c.client).GenerateHostCerts(ctx, &authserver.HostCertsRequest{
		NodeName:     nodeName,
		HostID:       nodeId,
		Role:         "Node",
		PublicSSHKey: publicKey,
		PublicTLSKey: publicTLSKey,
		AdditionalPrincipals: []string{
			nodeName,
			nodeId + "." + c.clusterName,
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

func (c *AuthServerClient) GenerateWindowsHostCertificate(ctx context.Context, csr []byte) (*WindowsHostCerts, error) {

	certs, err := (*c.client).GenerateWindowsDesktopCert(ctx, &authserver.WindowsDesktopCertRequest{
		CSR: csr,
		TTL: int64(1 * time.Hour),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate windows host certificate: %w", err)
	}

	return &WindowsHostCerts{
		Cert: certs.Cert,
	}, nil
}

func (c *AuthServerClient) GenerateDatabaseCertificate(ctx context.Context, csr []byte) (*DatabaseCerts, error) {
	certs, err := (*c.client).GenerateDatabaseCert(ctx, &authserver.DatabaseCertRequest{
		CSR: csr,
		TTL: int64(1 * time.Hour),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate database certificate: %w", err)
	}

	return &DatabaseCerts{
		Cert: certs.Cert,
	}, nil
}

func (c *AuthServerClient) GenerateUserAppCertificate(ctx context.Context, username string, appName string, publicAddress string, pubTlsKey []byte) (*AppCerts, error) {

	expiryTime := timestamppb.New(time.Now().Add(time.Hour * 24))

	certs, err := (*c.client).GenerateUserCerts(ctx, &authserver.UserCertsRequest{
		Username:     username,
		TLSPublicKey: pubTlsKey,
		Expires:      expiryTime,
		RouteToApp: &authserver.RouteToApp{
			Name:        appName,
			PublicAddr:  publicAddress,
			ClusterName: c.clusterName,
		},
		Usage: authserver.UserCertsRequest_App,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate app certificate: %w", err)
	}

	return &AppCerts{
		Cert: certs.TLS,
	}, nil
}
