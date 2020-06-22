package http_handler

import (
	"github.com/getto-systems/project-example-id/data"

	"time"
)

func RequestedAt() data.RequestedAt {
	return data.RequestedAt(time.Now().UTC())
}

type Logger interface {
	Debugf(*data.Request, string, ...interface{})
}
