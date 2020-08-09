package request

import (
	"time"

	"github.com/getto-systems/project-example-id/misc/expiration"
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
func (request Request) Expired(expires expiration.Expires) bool {
	return expires.Expired(time.Time(request.requestedAt))
}
func (request Request) Extend(limit expiration.ExtendLimit, second expiration.ExpireSecond) expiration.Expires {
	return expiration.Extend(time.Time(request.RequestedAt()), limit, second)
}

func (request Request) NewExpires(second expiration.ExpireSecond) expiration.Expires {
	return expiration.NewExpires(time.Time(request.RequestedAt()), second)
}
func (request Request) NewExtendLimit(second expiration.ExtendSecond) expiration.ExtendLimit {
	return expiration.NewExtendLimit(time.Time(request.RequestedAt()), second)
}

func (request Request) Route() Route {
	return request.route
}

func (route Route) RemoteAddr() RemoteAddr {
	return route.remoteAddr
}
