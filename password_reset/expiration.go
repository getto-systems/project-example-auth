package password_reset

import (
	"github.com/getto-systems/project-example-id/data/time"
)

type (
	expiration struct {
		expires time.Second
	}
)

func newExpiration(second time.Second) expiration {
	return expiration{
		expires: second,
	}
}
func (exp expiration) Expires(requestedAt time.RequestedAt) time.Expires {
	return requestedAt.Expires(exp.expires)
}
