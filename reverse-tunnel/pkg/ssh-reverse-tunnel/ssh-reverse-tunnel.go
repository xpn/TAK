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
	target   string
	username string
	port     int
	cert     string
	key      string
}

func startSSHOutboundHeartbeat(sshConn ssh.Conn) {
	ch, req, err := sshConn.OpenChannel("teleport-heartbeat", nil)
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
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

func startSSHOutboundTransport(target string, sshConn ssh.Conn) (net.Conn, error) {

	transportChannel, transportOOBRequests, err := sshConn.OpenChannel("teleport-transport", nil)
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
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
		log.Fatalf("failed to send request: %v", err)
		return nil, err
	}

	if !ok {
		log.Fatalf("teleport-transport-dial request denied")
		errMessageBytes, _ := io.ReadAll(transportChannel.Stderr())
		errMessage := string(bytes.TrimSpace(errMessageBytes))
		if errMessage != "" {
			log.Fatalf("error from server: %s", errMessage)
		}
		os.Exit(1)
	}

	log.Printf("teleport-transport-dial request sent")

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

func handleIncomingTeleportDiscoveryChannels(newChannel ssh.NewChannel) {
	discoveryChannel, discoveryOOBRequests, err := newChannel.Accept()
	if err != nil {
		log.Printf("failed to accept teleport-discovery channel: %v", err)
		return
	}

	DiscardChannelData(discoveryChannel)

	go func() {
		for req := range discoveryOOBRequests {
			log.Printf("Received request on teleport-discovery channel: %s", req.Type)
			log.Printf("Payload: %s", string(req.Payload))
		}
	}()
}

func pipeToTcpConn(hostname string, conn ssh.Channel) {
	//if hostname == "@local-node" {
	hostname = fmt.Sprintf("localhost:%d", 23)
	//}

	tcpConn, err := net.Dial("tcp", hostname)
	if err != nil {
		log.Printf("failed to dial tcp connection to %s: %v", hostname, err)
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

func handleIncomingTeleportTransportChannels(newChannel ssh.NewChannel) {
	transportChannel, transportOOBRequests, err := newChannel.Accept()
	if err != nil {
		log.Printf("failed to accept teleport-transport channel: %v", err)
		return
	}

	//DiscardChannelData(transportChannel)

	go func() {
		for req := range transportOOBRequests {
			log.Printf("Received request on teleport-transport channel: %s", req.Type)
			switch req.Type {
			case "teleport-transport-dial":
				log.Printf("Payload: %s", string(req.Payload))
				go func() {
					pipeToTcpConn(string(req.Payload), transportChannel)
				}()
				req.Reply(true, nil)
			default:
				log.Printf("Unknown request type: %s", req.Type)
			}
		}
	}()
}

func handleNewIncomingChannels(newInboundChannels <-chan ssh.NewChannel) {
	go func() {
		for newChannel := range newInboundChannels {
			log.Printf("New channel requested: %v", newChannel.ChannelType())

			switch newChannel.ChannelType() {
			case "teleport-discovery":
				handleIncomingTeleportDiscoveryChannels(newChannel)
			case "teleport-transport":
				handleIncomingTeleportTransportChannels(newChannel)
			default:
				log.Printf("Unknown channel type: %s", newChannel.ChannelType())
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
	handleNewIncomingChannels(newInboundChannels)

	// Handle outbound channels
	startSSHOutboundHeartbeat(sshConn)
	_, err = startSSHOutboundTransport(s.target, sshConn)
	if err != nil {
		return err
	}

	return nil
}

func New(username string, target string, port int, certificatePath string, keyPath string) *SSHReverseTunnel {
	return &SSHReverseTunnel{
		target:   target,
		username: username,
		port:     port,
		cert:     certificatePath,
		key:      keyPath,
	}
}
