package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"xpnsec.com/log-viewer/v2/pkg/sessions"
	"xpnsec.com/shared/v2/pkg/connection"
)

var sessionInfo []sessions.SessionInfo
var connectionOptions connection.Options

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded sessions",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := sessions.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
		if err != nil {
			fmt.Printf("Failed to create client: %v\n", err)
			return
		}

		sessionInfo, err := client.ListRecordedSessions()
		if err != nil {
			fmt.Printf("Failed to list recorded sessions: %v\n", err)
			return
		}

		for _, info := range sessionInfo {
			fmt.Printf("[%s]: %s@%s - %s - %d seconds\n", info.SessionId, info.Username, info.Hostname, info.Time, int64(info.Duration.Seconds()))
		}

	},
}

func init() {
	connectionOptions.AddProxyFlag(ListCmd.Flags())
	connectionOptions.AddClientCredentialFlags(ListCmd.Flags())
	connectionOptions.AddClusterNameFlag(ListCmd.Flags())

	cobra.MarkFlagRequired(ListCmd.Flags(), "client-cert")
	cobra.MarkFlagRequired(ListCmd.Flags(), "client-key")
	cobra.MarkFlagRequired(ListCmd.Flags(), "proxy")
	cobra.MarkFlagRequired(ListCmd.Flags(), "cluster-name")

}
