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
var clientKeyPath, clientCertPath string
var nodeName string
var username string
var appName string
var publicAddress string
var proxyAddress string
var clusterName string

var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "A certificate multi-tool for generating certificates required by Teleport",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if err := os.MkdirAll(outputDir, 0700); err != nil {
			fmt.Printf("[!] Error creating output dir: %x\n", err)
			return
		}

		cert, key, err := cert.GenerateAppPublicPrivateKey()
		if err != nil {
			fmt.Printf("[!] Error generating keys: %v\n", err)
			return
		}

		pubPath := path.Join(outputDir, "app.pub")
		keyPath := path.Join(outputDir, "app.key")

		// Write the CSR and key to files
		os.WriteFile(pubPath, cert, 0744)
		os.WriteFile(keyPath, key, 0600)

		fmt.Printf("[*] Keys generated:\n\t%s\n\t%s\n", pubPath, keyPath)

		// If Client Cert and Client Key provided, we actually request the certificate is signed!
		if clientCertPath != "" && clientKeyPath != "" {
			client, err := authserver.NewClient(clientCertPath, clientKeyPath, proxyAddress, clusterName)
			if err != nil {
				fmt.Printf("[!] Error creating auth server client: %v\n", err)
				return
			}
			certs, err := client.GenerateUserAppCertificate(context.Background(), username, appName, publicAddress, cert)
			if err != nil {
				fmt.Printf("[!] Error signing certificate: %v\n", err)
				return
			}

			tlsPubPath := path.Join(outputDir, "app_tls_signed.crt")

			os.WriteFile(tlsPubPath, certs.Cert, 0644)
		}
	},
}

func init() {
	AppCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to save generated certificates")
	AppCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "", "Existing Client Cert (makes gRPC call if included)")
	AppCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "", "Existing Client Key (makes gRPC call if included)")
	AppCmd.Flags().StringVarP(&username, "username", "u", "", "Username to use for certificate")
	AppCmd.Flags().StringVarP(&appName, "app-name", "n", "", "App name to use for certificate")
	AppCmd.Flags().StringVarP(&publicAddress, "public-address", "a", "", "Public address to use for certificate")
	AppCmd.Flags().StringVarP(&proxyAddress, "proxy-address", "p", "", "Proxy address:port to use for certificate")
	AppCmd.Flags().StringVarP(&clusterName, "cluster-name", "l", "", "Cluster name to use for certificate")

	AppCmd.MarkFlagRequired("output-dir")
}
