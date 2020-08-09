package password_reset

import (
	gotime "time"

	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/request"
)

type (
	Expires gotime.Time

	Expiration struct {
		expires time.Second
	}
)

func NewExpiration(second time.Second) Expiration {
	return Expiration{
		expires: second,
	}
}

func (exp Expiration) Expires(request request.Request) Expires {
	return NewExpires(request.RequestedAt(), exp.expires)
}

func NewExpires(requestedAt time.RequestedAt, second time.Second) Expires {
	return Expires(addSecond(requestedAt, second))
}
func addSecond(requestedAt time.RequestedAt, second time.Second) gotime.Time {
	duration := gotime.Duration(second * 1_000_000_000)
	return gotime.Time(requestedAt).Add(duration)
}

func (expires Expires) Expired(requestedAt time.RequestedAt) bool {
	return gotime.Time(expires).Before(gotime.Time(requestedAt))
}
