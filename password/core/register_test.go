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
)

// ログインID を取得
func Example_getLogin() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()
	h.registerLogin(db) // ログインID を登録

	request, user := h.context()

	registerer := newRegisterer(pub, db, gen)
	login, err := registerer.getLogin(request, user)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("login: %s\n", formatLogin(&login))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))

	// Output:
	// err: nil
	// login: {register-login}
	// debug: ["get login", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: []
}

// パスワードを取得
func Example_getLogin_fail_LoginNotFound() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()
	//h.registerLogin(db) // ログインID を登録しない

	request, user := h.context()

	registerer := newRegisterer(pub, db, gen)
	login, err := registerer.getLogin(request, user)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("login: %s\n", formatLogin(&login))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))

	// Output:
	// err: "login not found"
	// login: {}
	// debug: ["get login", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["get login failed", req: {register-remote}, user: {register-user}, err: "login not found"]
	// audit: []
}

// パスワードを保存したら audit: password registered
func Example_register() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("password")

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))
	fmt.Printf("db: %s\n", h.formatDB(db))

	// Output:
	// err: nil
	// debug: ["register password", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: ["password registered", req: {register-remote}, user: {register-user}, err: nil]
	// db: "password"
}

// 空のパスワードは保存できない
func Example_register_fail_EmptyPassword() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("") // 空のパスワード

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))
	fmt.Printf("db: %s\n", h.formatDB(db))

	// Output:
	// err: "password is empty"
	// debug: ["register password", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["register password failed", req: {register-remote}, user: {register-user}, err: "password is empty"]
	// audit: []
	// db: nil
}

// 長いパスワードは保存できない
func Example_register_fail_LongPassword() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))
	fmt.Printf("db: %s\n", h.formatDB(db))

	// Output:
	// err: "password is too long"
	// debug: ["register password", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["register password failed", req: {register-remote}, user: {register-user}, err: "password is too long"]
	// audit: []
	// db: nil
}

// ギリギリの長さのパスワードは保存できる
func Example_register_LongPassword() {
	h := newRegisterTestHelper()
	pub, db, gen, logger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	registerer := newRegisterer(pub, db, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", formatRegisterLog(logger.debug))
	fmt.Printf("info: %s\n", formatRegisterLog(logger.info))
	fmt.Printf("audit: %s\n", formatRegisterLog(logger.audit))
	fmt.Printf("db: %s\n", h.formatDB(db))

	// Output:
	// err: nil
	// debug: ["register password", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: ["password registered", req: {register-remote}, user: {register-user}, err: nil]
	// db: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}

type (
	registerTestGenerator struct{}

	registerTestHelper struct {
		gen registerTestGenerator

		request data.Request
		user    data.User
		login   password.Login
	}

	registerTestLogEntry struct {
		message string
		err     error
	}
)

func newRegisterTestGenerator() registerTestGenerator {
	return registerTestGenerator{}
}

func (registerTestGenerator) GeneratePassword(raw password.RawPassword) (password.HashedPassword, error) {
	return password.HashedPassword(raw), nil
}

func newRegisterTestHelper() registerTestHelper {
	gen := newRegisterTestGenerator()

	request := data.NewRequest(data.RequestedAt{}, data.RemoteAddr("register-remote"))
	user := data.NewUser("register-user")
	login := password.NewLogin("register-login")

	return registerTestHelper{
		gen: gen,

		request: request,
		user:    user,
		login:   login,
	}
}

func (h registerTestHelper) setup() (password.EventPublisher, *db.MemoryStore, password.Generator, *testLogger) {
	pub := pubsub.NewPubSub()
	logger := newTestLogger()
	log := password_event_log.NewEventLogger(logger)
	pub.Subscribe(log)

	db := db.NewMemoryStore()

	return pub, db, h.gen, logger
}

func (h registerTestHelper) registerLogin(db *db.MemoryStore) {
	db.RegisterLogin(h.user, h.login)
}

func (h registerTestHelper) context() (data.Request, data.User) {
	return h.request, h.user
}

func formatRegisterLog(entry event_log.Entry) string {
	if entry.Message == "" {
		return "[]"
	}

	return fmt.Sprintf(
		"[\"%s\", req: %s, user: %s, err: %s]",
		entry.Message,
		formatRequest(entry.Request),
		formatUser(entry.User),
		formatError(entry.Error),
	)
}

func (h registerTestHelper) formatDB(db *db.MemoryStore) string {
	password, ok := db.GetUserPassword(h.user)
	if !ok {
		return "nil"
	} else {
		return fmt.Sprintf("\"%s\"", password)
	}
}
