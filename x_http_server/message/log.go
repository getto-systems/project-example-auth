package message

import (
	"encoding/json"
	"log"
	"os"
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

func (log LogMessage) Send(message string) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	log.logger.Println(data)

	return nil
}
