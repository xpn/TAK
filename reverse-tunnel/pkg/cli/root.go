package cli

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	sshreversetunnel "xpnsec.com/reverse-tunnel/v2/pkg/ssh-reverse-tunnel"
)

var proxyHost string
var clientCertPath string
var clientKeyPath string
var username string
var connectHost string

var rootCmd = &cobra.Command{
	Use:   "reverse-tunnel",
	Short: "A Teleport SSH tunnel tool for establishing reverse tunnels",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := net.Dial("tcp", proxyHost)
		if err != nil {
			fmt.Printf("[!] Failed to create TCP connection: %v\n", err)
			return
		}
		defer conn.Close()

		tunnel := sshreversetunnel.New(username, "@remote-auth-server", connectHost, clientCertPath, clientKeyPath)
		err = tunnel.Connect(conn)
		if err != nil {
			fmt.Printf("[!] Failed to establish SSH reverse tunnel: %v\n", err)
			return
		}
		select {}
	},
}

func Execute() {

	rootCmd.Flags().StringVarP(&proxyHost, "proxy-host", "x", "", "Proxy host")
	rootCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "", "Client certificate path")
	rootCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "", "Client key path")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	rootCmd.Flags().StringVarP(&connectHost, "connect-host", "o", "", "Host to tunnel connections to")

	rootCmd.MarkFlagRequired("proxy-host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("client-cert")
	rootCmd.MarkFlagRequired("client-key")
	rootCmd.MarkFlagRequired("connect-host")

	cobra.CheckErr(rootCmd.Execute())

}
