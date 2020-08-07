package ticket

import (
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	Expiration struct {
		expires     time.Second
		extendLimit time.Second
	}

	ExpirationParam struct {
		Expires     time.Second
		ExtendLimit time.Second
	}
)

func NewExpiration(param ExpirationParam) Expiration {
	return Expiration{
		expires:     param.Expires,
		extendLimit: param.ExtendLimit,
	}
}
func (exp Expiration) ExpireSecond() time.Second {
	return exp.expires
}
func (exp Expiration) Expires(request request.Request) time.Expires {
	return request.RequestedAt().Expires(exp.expires)
}
func (exp Expiration) ExtendLimit(request request.Request) time.ExtendLimit {
	return request.RequestedAt().ExtendLimit(exp.extendLimit)
}
