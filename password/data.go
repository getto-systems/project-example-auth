package password

import (
	"github.com/getto-systems/project-example-id/data"
)

type (
	LoginID string
	Login   struct {
		id LoginID
	}

	RawPassword    string
	HashedPassword []byte

	Password struct {
		user   data.User
		hashed HashedPassword
	}

	ResetToken string
	ResetID    string

	Reset struct {
		id ResetID
	}

	Expiration struct {
		expires data.Second
	}
)

func NewLogin(loginID LoginID) Login {
	return Login{
		id: loginID,
	}
}

func (login Login) ID() LoginID {
	return login.id
}

func NewPassword(user data.User, hashed HashedPassword) Password {
	return Password{
		user:   user,
		hashed: hashed,
	}
}

func (password Password) User() data.User {
	return password.user
}

func NewReset(resetID ResetID) Reset {
	return Reset{
		id: resetID,
	}
}

func (reset Reset) ID() ResetID {
	return reset.id
}

func NewExpiration(second data.Second) Expiration {
	return Expiration{
		expires: second,
	}
}

func (exp Expiration) Expires(request data.Request) data.Expires {
	return request.RequestedAt().Expires(exp.expires)
}
