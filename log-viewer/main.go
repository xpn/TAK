package main

import (
	"github.com/fatih/color"
	"xpnsec.com/log-viewer/v2/pkg/cli"
)

// func dumpRecordedSession(client *auditlogv1.AuditLogServiceClient, sessionData string) error {

// 	streamingClient, err := (*client).StreamUnstructuredSessionEvents(context.Background(), &auditlogv1.StreamUnstructuredSessionEventsRequest{
// 		SessionId: sessionData,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		event, err := streamingClient.Recv()
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 		if event.Type == "desktop.recording" {
// 			return nil
// 		}
// 		if event.Type == "print" {
// 			data := event.Unstructured.Fields["data"].GetStringValue()
// 			dataBase64Decoded, err := base64.StdEncoding.DecodeString(data)
// 			if err != nil {
// 				fmt.Printf("Error decoding base64 data: %v\n", err)
// 				continue
// 			}

// 			fmt.Printf("%s", dataBase64Decoded)
// 		}
// 	}
// }

func main() {

	theme1 := color.RGB(247, 71, 130)
	theme2 := color.RGB(90, 142, 255)
	theme1.Println(`
_    ____ ____    _  _ _ ____ _ _ _ ____ ____
|    |  | | __ __ |  | | |___ | | | |___ |__/
|___ |__| |__]     \/  | |___ |_|_| |___ |  \`)
	theme2.Println(`           @_xpn_
`)

	cli.Execute()

	// certPath := "/tmp/user-node.crt"
	// keyPath := "/tmp/user-node.key"

	// // Load the client certificate and key
	// cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	// if err != nil {
	// 	fmt.Printf("Failed to load client cert: %v\n", err)
	// 	return
	// }

	// // Create TLS credentials
	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{cert},
	// }
	// credentials.NewTLS(tlsConfig)

	// conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	fmt.Printf("Failed to create gRPC client: %v\n", err)
	// 	return
	// }
	// defer conn.Close()

	// timestamp, err := time.Parse("2006-01-02T15:04:05Z", "2025-02-04T12:28:48.091Z")
	// if err != nil {
	// 	fmt.Printf("Error parsing time: %v\n", err)
	// 	return
	// }
	// startTimestamp := timestamppb.New(timestamp) // 24 hours ago
	// endTimestamp := timestamppb.Now()

	// client := auditlogv1.NewAuditLogServiceClient(conn)
	// unstructuredEvents, err := client.GetUnstructuredEvents(context.Background(), &auditlogv1.GetUnstructuredEventsRequest{
	// 	StartDate:  startTimestamp,
	// 	EndDate:    endTimestamp,
	// 	EventTypes: []string{"session.upload"},
	// })
	// if err != nil {
	// 	fmt.Printf("Error getting unstructured events: %v\n", err)
	// 	return
	// }

	// for _, event := range unstructuredEvents.Items {
	// 	fmt.Printf("Event: %v\n", event.Type)
	// 	if event.Type == "session.upload" {
	// 		fmt.Printf("Session Start Event Details: %v\n", event)
	// 		//sessionId := event.Unstructured.Fields["sid"].GetStringValue()
	// 		/*err := dumpRecordedSession(&client, sessionId)
	// 		if err != nil {
	// 			fmt.Printf("Error dumping recorded session: %v\n", err)
	// 		}*/
	// 	}
	// }

	// dumpRecordedSession(&client, "fd78d00c-29a6-46e4-92ab-dc0992e7e229")

	// //fmt.Printf("Unstructured Events: %v\n", unstructuredEvents)

}
