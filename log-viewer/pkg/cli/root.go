package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"xpnsec.com/log-viewer/v2/pkg/sessions"

	"xpnsec.com/log-viewer/v2/pkg/cli/list"
)

var proxyHost string
var clientCertPath string
var clientKeyPath string
var sessionID string

var rootCmd = &cobra.Command{
	Use:   "log-viewer",
	Short: "A tool for viewing logs from the auth server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := sessions.NewClient(clientCertPath, clientKeyPath, proxyHost)
		if err != nil {
			fmt.Printf("Failed to create client: %v\n", err)
			return
		}

		err = client.DumpRecordedSession(sessionID, true)
		if err != nil {
			fmt.Printf("Failed to dump recorded session: %v\n", err)
		}

	},
}

func Execute() {
	listCmd := list.ListCmd
	rootCmd.AddCommand(listCmd)

	rootCmd.Flags().StringVarP(&proxyHost, "proxy", "x", "", "Proxy host")
	rootCmd.Flags().StringVarP(&clientCertPath, "cert", "c", "", "Client certificate path")
	rootCmd.Flags().StringVarP(&clientKeyPath, "key", "k", "", "Client key path")
	rootCmd.Flags().StringVarP(&sessionID, "session-id", "i", "", "Session ID")

	rootCmd.MarkFlagRequired("proxy-host")
	rootCmd.MarkFlagRequired("client-cert")
	rootCmd.MarkFlagRequired("client-key")
	rootCmd.MarkFlagRequired("session-id")

	cobra.CheckErr(rootCmd.Execute())
}
