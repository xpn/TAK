package cli

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	"github.com/gravitational/teleport/api/utils/keys"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	authserver "xpnsec.com/certificate-tool/v2/pkg/api"
	"xpnsec.com/certificate-tool/v2/pkg/cert"
)

var outputDir string
var clientKeyPath, clientCertPath string
var nodeName string

var SSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A certificate multi-tool for generating certificates required by Teleport",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if err := os.MkdirAll(outputDir, 0700); err != nil {
			fmt.Printf("[!] Error creating output dir: %x\n", err)
			return
		}

		privateKey, err := cert.GenerateSSHKey()
		if err != nil {
			fmt.Printf("[!] Error generating private key: %v\n", err)
			return
		}

		der, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			fmt.Printf("[!] Error generating private key: %v\n", err)
			return
		}

		privPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: der,
		})

		sshPub, err := ssh.NewPublicKey(privateKey.Public())
		if err != nil {
			fmt.Printf("[!] Error generating public key: %v\n", err)
			return
		}

		tlsPub, err := keys.MarshalPublicKey(privateKey.Public())
		if err != nil {
			fmt.Printf("[!] Error generating public key: %v\n", err)
			return
		}

		publicSSH := ssh.MarshalAuthorizedKey(sshPub)

		publicSSHPath := path.Join(outputDir, "host_ssh.pub")
		privPEMPath := path.Join(outputDir, "host.key")
		tlsPubPath := path.Join(outputDir, "host_tls.pub")

		os.WriteFile(publicSSHPath, publicSSH, 0644)
		os.WriteFile(privPEMPath, privPEM, 0600)
		os.WriteFile(tlsPubPath, tlsPub, 0644)

		fmt.Printf("[*] Certificates generated:\n\t%s\n\t%s\n\t%s\n", publicSSHPath, privPEMPath, tlsPubPath)

		// If Client Cert and Client Key provided, we actually request the certificate is signed!
		if clientCertPath != "" && clientKeyPath != "" {
			client, err := authserver.NewClient(clientCertPath, clientKeyPath, "10.1.10.1:8443")
			if err != nil {
				fmt.Printf("[!] Error creating auth server client: %v\n", err)
				return
			}
			certs, err := client.GenerateSSHHostCertificate(context.Background(), nodeName, publicSSH, tlsPub)
			if err != nil {
				fmt.Printf("[!] Error generating host certificate: %v\n", err)
				return
			}

			tlsPubPath := path.Join(outputDir, "host_tls_signed.crt")
			sshPubPath := path.Join(outputDir, "host_ssh_signed.crt")

			os.WriteFile(tlsPubPath, certs.TLS, 0644)
			os.WriteFile(sshPubPath, certs.SSH, 0644)
		}
	},
}

func init() {
	SSHCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to save generated certificates")
	SSHCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "", "Existing Client Cert (makes gRPC call if included)")
	SSHCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "", "Existing Client Key (makes gRPC call if included)")
	SSHCmd.Flags().StringVarP(&nodeName, "node-name", "n", "", "Node name to use for certificate")
	SSHCmd.MarkFlagRequired("output-dir")
}
