package cli

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"xpnsec.com/shared/v2/pkg/connection"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/dialer/alpn"
	"xpnsec.com/teleport-tunnel-manager/v2/pkg/proxy"
)

var alpnProtocols string
var bindAddr string
var connectionOptions connection.Options

var AlpnCmd = &cobra.Command{
	Use:   "alpn",
	Short: "Start a TLS tunnel specifying ALPN protocols",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		alpnNames := strings.Split(alpnProtocols, ",")
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
				CertificatePath: connectionOptions.ClientCert,
				PrivateKeyPath:  connectionOptions.ClientKey,
			})

			alpnConn, err := alpnDialer.Dial(context.Background(), connectionOptions.Proxy)
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
	AlpnCmd.Flags().StringVarP(&bindAddr, "bind", "b", "", "Bind address (host:port)")
	connectionOptions.AddProxyFlag(AlpnCmd.Flags())
	connectionOptions.AddClientCredentialFlags(AlpnCmd.Flags())
	AlpnCmd.MarkFlagsRequiredTogether("client-cert", "client-key")

	AlpnCmd.MarkFlagRequired("protocols")
	AlpnCmd.MarkFlagRequired("proxy")
	AlpnCmd.MarkFlagRequired("bind")
}
