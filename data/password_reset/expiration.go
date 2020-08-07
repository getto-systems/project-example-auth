package password_reset

import (
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	Expiration struct {
		expires time.Second
	}
)

func NewExpiration(second time.Second) Expiration {
	return Expiration{
		expires: second,
	}
}
func (exp Expiration) Expires(requestedAt time.RequestedAt) time.Expires {
	return requestedAt.Expires(exp.expires)
}
