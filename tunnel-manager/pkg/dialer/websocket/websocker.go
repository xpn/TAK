package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketDialer struct {
	target string
}

type WebSocketConnection struct {
	Conn       *websocket.Conn
	ReadBuffer []byte
}

func (c *WebSocketConnection) Write(data []byte) (n int, err error) {
	fmt.Println("WebSocketConnection Write called")
	// Dump data being written for debugging

	err = c.Conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return len(data), nil
}

func (c *WebSocketConnection) Read(msg []byte) (n int, err error) {
	fmt.Println("WebSocketConnection Read called")
	// If there is data in the read buffer, serve from there first
	if len(c.ReadBuffer) > 0 {
		n = copy(msg, c.ReadBuffer)
		c.ReadBuffer = c.ReadBuffer[n:]
		return n, nil
	}

	_, message, err := c.Conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Copy data to msg and store any excess in ReadBuffer
	n = copy(msg, message)
	if n < len(message) {
		c.ReadBuffer = message[n:]
	}
	return n, nil
}
func (c *WebSocketConnection) SetDeadline(t time.Time) error {
	err1 := c.Conn.SetReadDeadline(t)
	err2 := c.Conn.SetWriteDeadline(t)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (c *WebSocketConnection) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *WebSocketConnection) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func (c *WebSocketConnection) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *WebSocketConnection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *WebSocketConnection) Close() error {
	return c.Conn.Close()
}

func (d *WebSocketDialer) Close() error {
	return nil
}

func New() *WebSocketDialer {
	return &WebSocketDialer{}
}

func (d *WebSocketDialer) Dial(ctx context.Context, target string) (net.Conn, error) {
	// Ignore SSL certificate errors
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	headers := make(map[string][]string)
	// Set any required headers here
	headers["Sec-WebSocket-Protocol"] = []string{"alpn"}
	c, _, err := dialer.Dial(target, headers)
	if err != nil {
		return nil, err
	}

	return &WebSocketConnection{Conn: c}, nil
}
