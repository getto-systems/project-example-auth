package core

import (
	"fmt"
	"strings"

	"github.com/getto-systems/project-example-id/password/log"
	repository_password "github.com/getto-systems/project-example-id/password/repository/password"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/event_log"
	"github.com/getto-systems/project-example-id/password"
)

// ログインID を取得
func Example_getLogin() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()
	h.registerLogin(passwords) // ログインID を登録

	request, user := h.context()

	registerer := newRegisterer(logger, passwords, gen)
	login, err := registerer.getLogin(request, user)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("login: %s\n", formatLogin(&login))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: nil
	// login: {register-login}
	// debug: ["Password/Register/TryToGetLogin", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: []
}

// パスワードを取得
func Example_getLogin_fail_LoginNotFound() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()
	//h.registerLogin(passwords) // ログインID を登録しない

	request, user := h.context()

	registerer := newRegisterer(logger, passwords, gen)
	login, err := registerer.getLogin(request, user)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("login: %s\n", formatLogin(&login))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))

	// Output:
	// err: "Password/Password/NotFound/Login"
	// login: {}
	// debug: ["Password/Register/TryToGetLogin", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["Password/Register/FailedToGetLogin", req: {register-remote}, user: {register-user}, err: "Password/Password/NotFound/Login"]
	// audit: []
}

// パスワードを保存したら audit: password registered
func Example_register() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("password")

	registerer := newRegisterer(logger, passwords, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))
	fmt.Printf("db: %s\n", h.formatDB(passwords))

	// Output:
	// err: nil
	// debug: ["Password/Register/TryToRegister", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: ["Password/Register/Registered", req: {register-remote}, user: {register-user}, err: nil]
	// db: "password"
}

// 空のパスワードは保存できない
func Example_register_fail_EmptyPassword() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("") // 空のパスワード

	registerer := newRegisterer(logger, passwords, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))
	fmt.Printf("db: %s\n", h.formatDB(passwords))

	// Output:
	// err: "Password/Password/Empty"
	// debug: ["Password/Register/TryToRegister", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["Password/Register/FailedToRegister", req: {register-remote}, user: {register-user}, err: "Password/Password/Empty"]
	// audit: []
	// db: nil
}

// 長いパスワードは保存できない
func Example_register_fail_LongPassword() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 73)) // 長いパスワード

	registerer := newRegisterer(logger, passwords, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))
	fmt.Printf("db: %s\n", h.formatDB(passwords))

	// Output:
	// err: "Password/Password/TooLong"
	// debug: ["Password/Register/TryToRegister", req: {register-remote}, user: {register-user}, err: nil]
	// info: ["Password/Register/FailedToRegister", req: {register-remote}, user: {register-user}, err: "Password/Password/TooLong"]
	// audit: []
	// db: nil
}

// ギリギリの長さのパスワードは保存できる
func Example_register_LongPassword() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()

	request, user := h.context()
	raw := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトまで許容

	registerer := newRegisterer(logger, passwords, gen)
	err := registerer.register(request, user, raw)

	fmt.Printf("err: %s\n", formatError(err))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit))
	fmt.Printf("db: %s\n", h.formatDB(passwords))

	// Output:
	// err: nil
	// debug: ["Password/Register/TryToRegister", req: {register-remote}, user: {register-user}, err: nil]
	// info: []
	// audit: ["Password/Register/Registered", req: {register-remote}, user: {register-user}, err: nil]
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

func newRegisterTestGenerator() (gen registerTestGenerator) {
	return
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

func (h registerTestHelper) setup() (password.Logger, *repository_password.MemoryStore, password.PasswordGenerator, *testLogger) {
	testLogger := newTestLogger()
	logger := log.NewLogger(testLogger)

	passwords := repository_password.NewMemoryStore()

	return logger, passwords, h.gen, testLogger
}

func (h registerTestHelper) registerLogin(passwords *repository_password.MemoryStore) {
	passwords.RegisterLogin(h.user, h.login)
}

func (h registerTestHelper) context() (data.Request, data.User) {
	return h.request, h.user
}

func (h registerTestHelper) formatLog(entry event_log.Entry) string {
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

func (h registerTestHelper) formatDB(passwords *repository_password.MemoryStore) string {
	password, ok := passwords.GetUserPassword(h.user)
	if !ok {
		return "nil"
	} else {
		return fmt.Sprintf("\"%s\"", password)
	}
}
