package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	api "xpnsec.com/node-hijack/v2/pkg/api"
)

var outputDir string
var username string
var hostname string
var clientCertPath string
var clientKeyPath string
var proxy string

var clientCert string
var clientKey string

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(clientCert, clientKey, "10.1.10.1:8443")
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
	ListCmd.Flags().StringVarP(&clientCert, "client-cert", "c", "", "Existing Node Cert")
	ListCmd.Flags().StringVarP(&clientKey, "client-key", "k", "", "Existing Node Key")
}
