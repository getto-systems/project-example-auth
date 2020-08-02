package password_reset

import (
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

type (
	SessionID string
	Token     string

	Session struct {
		id SessionID
	}

	SessionData struct {
		user        user.User
		login       user.Login
		requestedAt time.RequestedAt
		expires     time.Expires
	}

	Destination struct {
		// TODO 手抜きだ！
		Type string

		//SlackChannel string
		//MailAddress string
	}

	// TODO 見直す
	Status struct {
		// 宛先情報、待ち、送信中、完了、失敗のステータス
	}
)

func NewSession(id SessionID) Session {
	return Session{
		id: id,
	}
}
func (session Session) ID() SessionID {
	return session.id
}

func NewSessionData(user user.User, login user.Login, requestedAt time.RequestedAt, expires time.Expires) SessionData {
	return SessionData{
		user:        user,
		login:       login,
		requestedAt: requestedAt,
		expires:     expires,
	}
}
func (data SessionData) Data() (user.User, user.Login, time.RequestedAt, time.Expires) {
	return data.user, data.login, data.requestedAt, data.expires
}
func (data SessionData) User() user.User {
	return data.user
}
func (data SessionData) Login() user.Login {
	return data.login
}
func (data SessionData) Expires() time.Expires {
	return data.expires
}

func NewLogDestination() Destination {
	return Destination{
		Type: "Log",
	}
}

// TODO 見直す
func NewStatus() Status {
	return Status{}
}
