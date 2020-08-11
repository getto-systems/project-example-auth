package request

import (
	"time"
)

type (
	Request struct {
		requestedAt RequestedAt
		route       Route
	}
	RequestLog struct {
		RequestedAt string   `json:"requested_at"`
		Route       RouteLog `json:"route"`
	}

	RequestedAt time.Time

	Route struct {
		remoteAddr RemoteAddr
	}
	RouteLog struct {
		RemoteAddr string `json:"remote_addr"`
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

func (request Request) Log() RequestLog {
	return RequestLog{
		RequestedAt: time.Time(request.RequestedAt()).String(),
		Route:       request.route.Log(),
	}
}
func (route Route) Log() RouteLog {
	return RouteLog{
		RemoteAddr: string(route.remoteAddr),
	}
}
