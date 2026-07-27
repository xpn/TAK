package cli

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	authserver "xpnsec.com/certificate-tool/v2/pkg/api"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
)

var outputDir string
var username string
var clientCertPath string
var clientKeyPath string
var clusterName string
var proxyAddress string

var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database certificate utilities",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if err := os.MkdirAll(outputDir, 0700); err != nil {
			fmt.Printf("[!] Error creating output dir: %x\n", err)
			return
		}

		cert, key, err := cert.GenerateDatabaseCSR(username)
		if err != nil {
			fmt.Printf("[!] Error generating certificate: %v\n", err)
			return
		}

		csrPath := path.Join(outputDir, fmt.Sprintf("%s.csr", username))
		keyPath := path.Join(outputDir, fmt.Sprintf("%s.key", username))

		// Write the CSR and key to files
		os.WriteFile(csrPath, cert, 0744)
		os.WriteFile(keyPath, key, 0600)

		fmt.Printf("[*] Certificates generated:\n\t%s\n\t%s\n", csrPath, keyPath)

		// If Client Cert and Client Key provided, we actually request the certificate is signed!
		if clientCertPath != "" && clientKeyPath != "" {
			client, err := authserver.NewClient(clientCertPath, clientKeyPath, proxyAddress, clusterName)
			if err != nil {
				fmt.Printf("[!] Error creating auth server client: %v\n", err)
				return
			}
			certs, err := client.GenerateDatabaseCertificate(context.Background(), cert)
			if err != nil {
				fmt.Printf("[!] Error generating host certificate: %v\n", err)
				return
			}

			tlsPubPath := path.Join(outputDir, "database_tls_signed.crt")

			os.WriteFile(tlsPubPath, certs.Cert, 0644)
		}
	},
}

func init() {
	DatabaseCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to save generated certificates")
	DatabaseCmd.Flags().StringVarP(&username, "username", "u", "", "Username to generate the certificate for")
	DatabaseCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "", "Existing Client Cert (makes gRPC call if included)")
	DatabaseCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "", "Existing Client Key (makes gRPC call if included)")
	DatabaseCmd.Flags().StringVarP(&clusterName, "cluster-name", "l", "", "Cluster name to use for certificate")
	DatabaseCmd.Flags().StringVarP(&proxyAddress, "proxy-address", "a", "", "Proxy address:port to use for certificate")

	DatabaseCmd.MarkFlagRequired("output-dir")
	DatabaseCmd.MarkFlagRequired("username")
}
