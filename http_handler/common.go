package http_handler

import (
	"time"
)

func Now() time.Time {
	return time.Now().UTC()
}
