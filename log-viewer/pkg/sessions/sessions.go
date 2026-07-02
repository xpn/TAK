package sessions

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
	auditlogv1 "xpnsec.com/log-viewer/v2/pkg/grpc"
)

type SessionClient struct {
	client *auditlogv1.AuditLogServiceClient
	conn   *grpc.ClientConn
}

type SessionInfo struct {
	SessionId string        `json:"sessionId"`
	Time      time.Time     `json:"time"`
	Username  string        `json:"username"`
	Duration  time.Duration `json:"duration"`
	Hostname  string        `json:"hostname"`
}

func NewClient(certPath, keyPath, target string) (*SessionClient, error) {
	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	creds := credentials.NewTLS(&tls.Config{
		NextProtos:         []string{"teleport-auth@6578616d706c652e636f6d.teleport.cluster.local", "h2"},
		Certificates:       []tls.Certificate{clientCertificate},
		InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(creds),
	)

	client := auditlogv1.NewAuditLogServiceClient(conn)
	return &SessionClient{client: &client, conn: conn}, nil
}

func (c *SessionClient) ListRecordedSessions() ([]SessionInfo, error) {
	sessionInfo := make([]SessionInfo, 0)

	unstructuredEvents, err := (*c.client).GetUnstructuredEvents(context.Background(), &auditlogv1.GetUnstructuredEventsRequest{
		EventTypes: []string{"session.end"},
		StartDate:  timestamppb.New(time.Now().Add(time.Hour * -100)),
		EndDate:    timestamppb.New(time.Now()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get unstructured events: %w", err)
	}

	for _, event := range unstructuredEvents.Items {
		startTime := event.Unstructured.Fields["session_start"].GetStringValue()
		endTime := event.Unstructured.Fields["session_stop"].GetStringValue()
		parsedStartTime, err := time.Parse(time.RFC3339Nano, startTime)
		if err != nil {
			continue
		}

		parsedEndTime, err := time.Parse(time.RFC3339Nano, endTime)
		if err != nil {
			continue
		}
		duration := parsedEndTime.Sub(parsedStartTime)

		sessionInfo = append(sessionInfo, SessionInfo{
			SessionId: event.Unstructured.Fields["sid"].GetStringValue(),
			Username:  event.Unstructured.Fields["login"].GetStringValue(),
			Time:      event.Time.AsTime(),
			Duration:  duration,
			Hostname:  event.Unstructured.Fields["server_hostname"].GetStringValue(),
		})
	}
	return sessionInfo, nil
}

func (c *SessionClient) DumpRecordedSession(sessionId string, replay bool) error {

	streamingClient, err := (*c.client).StreamUnstructuredSessionEvents(context.Background(), &auditlogv1.StreamUnstructuredSessionEventsRequest{
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}

	player := NewSessionPlayer()
	channel := make(chan auditlogv1.EventUnstructured, 1)
	if replay {
		go player.Play(channel)
	}

	for {
		event, err := streamingClient.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if replay {
			channel <- *event
		}

		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		// if event.Type == "desktop.recording" {
		// 	return nil
		// }
		// if event.Type == "print" {
		// 	data := event.Unstructured.Fields["data"].GetStringValue()
		// 	dataBase64Decoded, err := base64.StdEncoding.DecodeString(data)
		// 	if err != nil {
		// 		fmt.Printf("Error decoding base64 data: %v\n", err)
		// 		continue
		// 	}

		// 	fmt.Printf("%s", dataBase64Decoded)
		// }
	}
}
