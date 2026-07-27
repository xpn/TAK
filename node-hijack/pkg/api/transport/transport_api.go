package transport_api

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	transport "xpnsec.com/node-hijack/v2/pkg/grpc/transport"
)

type TransportServerClient struct {
	client *transport.TransportServiceClient
}

func NewClient(target string) (*TransportServerClient, error) {
	clientCertificate, err := tls.LoadX509KeyPair("/tmp/.tsh/keys/xpn-teleport-server/regular-user.crt", "/tmp/.tsh/keys/xpn-teleport-server/regular-user.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		NextProtos:         []string{"teleport-proxy-ssh-grpc"},
		Certificates:       []tls.Certificate{clientCertificate},
		InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		return nil, err
	}

	client := transport.NewTransportServiceClient(conn)

	return &TransportServerClient{
		client: &client,
	}, nil

}

func (c *TransportServerClient) CreateSSHConnection(ctx context.Context, cluster string, hostport string) (*sshConn, error) {
	stream, err := (*c.client).ProxySSH(ctx)
	if err != nil {
		return nil, err
	}

	if err = stream.Send(&transport.ProxySSHRequest{DialTarget: &transport.TargetHost{
		HostPort: hostport,
		Cluster:  cluster,
	}}); err != nil {
		return nil, err
	}

	data, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	fmt.Println(data)

	ctx, cancel := context.WithCancel(context.Background())
	conn := newSSHConn(cancel, stream)

	return conn, nil
}

type sshConn struct {
	stream grpc.BidiStreamingClient[transport.ProxySSHRequest, transport.ProxySSHResponse]
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once

	readMu   sync.Mutex
	leftover []byte // bytes from a Recv not yet consumed by Read

	writeMu sync.Mutex
}

func newSSHConn(cancel context.CancelFunc, grpcStream grpc.BidiStreamingClient[transport.ProxySSHRequest, transport.ProxySSHResponse]) *sshConn {
	return &sshConn{
		stream: grpcStream,
		cancel: cancel,
	}
}

func (c *sshConn) Read(bytes []byte) (int, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()

	for len(c.leftover) == 0 {
		msg, err := c.stream.Recv()
		if err != nil {
			return 0, err // io.EOF or a gRPC status; SSH treats non-nil as conn dead
		}

		c.leftover = msg.GetSsh().GetPayload()
	}

	n := copy(bytes, c.leftover)
	c.leftover = c.leftover[n:]
	return n, nil
}

func (c *sshConn) Write(bytes []byte) (int, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	// Copy: x/crypto/ssh reuses its packet write buffer, and you must not
	// hand a slice that may be mutated to Send.
	data := append([]byte(nil), bytes...)
	err := c.stream.Send(&transport.ProxySSHRequest{Frame: &transport.ProxySSHRequest_Ssh{Ssh: &transport.Frame{Payload: data}}})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (c *sshConn) Close() error {
	c.once.Do(func() { c.cancel() })
	return nil
}

// Addr + deadline stubs — SSH's handshake path doesn't require real deadlines.
type grpcAddr struct{}

func (grpcAddr) Network() string { return "grpc" }
func (grpcAddr) String() string  { return "grpc" }

func (c *sshConn) LocalAddr() net.Addr                { return grpcAddr{} }
func (c *sshConn) RemoteAddr() net.Addr               { return grpcAddr{} }
func (c *sshConn) SetDeadline(t time.Time) error      { return nil }
func (c *sshConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *sshConn) SetWriteDeadline(t time.Time) error { return nil }
