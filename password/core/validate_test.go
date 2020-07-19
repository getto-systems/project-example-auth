package core

import (
	"fmt"
	"strings"

	"github.com/getto-systems/project-example-id/password/db"
	password_event_log "github.com/getto-systems/project-example-id/password/event_log"
	"github.com/getto-systems/project-example-id/password/pubsub"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"

	"errors"
	"testing"
)

// パスワードが一致したら audit: authenticated by password
func TestValidate(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, user := h.context()
	raw := password.RawPassword("password") // 保存されているものと同じパスワード

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, nil)
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.log("authenticated by password"),
	)
}

// パスワードが一致しなかったら audit: validate password failed
func TestValidateFailed(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, user := h.context()
	raw := password.RawPassword("different-password") // 保存されているものと違うパスワード

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, errors.New("password not matched"))
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.logFailed("validate password failed", err),
	)
}

// パスワードが見つからない場合は認証失敗
func TestValidateFailedWhenPasswordNotFound(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	//h.registerPassword(db, password.HashedPassword("password")) // パスワードの登録はしない

	request, user := h.context()
	raw := password.RawPassword("password") // このユーザーのパスワードは登録されていない

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, errors.New("user password not found"))
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.logFailed("validate password failed", err),
	)
}

// 空のパスワードの場合、必ず失敗する
func TestValidateFailedWhenEmptyPassword(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, user := h.context()
	raw := password.RawPassword("") // 空のパスワード

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, errors.New("password is empty"))
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.logFailed("validate password failed", err),
	)
}

// 長いパスワードの場合、必ず失敗する
func TestValidateFailedWhenLongPassword(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword("password"))

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, errors.New("password is too long"))
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.logFailed("validate password failed", err),
	)
}

// ギリギリの長さのパスワードの場合、成功する
func TestValidateWhenLongPassword(t *testing.T) {
	h := newValidateTestHelper(t)
	pub, db, matcher, logger := h.setup()
	h.registerPassword(db, password.HashedPassword(strings.Repeat("a", 72)))

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	validator := newValidator(pub, db, matcher)
	err := validator.validate(request, user, raw)

	h.assertResult(err, nil)
	h.assertLog(
		logger,
		h.log("validate password"),
		h.logEmpty(),
		h.log("authenticated by password"),
	)
}

type (
	validateTestMatcher struct{}

	validateTestHelper struct {
		t       *testing.T
		matcher validateTestMatcher

		request data.Request
		user    data.User
	}

	validateTestLogEntry struct {
		message string
		err     error
	}
)

func newValidateTestMatcher() validateTestMatcher {
	return validateTestMatcher{}
}

func (validateTestMatcher) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) error {
	if string(raw) != string(hashed) {
		return errors.New("password not matched")
	}
	return nil
}

func newValidateTestHelper(t *testing.T) *validateTestHelper {
	matcher := newValidateTestMatcher()

	request := data.Request{}
	user := data.NewUser("validate-test")

	return &validateTestHelper{
		t:       t,
		matcher: matcher,

		request: request,
		user:    user,
	}
}

func (h *validateTestHelper) setup() (password.EventPublisher, *db.MemoryStore, password.Matcher, *testLogger) {
	pub := pubsub.NewPubSub()
	logger := newTestLogger()
	log := password_event_log.NewEventLogger(logger)
	pub.Subscribe(log)

	db := db.NewMemoryStore()

	return pub, db, h.matcher, logger
}
func (h validateTestHelper) context() (data.Request, data.User) {
	return h.request, h.user
}

func (h validateTestHelper) assertResult(got error, expected error) {
	assertError(h.t, "validate", got, expected)
}

func (h validateTestHelper) assertLog(
	logger *testLogger,
	debug validateTestLogEntry,
	info validateTestLogEntry,
	audit validateTestLogEntry,
) {
	h.assertLogEntry("debug", logger.debug, debug)
	h.assertLogEntry("info", logger.info, info)
	h.assertLogEntry("audit", logger.audit, audit)
}
func (h validateTestHelper) assertLogEntry(label string, got event_log.Entry, expected validateTestLogEntry) {
	if expected.message == "" {
		if got.Message != "" {
			h.t.Errorf("validate %s log message is not expected: %s (empty expected)", label, got.Message)
		}
	} else {
		if got.Message != expected.message {
			h.t.Errorf("validate %s log message is not expected: %s (expected: %s)", label, got.Message, expected.message)
		}
		if got.Request != h.request {
			h.t.Errorf("validate %s log request is not expected: %v (expected: %v)", label, got.Request, h.request)
		}
		if *got.User != h.user {
			h.t.Errorf("validate %s log user is not expected: %v (expected: %v)", label, *got.User, h.user)
		}
		if expected.err != nil {
			assertError(h.t, fmt.Sprintf("validate %s log", label), got.Error, expected.err)
		}
	}
}

func (h validateTestHelper) log(message string) validateTestLogEntry {
	return validateTestLogEntry{
		message: message,
	}
}
func (h validateTestHelper) logFailed(message string, err error) validateTestLogEntry {
	return validateTestLogEntry{
		message: message,
		err:     err,
	}
}
func (h validateTestHelper) logEmpty() validateTestLogEntry {
	return validateTestLogEntry{}
}

func (h validateTestHelper) registerPassword(db *db.MemoryStore, password password.HashedPassword) {
	db.RegisterUserPassword(h.user, password)
}
