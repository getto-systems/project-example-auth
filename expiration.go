package main

import (
	"github.com/getto-systems/project-example-id/data"
)

type expiration struct {
	expires     data.Second
	extendLimit data.Second
}

func (exp expiration) Expires(request data.Request) data.Expires {
	return request.RequestedAt().Expires(exp.expires)
}

func (exp expiration) ExtendLimit(request data.Request) data.ExtendLimit {
	return request.RequestedAt().ExtendLimit(exp.extendLimit)
}
