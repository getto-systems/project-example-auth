package data

type (
	Request struct {
		requestedAt RequestedAt
		route       Route
	}

	Route struct {
		remoteAddr RemoteAddr
	}

	RemoteAddr string
)

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
