package cli

import (
	"github.com/spf13/cobra"
	clean "xpnsec.com/node-hijack/v2/pkg/cli/clean"
	hijack "xpnsec.com/node-hijack/v2/pkg/cli/hijack"
	list "xpnsec.com/node-hijack/v2/pkg/cli/list"
	mitm "xpnsec.com/node-hijack/v2/pkg/cli/mitm"
)

var rootCmd = &cobra.Command{
	Use:   "node-hijack",
	Short: "A tool designed to make Teleport Node hijacking easier",
	Long:  ``,
}

func Execute() {
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(hijack.HijackCmd)
	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(mitm.MITMCmd)

	cobra.CheckErr(rootCmd.Execute())
}
