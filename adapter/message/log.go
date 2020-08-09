package message

import (
	"encoding/json"
	"log"
	"os"

	"github.com/getto-systems/project-example-id/gateway/message"
)

type (
	LogMessage struct {
		logger *log.Logger
	}
)

func NewLogMessage() LogMessage {
	return LogMessage{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (log LogMessage) log(message string) message.LogMessage {
	return log
}

func (log LogMessage) Send(message string) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	log.logger.Println(data)

	return nil
}
