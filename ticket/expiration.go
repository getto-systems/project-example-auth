package ticket

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	expiration struct {
		expires     time.Second
		extendLimit time.Second
	}
)

func newExpiration(param ticket.ExpirationParam) expiration {
	return expiration{
		expires:     param.Expires,
		extendLimit: param.ExtendLimit,
	}
}
func (exp expiration) Expires(request request.Request) time.Expires {
	return request.RequestedAt().Expires(exp.expires)
}
func (exp expiration) ExtendLimit(request request.Request) time.ExtendLimit {
	return request.RequestedAt().ExtendLimit(exp.extendLimit)
}
