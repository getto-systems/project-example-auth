package http_handler

import (
	"time"

	"github.com/getto-systems/project-example-id/data"
)

func RequestedAt() data.RequestedAt {
	return data.RequestedAt(time.Now().UTC())
}

type Logger interface {
	DebugMessage(*data.Request, string)
	DebugError(*data.Request, string, error)
}
