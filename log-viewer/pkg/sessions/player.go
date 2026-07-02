package sessions

import (
	"encoding/base64"
	"fmt"
	"time"

	auditlogv1 "xpnsec.com/log-viewer/v2/pkg/grpc"
)

type SessionPlayer struct {
}

func NewSessionPlayer() *SessionPlayer {
	return &SessionPlayer{}
}

func (s *SessionPlayer) Play(event <-chan auditlogv1.EventUnstructured) {
	startms := time.Now().UnixMilli()
	for e := range event {
		if e.Type == "session.end" || e.Type == "session.leave" {
			return
		}

		if e.Type != "print" {
			continue
		}

		targetms := int64(e.Unstructured.Fields["ms"].GetNumberValue())

		// If the duration hasn't elapsed yet, sleep for the remaining time
		if startms != 0 && targetms > startms {
			time.Sleep(time.Duration(targetms-startms) * time.Millisecond)
		}
		startms = targetms

		// Print the event
		stdoutData := e.Unstructured.Fields["data"].GetStringValue()

		// Base64 decode the data
		decodedData, err := base64.StdEncoding.DecodeString(stdoutData)
		if err != nil {
			fmt.Println("Error decoding data:", err)
		} else {
			fmt.Printf("%s", decodedData)
		}

	}
}
