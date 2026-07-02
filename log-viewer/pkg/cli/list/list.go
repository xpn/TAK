package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"xpnsec.com/log-viewer/v2/pkg/sessions"
)

var clientCertPath, clientKeyPath, proxyHost string
var sessionInfo []sessions.SessionInfo

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded sessions",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := sessions.NewClient(clientCertPath, clientKeyPath, proxyHost)
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
	ListCmd.Flags().StringVar(&clientCertPath, "cert", "", "Path to client certificate")
	ListCmd.Flags().StringVar(&clientKeyPath, "key", "", "Path to client key")
	ListCmd.Flags().StringVar(&proxyHost, "proxy", "", "Proxy host")

	cobra.MarkFlagRequired(ListCmd.Flags(), "cert")
	cobra.MarkFlagRequired(ListCmd.Flags(), "key")
	cobra.MarkFlagRequired(ListCmd.Flags(), "proxy")

}
