package ticket

import (
	"github.com/getto-systems/project-example-id/data"
)

type (
	Ticket       []byte
	ApiToken     []byte
	ContentToken interface {
		Policy() string
		Signature() string
	}

	Nonce string

	Expiration struct {
		expires     data.Second
		extendLimit data.Second
	}
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
