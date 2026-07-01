package main

import (
	"fmt"

	"xpnsec.com/teleport-tunnel-manager/v2/pkg/cli"
)

// func unixSocketToConnProxy(socketPath string, conn net.Conn) {
// 	listener, err := net.Listen("unix", socketPath)
// 	if err != nil {
// 		fmt.Printf("[!] Failed to start Unix socket server: %v\n", err)
// 		return
// 	}
// 	defer listener.Close()
// 	fmt.Printf("[*] Unix socket server listening on %s\n", socketPath)
// 	for {
// 		localConn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Printf("[!] Failed to accept Unix socket connection: %v\n", err)
// 			continue
// 		}

// 		go func() {
// 			defer localConn.Close()
// 			// Forward data between localConn and conn
// 			go func() {
// 				_, _ = io.Copy(conn, localConn)
// 			}()
// 			_, _ = io.Copy(localConn, conn)
// 		}()
// 	}
// }

func main() {

	fmt.Printf(`
╔╦╗┬ ┬┌┐┌┌┐┌┌─┐┬   ╔╦╗┌─┐┌┐┌┌─┐┌─┐┌─┐┬─┐
 ║ │ │││││││├┤ │───║║║├─┤│││├─┤│ ┬├┤ ├┬┘
 ╩ └─┘┘└┘┘└┘└─┘┴─┘ ╩ ╩┴ ┴┘└┘┴ ┴└─┘└─┘┴└─
                                  @_xpn_
`)

	cli.Execute()

	// } else {
	// 	// Unix socket mode
	// 	socketPath = localTarget
	// 	unixSocketToConnProxy(socketPath, conn)
	// }
}
