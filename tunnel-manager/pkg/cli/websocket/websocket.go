package cli

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"xpnsec.com/shared/v2/pkg/connection"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/dialer/websocket"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/proxy"
)

var bindAddr string
var connectionOptions connection.Options

var WebSocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "Start a WebSocket tunnel",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		bindHost, bindPort, err := net.SplitHostPort(bindAddr)
		if err != nil {
			return err
		}

		bindPortNum, err := strconv.Atoi(bindPort)
		if err != nil {
			return err
		}

		proxy.TCPServerToConnProxy(bindHost, bindPortNum, func() (net.Conn, error) {
			wsDialer := websocket.New()

			targetUrl := fmt.Sprintf("wss://%s/webapi/connectionupgrade", connectionOptions.Proxy)

			wsConn, err := wsDialer.Dial(context.Background(), targetUrl)
			return wsConn, err
		})

		return nil
	},
}

func init() {
	WebSocketCmd.Flags().StringVarP(&bindAddr, "bind", "b", "", "Bind address (host:port)")
	connectionOptions.AddProxyFlag(WebSocketCmd.Flags())

	WebSocketCmd.MarkFlagRequired("proxy")
	WebSocketCmd.MarkFlagRequired("bind")
}
