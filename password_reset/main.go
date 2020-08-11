package password_reset

import (
	"time"

	"github.com/getto-systems/project-example-id/z_external/errors"
	"github.com/getto-systems/project-example-id/z_external/expiration"

	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

var (
	ErrCreateSessionNotFoundDestination = errors.NewError("PasswordReset.CreateSession", "NotFound.Destination")

	ErrGetStatusNotFoundSession  = errors.NewError("PasswordReset.GetStatus", "NotFound.Session")
	ErrGetStatusMatchFailedLogin = errors.NewError("PasswordReset.GetStatus", "MatchFailed.Login")

	ErrValidateNotFoundSession  = errors.NewError("PasswordReset.Validate", "NotFound.Session")
	ErrValidateMatchFailedLogin = errors.NewError("PasswordReset.Validate", "MatchFailed.Login")
	ErrValidateAlreadyExpired   = errors.NewError("PasswordReset.Validate", "AlreadyExpired")
	ErrValidateAlreadyClosed    = errors.NewError("PasswordReset.Validate", "AlreadyClosed")
)

type (
	SessionID string
	Token     string

	ExpireSecond expiration.ExpireSecond
	Expires      expiration.Expires

	Session struct {
		id SessionID
	}

	SessionData struct {
		user        user.User
		login       user.Login
		requestedAt request.RequestedAt
		expires     Expires
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

func NewExpires(request request.Request, second ExpireSecond) Expires {
	return Expires(expiration.NewExpires(
		time.Time(request.RequestedAt()),
		expiration.ExpireSecond(second),
	))
}
func (expires Expires) Expired(request request.Request) bool {
	return expiration.Expires(expires).Expired(time.Time(request.RequestedAt()))
}

func NewSessionData(user user.User, login user.Login, requestedAt request.RequestedAt, expires Expires) SessionData {
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
func (data SessionData) Expires() Expires {
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

func ExpireMinute(minutes int64) ExpireSecond {
	// セッションの有効期限は「分」の単位で設定するべき
	return ExpireSecond(minutes * 60)
}
