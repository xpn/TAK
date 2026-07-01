package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	alpn "xpnsec.com/teleport-tunnel-manager/v2/pkg/cli/alpn"
	websocket "xpnsec.com/teleport-tunnel-manager/v2/pkg/cli/websocket"
)

var rootCmd = &cobra.Command{
	Use:   "tunnel-manager",
	Short: "Tunnel Manager is a research CLI tool for managing Teleport tunnels",
	Long:  `Tunnel Manager is a research CLI tool for managing Teleport tunnels`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {

	rootCmd.AddCommand(alpn.AlpnCmd)
	rootCmd.AddCommand(websocket.WebSocketCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
