package http_handler

import (
	"github.com/getto-systems/project-example-id/basic"

	"time"
)

func Now() basic.Time {
	return basic.Time(time.Now().UTC())
}
