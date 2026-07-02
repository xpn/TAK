package alpn

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

type ALPNDialer struct {
	config ALPNConfig
	target string
}

type ALPNConnection struct {
	Conn       *tls.Conn
	ReadBuffer []byte
}

type ALPNConfig struct {
	NextProtos      []string
	CertificatePath string
	PrivateKeyPath  string
}

func (c *ALPNConnection) Write(data []byte) (n int, err error) {
	return c.Conn.Write(data)
}

func (c *ALPNConnection) Read(msg []byte) (n int, err error) {
	return c.Conn.Read(msg)
}

func (c *ALPNConnection) SetDeadline(t time.Time) error {
	err1 := c.Conn.SetReadDeadline(t)
	err2 := c.Conn.SetWriteDeadline(t)
	if err1 != nil {
		return err1
	}
	return err2
}

func (c *ALPNConnection) Close() error {
	return c.Conn.Close()
}

func (c *ALPNConnection) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *ALPNConnection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *ALPNConnection) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *ALPNConnection) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func New(config ALPNConfig) *ALPNDialer {
	return &ALPNDialer{
		config: config,
	}
}

func (d *ALPNDialer) Dial(ctx context.Context, target string) (net.Conn, error) {

	// If client certificate and key paths are provided, load them
	clientCertificate := tls.Certificate{}
	var err error
	if d.config.CertificatePath != "" && d.config.PrivateKeyPath != "" {
		clientCertificate, err = tls.LoadX509KeyPair(d.config.CertificatePath, d.config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
	}

	dialer := tls.Dialer{
		Config: &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         d.config.NextProtos,
			Certificates:       []tls.Certificate{clientCertificate},
		},
	}
	conn, err := dialer.DialContext(ctx, "tcp", target)
	if err != nil {
		return nil, err
	}
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return nil, err
	}

	return &ALPNConnection{Conn: tlsConn}, nil
}
