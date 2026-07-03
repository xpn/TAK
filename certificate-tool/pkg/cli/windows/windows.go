package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
)

var outputDir string
var username string
var hostname string

var WindowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "Windows certificate utilities",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if err := os.MkdirAll(outputDir, 0700); err != nil {
			fmt.Printf("[!] Error creating output dir: %x\n", err)
			return
		}

		cert, key, err := cert.GenerateWindowsCSR(username, hostname)
		if err != nil {
			fmt.Printf("[!] Error generating certificate: %v\n", err)
			return
		}

		csrPath := path.Join(outputDir, fmt.Sprintf("%s-%s.csr", username, hostname))
		keyPath := path.Join(outputDir, fmt.Sprintf("%s-%s.key", username, hostname))

		// Write the CSR and key to files
		os.WriteFile(csrPath, cert, 0744)
		os.WriteFile(keyPath, key, 0600)

		fmt.Printf("[*] Certificates generated:\n\t%s\n\t%s\n", csrPath, keyPath)
	},
}

func init() {
	WindowsCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to save generated certificates")
	WindowsCmd.Flags().StringVarP(&username, "username", "u", "", "Username to generate the certificate for")
	WindowsCmd.Flags().StringVarP(&hostname, "hostname", "t", "", "Hostname to generate the certificate for")

	WindowsCmd.MarkFlagRequired("output-dir")
	WindowsCmd.MarkFlagRequired("username")
	WindowsCmd.MarkFlagRequired("hostname")
}
