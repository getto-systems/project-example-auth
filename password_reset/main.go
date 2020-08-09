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

	Status struct {
		waiting  *StatusWaiting
		sending  *StatusSending
		complete *StatusComplete
		failed   *StatusFailed
	}
	StatusWaiting struct {
		since time.Time
	}
	StatusSending struct {
		since time.Time
	}
	StatusComplete struct {
		at time.Time
	}
	StatusFailed struct {
		at     time.Time
		reason string
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
func (data SessionData) User() user.User {
	return data.user
}
func (data SessionData) Login() user.Login {
	return data.login
}
func (data SessionData) RequestedAt() time.RequestedAt {
	return data.requestedAt
}
func (data SessionData) Expires() time.Expires {
	return data.expires
}

func NewLogDestination() Destination {
	return Destination{
		Type: "Log",
	}
}

func NewStatusWaiting(since time.Time) Status {
	waiting := newStatusWaiting(since)
	return Status{
		waiting: &waiting,
	}
}
func NewStatusSending(since time.Time) Status {
	sending := newStatusSending(since)
	return Status{
		sending: &sending,
	}
}
func NewStatusComplete(at time.Time) Status {
	complete := newStatusComplete(at)
	return Status{
		complete: &complete,
	}
}
func NewStatusFailed(at time.Time, reason string) Status {
	failed := newStatusFailed(at, reason)
	return Status{
		failed: &failed,
	}
}

func newStatusWaiting(since time.Time) StatusWaiting {
	return StatusWaiting{since: since}
}
func newStatusSending(since time.Time) StatusSending {
	return StatusSending{since: since}
}
func newStatusComplete(at time.Time) StatusComplete {
	return StatusComplete{at: at}
}
func newStatusFailed(at time.Time, reason string) StatusFailed {
	return StatusFailed{at: at, reason: reason}
}

func (status Status) Waiting() bool {
	return status.waiting != nil
}
func (status Status) WaitingSince() (_ time.Time) {
	if status.waiting == nil {
		return
	}
	return status.waiting.since
}
func (status Status) Sending() bool {
	return status.sending != nil
}
func (status Status) SendingSince() (_ time.Time) {
	if status.sending == nil {
		return
	}
	return status.sending.since
}
func (status Status) Complete() bool {
	return status.complete != nil
}
func (status Status) CompleteAt() (_ time.Time) {
	if status.complete == nil {
		return
	}
	return status.complete.at
}
func (status Status) Failed() bool {
	return status.failed != nil
}
func (status Status) FailedAtAndReason() (_ time.Time, _ string) {
	if status.failed == nil {
		return
	}
	return status.failed.at, status.failed.reason
}
