package cli

import (
	"github.com/spf13/cobra"
	ssh "xpnsec.com/certificate-tool/v2/pkg/cli/ssh"
	windows "xpnsec.com/certificate-tool/v2/pkg/cli/windows"
)

var proxyHost string
var clientCertPath string
var clientKeyPath string
var username string
var connectHost string

var rootCmd = &cobra.Command{
	Use:   "certificate-tool",
	Short: "A certificate multi-tool for generating certificates required by Teleport",
	Long:  ``,
}

func Execute() {

	rootCmd.AddCommand(windows.WindowsCmd)
	rootCmd.AddCommand(ssh.SSHCmd)

	cobra.CheckErr(rootCmd.Execute())
}
