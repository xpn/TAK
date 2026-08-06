package cli

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	authserver "xpnsec.com/certificate-tool/v2/pkg/api"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
	"xpnsec.com/shared/v2/pkg/connection"
)

var outputDir string
var nodeName string
var username string
var appName string
var publicAddress string
var connectionOptions connection.Options

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
		if cmd.Flags().Changed("client-cert") {
			client, err := authserver.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
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
	AppCmd.Flags().StringVarP(&username, "username", "u", "", "Username to use for certificate")
	AppCmd.Flags().StringVarP(&appName, "app-name", "n", "", "App name to use for certificate")
	AppCmd.Flags().StringVarP(&publicAddress, "public-address", "a", "", "Public address to use for certificate")
	connectionOptions.AddProxyFlag(AppCmd.Flags())
	connectionOptions.AddClientCredentialFlags(AppCmd.Flags())
	connectionOptions.AddClusterNameFlag(AppCmd.Flags())
	AppCmd.MarkFlagsRequiredTogether("client-cert", "client-key", "proxy", "cluster-name", "username", "app-name", "public-address")

	AppCmd.MarkFlagRequired("output-dir")
}
