package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"xpnsec.com/log-viewer/v2/pkg/sessions"
	"xpnsec.com/shared/v2/pkg/connection"

	"xpnsec.com/log-viewer/v2/pkg/cli/list"
)

var sessionID string
var connectionOptions connection.Options

var rootCmd = &cobra.Command{
	Use:   "log-viewer",
	Short: "A tool for viewing logs from the auth server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := sessions.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
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

	rootCmd.Flags().StringVarP(&sessionID, "session-id", "i", "", "Session ID")
	connectionOptions.AddProxyFlag(rootCmd.Flags())
	connectionOptions.AddClientCredentialFlags(rootCmd.Flags())
	connectionOptions.AddClusterNameFlag(rootCmd.Flags())

	rootCmd.MarkFlagRequired("proxy")
	rootCmd.MarkFlagRequired("client-cert")
	rootCmd.MarkFlagRequired("client-key")
	rootCmd.MarkFlagRequired("cluster-name")
	rootCmd.MarkFlagRequired("session-id")

	cobra.CheckErr(rootCmd.Execute())
}
