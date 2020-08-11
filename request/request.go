package request

import (
	"time"
)

type (
	Request struct {
		requestedAt RequestedAt
		route       Route
	}

	RequestedAt time.Time

	Route struct {
		remoteAddr RemoteAddr
	}

	RemoteAddr string
)

func RequestedAtNow() RequestedAt {
	return RequestedAt(time.Now().UTC())
}

func NewRequest(requestedAt RequestedAt, remoteAddr RemoteAddr) Request {
	return Request{
		requestedAt: requestedAt,
		route: Route{
			remoteAddr: remoteAddr,
		},
	}
}

func (request Request) RequestedAt() RequestedAt {
	return request.requestedAt
}

func (request Request) Route() Route {
	return request.route
}

func (route Route) RemoteAddr() RemoteAddr {
	return route.remoteAddr
}
