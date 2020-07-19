package core

/*
import (
	"testing"
	"time"

	"github.com/getto-systems/project-example-id/data"

	"errors"
)

// チケットが検証されたら AuthenticatedByTicket イベントが発行される
func TestValidate(t *testing.T) {
	signUser := data.NewUser("user-id")
	request, expires := validRequestAndExpires()

	pub := newValidateTestEventPublisher()
	signer := newValidateTestSigner(signUser, expires)

	// validate
	validater := NewValidater(pub, signer)
	user, err := validater.validate(request, Ticket{}, Nonce(""))

	h := newValidateTestHelper(t, pub, user, err)
	h.checkUser(signUser)
	h.checkValidateError(nil)
	h.checkValidateTicketEvent("fired")
	h.checkValidateTicketFailedEvent(nil)
	h.checkAuthenticatedByTicketEvent("fired")
}

// 署名が検証できなかったら ValidateTicketFailed イベントが発行される
func TestValidateFailed(t *testing.T) {
	pub := newValidateTestEventPublisher()
	signer := newValidateTestInvalidSigner()

	// validate
	validater := NewValidater(pub, signer)
	user, err := validater.validate(data.Request{}, Ticket{}, Nonce(""))

	h := newValidateTestHelper(t, pub, user, err)
	h.checkValidateError(errors.New("ticket parse error"))
	h.checkValidateTicketEvent("fired")
	h.checkValidateTicketFailedEvent(errors.New("ticket parse error"))
	h.checkAuthenticatedByTicketEvent("never")
}

// リクエスト時刻が有効期限を過ぎていたら ValidateTicketFailed イベントが発行される
func TestValidateFailedWhenExpiredTicket(t *testing.T) {
	signUser := data.NewUser("user-id")
	request, expires := invalidRequestAndExpires()

	pub := newValidateTestEventPublisher()
	signer := newValidateTestSigner(signUser, expires)

	// validate
	validater := NewValidater(pub, signer)
	user, err := validater.validate(request, Ticket{}, Nonce(""))

	h := newValidateTestHelper(t, pub, user, err)
	h.checkInvalidUser(signUser)
	h.checkValidateError(errors.New("ticket already expired"))
	h.checkValidateTicketEvent("fired")
	h.checkValidateTicketFailedEvent(errors.New("ticket already expired"))
	h.checkAuthenticatedByTicketEvent("never")
}

type (
	validateTestEventPublisher struct {
		validateTicket        string
		validateTicketFailed  error
		authenticatedByTicket string
	}

	validateTestSigner struct {
		user    data.User
		expires data.Expires
	}

	validateTestInvalidSigner struct{}

	validateTestHelper struct {
		t    *testing.T
		pub  *validateTestEventPublisher
		user data.User
		err  error
	}
)

func validRequestAndExpires() (data.Request, data.Expires) {
	requestedAt := data.RequestedAt(time.Now())
	expires := requestedAt.Expires(data.Second(10))
	request := data.NewRequest(requestedAt, "")

	return request, expires
}

func invalidRequestAndExpires() (data.Request, data.Expires) {
	requestedAt := data.RequestedAt(time.Now())
	expires := requestedAt.Expires(data.Second(-10))
	request := data.NewRequest(requestedAt, "")

	return request, expires
}

func newValidateTestEventPublisher() *validateTestEventPublisher {
	return &validateTestEventPublisher{
		validateTicket:        "never",
		authenticatedByTicket: "never",
	}
}

func (pub *validateTestEventPublisher) ValidateTicket(request data.Request) {
	pub.validateTicket = "fired"
}
func (pub *validateTestEventPublisher) ValidateTicketFailed(request data.Request, err error) {
	pub.validateTicketFailed = err
}
func (pub *validateTestEventPublisher) AuthenticatedByTicket(request data.Request, user data.User) {
	pub.authenticatedByTicket = "fired"
}

func newValidateTestSigner(user data.User, expires data.Expires) validateTestSigner {
	return validateTestSigner{
		user:    user,
		expires: expires,
	}
}

func (signer validateTestSigner) Parse(Ticket) (Nonce, data.User, data.Expires, error) {
	return Nonce(""), signer.user, signer.expires, nil
}

func (validateTestSigner) Sign(Nonce, data.User, data.Expires) (Ticket, error) {
	return nil, nil
}

func newValidateTestInvalidSigner() validateTestInvalidSigner {
	return validateTestInvalidSigner{}
}

func (validateTestInvalidSigner) Parse(Ticket) (Nonce, data.User, data.Expires, error) {
	return Nonce(""), data.User{}, data.Expires{}, errors.New("ticket parse error")
}

func (validateTestInvalidSigner) Sign(Nonce, data.User, data.Expires) (Ticket, error) {
	return nil, nil
}

func newValidateTestHelper(t *testing.T, pub *validateTestEventPublisher, user data.User, err error) validateTestHelper {
	return validateTestHelper{
		t:    t,
		pub:  pub,
		user: user,
		err:  err,
	}
}

func (h validateTestHelper) checkUser(user data.User) {
	if h.user != user {
		h.t.Errorf("different user: %v (expected: %v)", h.user, user)
	}
}
func (h validateTestHelper) checkInvalidUser(user data.User) {
	if h.user == user {
		h.t.Errorf("receive valid user: %v", h.user)
	}
}
func (h validateTestHelper) checkValidateError(err error) {
	if err == nil {
		if h.err != nil {
			h.t.Errorf("validate fired: %s", h.err)
		}
	} else {
		if h.err == nil {
			h.t.Error("validate success")
		} else {
			if h.err.Error() != err.Error() {
				h.t.Errorf("validate error message is not matched: %s (expected: %s)", h.err, err)
			}
		}
	}
}
func (h validateTestHelper) checkValidateTicketEvent(event string) {
	if h.pub.validateTicket != event {
		h.t.Errorf("ValidateTicket event not match: %s (expected: %s)", h.pub.validateTicket, event)
	}
}
func (h validateTestHelper) checkValidateTicketFailedEvent(err error) {
	if err == nil {
		if h.pub.validateTicketFailed != nil {
			h.t.Errorf("ValidateTicketFailed event fired: %s", h.pub.validateTicketFailed)
		}
	} else {
		if h.pub.validateTicketFailed == nil {
			h.t.Error("ValidateTicketFailed event not fired")
		} else {
			if h.pub.validateTicketFailed.Error() != err.Error() {
				h.t.Errorf("ValidateTicketFailed error message is not matched: %s (expected: %s)", h.pub.validateTicketFailed, err)
			}
		}
	}
}
func (h validateTestHelper) checkAuthenticatedByTicketEvent(event string) {
	if h.pub.authenticatedByTicket != event {
		h.t.Errorf("AuthenticatedByTicket event not match: %s (expected: %s)", h.pub.authenticatedByTicket, event)
	}
}
*/
