package request

import (
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	Request struct {
		requestedAt time.RequestedAt
		route       Route
	}

	Route struct {
		remoteAddr RemoteAddr
	}

	RemoteAddr string
)

func NewRequest(requestedAt time.RequestedAt, remoteAddr RemoteAddr) Request {
	return Request{
		requestedAt: requestedAt,
		route: Route{
			remoteAddr: remoteAddr,
		},
	}
}

func (request Request) RequestedAt() time.RequestedAt {
	return request.requestedAt
}

func (request Request) Route() Route {
	return request.route
}

func (route Route) RemoteAddr() RemoteAddr {
	return route.remoteAddr
}
