package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type ExpirationParam struct {
	Expires     data.Second
	ExtendLimit data.Second
}

func NewExpiration(param ExpirationParam) Expiration {
	return Expiration{
		expires:     param.Expires,
		extendLimit: param.ExtendLimit,
	}
}

func (exp Expiration) Expires(request data.Request) data.Expires {
	return request.RequestedAt().Expires(exp.expires)
}

func (exp Expiration) ExtendLimit(request data.Request) data.ExtendLimit {
	return request.RequestedAt().ExtendLimit(exp.extendLimit)
}
