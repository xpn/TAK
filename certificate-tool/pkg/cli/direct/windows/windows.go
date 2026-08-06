package cli

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	authserver "xpnsec.com/certificate-tool/v2/pkg/api"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
	"xpnsec.com/shared/v2/pkg/connection"
)

var outputDir string
var username string
var hostname string
var connectionOptions connection.Options

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

		// If Client Cert and Client Key provided, we actually request the certificate is signed!
		if cmd.Flags().Changed("client-cert") {
			client, err := authserver.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
			if err != nil {
				fmt.Printf("[!] Error creating auth server client: %v\n", err)
				return
			}
			certs, err := client.GenerateWindowsHostCertificate(context.Background(), cert)
			if err != nil {
				fmt.Printf("[!] Error generating host certificate: %v\n", err)
				return
			}

			tlsPubPath := path.Join(outputDir, "windows_tls_signed.crt")

			os.WriteFile(tlsPubPath, certs.Cert, 0644)

			fmt.Printf("[*] Signed certificate saved to %s\n", tlsPubPath)

			fmt.Println("[?] To convert the certificate to DER format, use the following command:")
			fmt.Printf("  openssl x509 -in %s -outform DER -out %s\n", tlsPubPath, strings.TrimSuffix(tlsPubPath, ".crt")+".der")

		}
	},
}

func init() {
	WindowsCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to save generated certificates")
	WindowsCmd.Flags().StringVarP(&username, "username", "u", "", "Username to generate the certificate for")
	WindowsCmd.Flags().StringVarP(&hostname, "hostname", "t", "", "Hostname to generate the certificate for")
	connectionOptions.AddProxyFlag(WindowsCmd.Flags())
	connectionOptions.AddClientCredentialFlags(WindowsCmd.Flags())
	connectionOptions.AddClusterNameFlag(WindowsCmd.Flags())
	WindowsCmd.MarkFlagsRequiredTogether("client-cert", "client-key", "proxy", "cluster-name")

	WindowsCmd.MarkFlagRequired("output-dir")
	WindowsCmd.MarkFlagRequired("username")
	WindowsCmd.MarkFlagRequired("hostname")
}
