package http_handler

import (
	"net/http"
	"time"

	"github.com/getto-systems/project-example-id/data"
)

func Request(r *http.Request) data.Request {
	return data.Request{
		RequestedAt: data.RequestedAt(time.Now().UTC()),
		Route: data.Route{
			RemoteAddr: data.RemoteAddr(r.RemoteAddr),
		},
	}
}
