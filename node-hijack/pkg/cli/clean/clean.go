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

var name string
var clientCert string
var clientKey string

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up all hijacked entries",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(clientCert, clientKey, "10.1.10.1:8443")
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
	CleanCmd.Flags().StringVarP(&clientCert, "client-cert", "c", "", "Existing Node Cert")
	CleanCmd.Flags().StringVarP(&clientKey, "client-key", "k", "", "Existing Node Key")

}
