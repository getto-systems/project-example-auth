package password

import (
	"github.com/getto-systems/project-example-id/data"
)

var (
	ErrResetSessionNotFoundResetStatus  = newError("ResetSession/NotFound/ResetStatus")
	ErrResetSessionNotFoundResetSession = newError("ResetSession/NotFound/ResetSession")

	ErrResetSessionLoginNotMatched = newError("ResetSession/LoginNotMatched")
	ErrResetSessionAlreadyExpired  = newError("ResetSession/AlreadyExpired")
)

type (
	ResetToken     string
	ResetSessionID string

	ResetSession struct {
		id ResetSessionID
	}

	ResetSessionData struct {
		user        data.User
		login       Login
		requestedAt data.RequestedAt
		expires     data.Expires
	}

	// TODO 見直す
	ResetStatus struct {
	}

	ResetSessionExpiration struct {
		expires data.Second
	}

	ResetSessionRepository interface {
		FindResetStatus(ResetSession) (ResetStatus, error)     // 見つからない場合は ErrResetSessionNotFoundResetStatus
		FindResetSession(ResetToken) (ResetSessionData, error) // 見つからない場合は ErrResetSessionNotFoundResetSession

		RegisterResetSession(ResetSessionGenerator, ResetSessionData) (ResetSession, ResetToken, error)
	}

	ResetSessionGenerator interface {
		GenerateSession() (ResetSessionID, ResetToken, error)
	}
)

func NewResetSession(id ResetSessionID) ResetSession {
	return ResetSession{
		id: id,
	}
}

func (session ResetSession) ID() ResetSessionID {
	return session.id
}

func NewResetSessionData(
	user data.User,
	login Login,
	requestedAt data.RequestedAt,
	expires data.Expires,
) ResetSessionData {
	return ResetSessionData{
		user:        user,
		login:       login,
		requestedAt: requestedAt,
		expires:     expires,
	}
}

func (data ResetSessionData) Data() (data.User, Login, data.RequestedAt, data.Expires) {
	return data.user, data.login, data.requestedAt, data.expires
}

func (data ResetSessionData) User() data.User {
	return data.user
}

func (data ResetSessionData) Login() Login {
	return data.login
}

func (data ResetSessionData) Expires() data.Expires {
	return data.expires
}

func NewResetSessionExpiration(second data.Second) ResetSessionExpiration {
	return ResetSessionExpiration{
		expires: second,
	}
}

func (exp ResetSessionExpiration) Expires(requestedAt data.RequestedAt) data.Expires {
	return requestedAt.Expires(exp.expires)
}

// TODO 見直す
func NewResetStatus() ResetStatus {
	return ResetStatus{}
}
