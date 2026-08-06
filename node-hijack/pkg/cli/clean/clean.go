package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	api "xpnsec.com/node-hijack/v2/pkg/api"
	"xpnsec.com/shared/v2/pkg/connection"
)

var name string
var connectionOptions connection.Options

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up all hijacked entries",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
		if err != nil {
			panic(err)
		}

		// nodes, err := client.GetAllNodes(context.Background())
		// if err != nil {
		// 	panic(err)
		// }

		fmt.Println("[*] Cleaning up node")
		client.DeleteNode(context.Background(), name)
		fmt.Println("[*] Node cleaned up")
	},
}

func init() {
	CleanCmd.Flags().StringVarP(&name, "name", "n", "", "Node name to clean")
	connectionOptions.AddProxyFlag(CleanCmd.Flags())
	connectionOptions.AddClientCredentialFlags(CleanCmd.Flags())
	connectionOptions.AddClusterNameFlag(CleanCmd.Flags())
	CleanCmd.MarkFlagRequired("name")
	CleanCmd.MarkFlagRequired("proxy")
	CleanCmd.MarkFlagRequired("client-cert")
	CleanCmd.MarkFlagRequired("client-key")
	CleanCmd.MarkFlagRequired("cluster-name")
}
