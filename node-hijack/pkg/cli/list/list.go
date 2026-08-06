package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	api "xpnsec.com/node-hijack/v2/pkg/api"
	"xpnsec.com/shared/v2/pkg/connection"
)

var connectionOptions connection.Options

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
		if err != nil {
			panic(err)
		}

		nodes, err := client.GetAllNodes(context.Background())
		if err != nil {
			panic(err)
		}

		for _, node := range nodes {
			fmt.Printf("[%s] %s (%s)\n", node.Name, node.Hostname, node.IP)
		}
	},
}

func init() {
	connectionOptions.AddProxyFlag(ListCmd.Flags())
	connectionOptions.AddClientCredentialFlags(ListCmd.Flags())
	connectionOptions.AddClusterNameFlag(ListCmd.Flags())
	ListCmd.MarkFlagRequired("proxy")
	ListCmd.MarkFlagRequired("client-cert")
	ListCmd.MarkFlagRequired("client-key")
	ListCmd.MarkFlagRequired("cluster-name")
}
