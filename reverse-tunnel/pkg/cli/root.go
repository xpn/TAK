package cli

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	sshreversetunnel "xpnsec.com/reverse-tunnel/v2/pkg/ssh-reverse-tunnel"
	"xpnsec.com/shared/v2/pkg/connection"
)

var clientSSHCertPath string
var username string
var connectHost string
var connectionOptions connection.Options

var rootCmd = &cobra.Command{
	Use:   "reverse-tunnel",
	Short: "A Teleport SSH tunnel tool for establishing reverse tunnels",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// conn, err := net.Dial("tcp", proxyHost)
		// if err != nil {
		// 	fmt.Printf("[!] Failed to create TCP connection: %v\n", err)
		// 	return
		// }
		// defer conn.Close()
		clientCertificate, err := tls.LoadX509KeyPair(connectionOptions.ClientCert, connectionOptions.ClientKey)
		if err != nil {
			fmt.Printf("[!] Failed to load client certificate: %v\n", err)
			return
		}

		config := &tls.Config{
			NextProtos:         []string{"teleport-reversetunnel"},
			Certificates:       []tls.Certificate{clientCertificate},
			InsecureSkipVerify: true,
		}

		conn, err := tls.Dial("tcp", connectionOptions.Proxy, config)
		if err != nil {
			fmt.Printf("[!] Failed to create TLS connection: %v\n", err)
			return
		}
		defer conn.Close()

		tunnel := sshreversetunnel.New(username, "@remote-auth-server", connectHost, clientSSHCertPath, connectionOptions.ClientKey)
		err = tunnel.Connect(conn)
		if err != nil {
			fmt.Printf("[!] Failed to establish SSH reverse tunnel: %v\n", err)
			return
		}
		select {}
	},
}

func Execute() {

	rootCmd.Flags().StringVarP(&clientSSHCertPath, "ssh-cert", "s", "", "SSH certificate path")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	rootCmd.Flags().StringVarP(&connectHost, "connect-host", "o", "", "Host to tunnel connections to")
	connectionOptions.AddProxyFlag(rootCmd.Flags())
	connectionOptions.AddClientCredentialFlags(rootCmd.Flags())

	rootCmd.MarkFlagRequired("proxy")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("client-cert")
	rootCmd.MarkFlagRequired("ssh-cert")
	rootCmd.MarkFlagRequired("client-key")
	rootCmd.MarkFlagRequired("connect-host")

	cobra.CheckErr(rootCmd.Execute())

}
