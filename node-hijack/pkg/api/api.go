package authserver

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/gravitational/teleport/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authserver "xpnsec.com/node-hijack/v2/pkg/grpc"
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

func (c *AuthServerClient) GetNode(ctx context.Context, nodeName string) (*types.ServerV2, error) {
	result, err := (*c.client).GetNode(ctx, &types.ResourceInNamespaceRequest{
		Name:      nodeName,
		Namespace: "default",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	return result, nil
}

func (c *AuthServerClient) GetAllNodes(ctx context.Context) error {
	result, err := (*c.client).ListUnifiedResources(ctx, &authserver.ListUnifiedResourcesRequest{
		Kinds: []string{"node"},
	})
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}
	fmt.Println(result)
	return nil
}

func (c *AuthServerClient) UpdateNode(ctx context.Context, updated *types.ServerV2) error {
	_, err := (*c.client).UpsertNode(ctx, updated)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}
	return nil
}

func (c *AuthServerClient) UpsertNode(ctx context.Context, nodeName string) error {
	result, err := (*c.client).UpsertNode(ctx, &types.ServerV2{
		Kind: "Node",
		Spec: types.ServerSpecV2{
			//	Addr:     "1.1.1.1",
			Hostname:  nodeName,
			UseTunnel: true,
			//	PeerAddr: "1.2.3.4",
		},
		Metadata: types.Metadata{
			Name: nodeName,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upsert node: %w", err)
	}

	fmt.Println(result.Namespace)
	return nil
}

func (c *AuthServerClient) DeleteNode(ctx context.Context, nodeName string) error {
	_, err := (*c.client).DeleteNode(ctx, &types.ResourceInNamespaceRequest{
		Name:      nodeName,
		Namespace: "default",
	})
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}
