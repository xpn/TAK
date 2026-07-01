package proxy

import (
	"fmt"
	"io"
	"net"
)

func TCPServerToConnProxy(localAddr string, localPort int, callback func() (net.Conn, error)) {
	fmt.Printf("[*] Starting TCP server to connection proxy: [%s:%d]\n", localAddr, localPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localAddr, localPort))
	if err != nil {
		fmt.Printf("Failed to start local server: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("[*] Local server listening on %s:%d\n", localAddr, localPort)

	for {
		localConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[!] Failed to accept local connection: %v\n", err)
			continue
		}

		go func() {
			defer localConn.Close()

			conn, err := callback()
			if err != nil {
				fmt.Printf("[!] Failed to establish remote connection: %v\n", err)
				return
			}

			fmt.Println("[*] Remote connection established")
			defer conn.Close()

			// Forward data between localConn and conn
			go func() {
				_, _ = io.Copy(conn, localConn)
			}()
			_, _ = io.Copy(localConn, conn)
		}()
	}
}
