package http_handler

import (
	"github.com/getto-systems/project-example-id/basic"

	"time"
)

func RequestedAt() basic.RequestedAt {
	return basic.RequestedAt(time.Now().UTC())
}

type Logger interface {
	Debugf(basic.Request, string, ...interface{})
}
