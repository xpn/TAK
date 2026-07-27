package cli

import (
	"github.com/spf13/cobra"
	app "xpnsec.com/certificate-tool/v2/pkg/cli/direct/app"
	database "xpnsec.com/certificate-tool/v2/pkg/cli/direct/database"
	windows "xpnsec.com/certificate-tool/v2/pkg/cli/direct/windows"
	ssh "xpnsec.com/certificate-tool/v2/pkg/cli/ssh"
)

var rootCmd = &cobra.Command{
	Use:   "certificate-tool",
	Short: "A certificate multi-tool for generating certificates required by Teleport",
	Long:  ``,
}

func Execute() {

	// expose two subcommands, "direct" and "teleport"
	directCmd := &cobra.Command{
		Use:   "direct",
		Short: "Generate certificates which must be presented to the server directly",
		Long:  ``,
	}
	teleportCmd := &cobra.Command{
		Use:   "teleport",
		Short: "Generate certificates which are used to authenticate with the Teleport server",
		Long:  ``,
	}

	directCmd.AddCommand(windows.WindowsCmd)
	directCmd.AddCommand(database.DatabaseCmd)
	directCmd.AddCommand(app.AppCmd)

	teleportCmd.AddCommand(ssh.SSHCmd)

	rootCmd.AddCommand(directCmd)
	rootCmd.AddCommand(teleportCmd)

	cobra.CheckErr(rootCmd.Execute())
}
