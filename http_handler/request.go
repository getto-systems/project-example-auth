package http_handler

import (
	"net/http"

	"github.com/getto-systems/project-example-id/data"
)

func Request(r *http.Request) data.Request {
	return data.NewRequest(data.Now(), data.RemoteAddr(r.RemoteAddr))
}
