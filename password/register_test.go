package password

import (
	"errors"
	"strings"
	"testing"

	"github.com/getto-systems/project-example-id/data"
)

// パスワードを保存したら PasswordRegistered イベントが発行される
func TestRegister(t *testing.T) {
	raw := data.RawPassword("password")

	pub := newRegisterTestEventPublisher()
	db := newRegisterTestDB()
	gen := newRegisterTestGenerator()

	// register
	register := NewRegister(pub, db, gen)
	err := register.register(data.Request{}, data.User{}, raw)

	h := newRegisterTestHelper(t, pub, err)
	h.checkRegisterError(nil)
	h.checkRegisterPasswordEvent("fired")
	h.checkRegisterPasswordFailedEvent(nil)
	h.checkPasswordRegisteredEvent("fired")
}

// 空のパスワードは保存できない
func TestRegisterFailedWhenEmptyPassword(t *testing.T) {
	raw := data.RawPassword("")

	pub := newRegisterTestEventPublisher()
	db := newRegisterTestDB()
	gen := newRegisterTestGenerator()

	// register
	register := NewRegister(pub, db, gen)
	err := register.register(data.Request{}, data.User{}, raw)

	h := newRegisterTestHelper(t, pub, err)
	h.checkRegisterError(errors.New("password is empty"))
	h.checkRegisterPasswordEvent("fired")
	h.checkRegisterPasswordFailedEvent(errors.New("password is empty"))
	h.checkPasswordRegisteredEvent("never")
}

// 長いパスワードは保存できない
func TestRegisterFailedWhenLongPassword(t *testing.T) {
	length := 73
	raw := data.RawPassword(strings.Repeat("a", length))

	pub := newRegisterTestEventPublisher()
	db := newRegisterTestDB()
	gen := newRegisterTestGenerator()

	// register
	register := NewRegister(pub, db, gen)
	err := register.register(data.Request{}, data.User{}, raw)

	h := newRegisterTestHelper(t, pub, err)
	h.checkRegisterError(errors.New("password is too long"))
	h.checkRegisterPasswordEvent("fired")
	h.checkRegisterPasswordFailedEvent(errors.New("password is too long"))
	h.checkPasswordRegisteredEvent("never")
}

// ギリギリの長さのパスワードは保存できる
func TestRegisterWhenLongPassword(t *testing.T) {
	length := 72
	raw := data.RawPassword(strings.Repeat("a", length))

	pub := newRegisterTestEventPublisher()
	db := newRegisterTestDB()
	gen := newRegisterTestGenerator()

	// register
	register := NewRegister(pub, db, gen)
	err := register.register(data.Request{}, data.User{}, raw)

	h := newRegisterTestHelper(t, pub, err)
	h.checkRegisterError(nil)
	h.checkRegisterPasswordEvent("fired")
	h.checkRegisterPasswordFailedEvent(nil)
	h.checkPasswordRegisteredEvent("fired")
}

type (
	registerTestEventPublisher struct {
		registerPassword       string
		registerPasswordFailed error
		passwordRegistered     string
	}

	registerTestDB struct{}

	registerTestGenerator struct{}

	registerTestHelper struct {
		t   *testing.T
		pub *registerTestEventPublisher
		err error
	}
)

func newRegisterTestEventPublisher() *registerTestEventPublisher {
	return &registerTestEventPublisher{
		registerPassword:   "never",
		passwordRegistered: "never",
	}
}

func (pub *registerTestEventPublisher) RegisterPassword(request data.Request, user data.User) {
	pub.registerPassword = "fired"
}
func (pub *registerTestEventPublisher) RegisterPasswordFailed(request data.Request, user data.User, err error) {
	pub.registerPasswordFailed = err
}
func (pub *registerTestEventPublisher) PasswordRegistered(request data.Request, user data.User) {
	pub.passwordRegistered = "fired"
}

func newRegisterTestDB() registerTestDB {
	return registerTestDB{}
}

func (registerTestDB) RegisterUserPassword(data.User, data.HashedPassword) error {
	return nil
}

func newRegisterTestGenerator() registerTestGenerator {
	return registerTestGenerator{}
}

func (registerTestGenerator) GeneratePassword(raw data.RawPassword) (data.HashedPassword, error) {
	return data.HashedPassword(raw), nil
}

func newRegisterTestHelper(t *testing.T, pub *registerTestEventPublisher, err error) registerTestHelper {
	return registerTestHelper{
		t:   t,
		pub: pub,
		err: err,
	}
}

func (h registerTestHelper) checkRegisterError(err error) {
	if err == nil {
		if h.err != nil {
			h.t.Errorf("register fired: %s", h.err)
		}
	} else {
		if h.err == nil {
			h.t.Error("register success")
		} else {
			if h.err.Error() != err.Error() {
				h.t.Errorf("register error message is not matched: %s (expected: %s)", h.err, err)
			}
		}
	}
}
func (h registerTestHelper) checkRegisterPasswordEvent(event string) {
	if h.pub.registerPassword != event {
		h.t.Errorf("RegisterPassword event not match: %s (expected: %s)", h.pub.registerPassword, event)
	}
}
func (h registerTestHelper) checkRegisterPasswordFailedEvent(err error) {
	if err == nil {
		if h.pub.registerPasswordFailed != nil {
			h.t.Errorf("RegisterPasswordFailed event fired: %s", h.pub.registerPasswordFailed)
		}
	} else {
		if h.pub.registerPasswordFailed == nil {
			h.t.Error("RegisterPasswordFailed event not fired")
		} else {
			if h.pub.registerPasswordFailed.Error() != err.Error() {
				h.t.Errorf("RegisterPasswordFailed error message is not matched: %s (expected: %s)", h.pub.registerPasswordFailed, err)
			}
		}
	}
}
func (h registerTestHelper) checkPasswordRegisteredEvent(event string) {
	if h.pub.passwordRegistered != event {
		h.t.Errorf("PasswordRegistered event not match: %s (expected: %s)", h.pub.passwordRegistered, event)
	}
}
