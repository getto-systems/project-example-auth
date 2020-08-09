package password_reset

import (
	"time"

	"github.com/getto-systems/project-example-id/_misc/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
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
		requestedAt request.RequestedAt
		expires     expiration.Expires
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
		since WaitingSince
	}
	StatusSending struct {
		since SendingSince
	}
	StatusComplete struct {
		at CompleteAt
	}
	StatusFailed struct {
		at     FailedAt
		reason string
	}

	WaitingSince time.Time
	SendingSince time.Time
	CompleteAt   time.Time
	FailedAt     time.Time
)

func NewSession(id SessionID) Session {
	return Session{
		id: id,
	}
}
func (session Session) ID() SessionID {
	return session.id
}

func NewSessionData(user user.User, login user.Login, requestedAt request.RequestedAt, expires expiration.Expires) SessionData {
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
func (data SessionData) RequestedAt() request.RequestedAt {
	return data.requestedAt
}
func (data SessionData) Expires() expiration.Expires {
	return data.expires
}

func NewLogDestination() Destination {
	return Destination{
		Type: "Log",
	}
}

func NewStatusWaiting(requestedAt request.RequestedAt) Status {
	waiting := newStatusWaiting(requestedAt)
	return Status{
		waiting: &waiting,
	}
}
func NewStatusSending(requestedAt request.RequestedAt) Status {
	sending := newStatusSending(requestedAt)
	return Status{
		sending: &sending,
	}
}
func NewStatusComplete(requestedAt request.RequestedAt) Status {
	complete := newStatusComplete(requestedAt)
	return Status{
		complete: &complete,
	}
}
func NewStatusFailed(requestedAt request.RequestedAt, reason string) Status {
	failed := newStatusFailed(requestedAt, reason)
	return Status{
		failed: &failed,
	}
}

func newStatusWaiting(requestedAt request.RequestedAt) StatusWaiting {
	return StatusWaiting{since: WaitingSince(requestedAt)}
}
func newStatusSending(requestedAt request.RequestedAt) StatusSending {
	return StatusSending{since: SendingSince(requestedAt)}
}
func newStatusComplete(requestedAt request.RequestedAt) StatusComplete {
	return StatusComplete{at: CompleteAt(requestedAt)}
}
func newStatusFailed(requestedAt request.RequestedAt, reason string) StatusFailed {
	return StatusFailed{at: FailedAt(requestedAt), reason: reason}
}

func (status Status) Waiting() bool {
	return status.waiting != nil
}
func (status Status) WaitingSince() (_ WaitingSince) {
	if status.waiting == nil {
		return
	}
	return status.waiting.since
}
func (status Status) Sending() bool {
	return status.sending != nil
}
func (status Status) SendingSince() (_ SendingSince) {
	if status.sending == nil {
		return
	}
	return status.sending.since
}
func (status Status) Complete() bool {
	return status.complete != nil
}
func (status Status) CompleteAt() (_ CompleteAt) {
	if status.complete == nil {
		return
	}
	return status.complete.at
}
func (status Status) Failed() bool {
	return status.failed != nil
}
func (status Status) FailedAtAndReason() (_ FailedAt, _ string) {
	if status.failed == nil {
		return
	}
	return status.failed.at, status.failed.reason
}
