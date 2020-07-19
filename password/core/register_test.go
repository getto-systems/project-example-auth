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

// パスワードを保存したら audit: password registered
func TestRegister(t *testing.T) {
	h := newRegisterTestHelper(t)
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("password")

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	h.assertResult(err, nil)
	h.assertLog(
		logger,
		h.log("register password"),
		h.logEmpty(),
		h.log("password registered"),
	)
	h.assertRegistered(db, raw)
}

// 空のパスワードは保存できない
func TestRegisterFailedWhenEmptyPassword(t *testing.T) {
	h := newRegisterTestHelper(t)
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("") // 空のパスワード

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	h.assertResult(err, errors.New("password is empty"))
	h.assertLog(
		logger,
		h.log("register password"),
		h.logFailed("register password failed", err),
		h.logEmpty(),
	)
	h.assertNotRegistered(db)
}

// 長いパスワードは保存できない
func TestRegisterFailedWhenLongPassword(t *testing.T) {
	h := newRegisterTestHelper(t)
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	h.assertResult(err, errors.New("password is too long"))
	h.assertLog(
		logger,
		h.log("register password"),
		h.logFailed("register password failed", err),
		h.logEmpty(),
	)
	h.assertNotRegistered(db)
}

// ギリギリの長さのパスワードは保存できる
func TestRegisterWhenLongPassword(t *testing.T) {
	h := newRegisterTestHelper(t)
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	h.assertResult(err, nil)
	h.assertLog(
		logger,
		h.log("register password"),
		h.logEmpty(),
		h.log("password registered"),
	)
	h.assertRegistered(db, raw)
}

// DB エラーの場合、登録は失敗する
func TestRegisterWithFailureDB(t *testing.T) {
	h := newRegisterTestHelper(t)
	pub, _, gen, logger := h.setup()
	failure_db := newRegisterTestFailureDB() // 登録で error を返す DB

	request, user := h.context()
	raw := password.RawPassword("password")

	registerer := newRegisterer(pub, failure_db, gen)
	err := registerer.register(request, user, raw)

	h.assertResult(err, errors.New("db error"))
	h.assertLog(
		logger,
		h.log("register password"),
		h.logFailed("register password failed", err),
		h.logEmpty(),
	)
}

type (
	registerTestFailureDB struct{}

	registerTestGenerator struct{}

	registerTestHelper struct {
		t   *testing.T
		gen registerTestGenerator

		request data.Request
		user    data.User
	}

	registerTestLogEntry struct {
		message string
		err     error
	}
)

func newRegisterTestFailureDB() registerTestFailureDB {
	return registerTestFailureDB{}
}

func (registerTestFailureDB) RegisterUserPassword(data.User, password.HashedPassword) error {
	return errors.New("db error")
}

func newRegisterTestGenerator() registerTestGenerator {
	return registerTestGenerator{}
}

func (registerTestGenerator) GeneratePassword(raw password.RawPassword) (password.HashedPassword, error) {
	return password.HashedPassword(raw), nil
}

func newRegisterTestHelper(t *testing.T) *registerTestHelper {
	gen := newRegisterTestGenerator()

	request := data.Request{}
	user := data.NewUser("register-test")

	return &registerTestHelper{
		t:   t,
		gen: gen,

		request: request,
		user:    user,
	}
}

func (h *registerTestHelper) setup() (password.EventPublisher, *db.MemoryStore, password.Generator, *testLogger) {
	pub := pubsub.NewPubSub()
	logger := newTestLogger()
	log := password_event_log.NewEventLogger(logger)
	pub.Subscribe(log)

	db := db.NewMemoryStore()

	return pub, db, h.gen, logger
}
func (h registerTestHelper) context() (data.Request, data.User) {
	return h.request, h.user
}

func (h registerTestHelper) assertResult(got error, expected error) {
	assertError(h.t, "register", got, expected)
}

func (h registerTestHelper) assertLog(
	logger *testLogger,
	debug registerTestLogEntry,
	info registerTestLogEntry,
	audit registerTestLogEntry,
) {
	h.assertLogEntry("debug", logger.debug, debug)
	h.assertLogEntry("info", logger.info, info)
	h.assertLogEntry("audit", logger.audit, audit)
}
func (h registerTestHelper) assertLogEntry(label string, got event_log.Entry, expected registerTestLogEntry) {
	if expected.message == "" {
		if got.Message != "" {
			h.t.Errorf("register %s log message is not expected: %s (empty expected)", label, got.Message)
		}
	} else {
		if got.Message != expected.message {
			h.t.Errorf("register %s log message is not expected: %s (expected: %s)", label, got.Message, expected.message)
		}
		if got.Request != h.request {
			h.t.Errorf("register %s log request is not expected: %v (expected: %v)", label, got.Request, h.request)
		}
		if *got.User != h.user {
			h.t.Errorf("register %s log user is not expected: %v (expected: %v)", label, *got.User, h.user)
		}
		if expected.err != nil {
			assertError(h.t, fmt.Sprintf("register %s log", label), got.Error, expected.err)
		}
	}
}

func (h registerTestHelper) log(message string) registerTestLogEntry {
	return registerTestLogEntry{
		message: message,
	}
}
func (h registerTestHelper) logFailed(message string, err error) registerTestLogEntry {
	return registerTestLogEntry{
		message: message,
		err:     err,
	}
}
func (h registerTestHelper) logEmpty() registerTestLogEntry {
	return registerTestLogEntry{}
}

func (h registerTestHelper) assertRegistered(db *db.MemoryStore, expected password.RawPassword) {
	hashed, _ := h.gen.GeneratePassword(expected)

	got, ok := db.GetUserPassword(h.user)
	if !ok {
		h.t.Error("password not registered")
	} else {
		if string(got) != string(hashed) {
			h.t.Errorf("registered password not match: %s (expected: %s)", got, hashed)
		}
	}
}
func (h registerTestHelper) assertNotRegistered(db *db.MemoryStore) {
	_, ok := db.GetUserPassword(h.user)
	if ok {
		h.t.Error("password registered")
	}
}
