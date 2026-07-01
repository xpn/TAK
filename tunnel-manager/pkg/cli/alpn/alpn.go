package cli

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/dialer/alpn"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/proxy"
)

var alpnProtocols string
var clientCertPath string
var clientKeyPath string
var proxyAddr string
var bindAddr string

var AlpnCmd = &cobra.Command{
	Use:   "alpn",
	Short: "Start a TLS tunnel specifying ALPN protocols",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		alpnNames := strings.Split(alpnProtocols, ",")
		certPath := clientCertPath
		keyPath := clientKeyPath
		bindHost, bindPort, err := net.SplitHostPort(bindAddr)
		if err != nil {
			return err
		}

		bindPortNum, err := strconv.Atoi(bindPort)
		if err != nil {
			return err
		}

		proxy.TCPServerToConnProxy(bindHost, bindPortNum, func() (net.Conn, error) {
			alpnDialer := alpn.New(alpn.ALPNConfig{
				NextProtos:      alpnNames,
				CertificatePath: certPath,
				PrivateKeyPath:  keyPath,
			})

			alpnConn, err := alpnDialer.Dial(context.Background(), proxyAddr)
			if err != nil {
				return nil, err
			}
			return alpnConn, nil
		})

		return nil
	},
}

func init() {
	AlpnCmd.Flags().StringVarP(&alpnProtocols, "protocols", "p", "", "Comma-separated list of ALPN protocols")
	AlpnCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "", "Path to client certificate")
	AlpnCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "", "Path to client key")
	AlpnCmd.Flags().StringVarP(&proxyAddr, "proxy", "x", "", "Proxy address (host:port)")
	AlpnCmd.Flags().StringVarP(&bindAddr, "bind", "b", "", "Bind address (host:port)")

	AlpnCmd.MarkFlagRequired("protocols")
	AlpnCmd.MarkFlagRequired("client-cert")
	AlpnCmd.MarkFlagRequired("client-key")
	AlpnCmd.MarkFlagRequired("proxy")
	AlpnCmd.MarkFlagRequired("bind")
}
