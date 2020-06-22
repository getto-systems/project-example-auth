package http_handler

import (
	"github.com/getto-systems/project-example-id/data"

	"time"
)

func RequestedAt() data.RequestedAt {
	return data.RequestedAt(time.Now().UTC())
}

type Logger interface {
	DebugMessage(*data.Request, string)
	DebugError(*data.Request, string, error)
}
