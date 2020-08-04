package client

import (
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func ExamplePasswordLogin_login_renew_logout() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	h.newRequest("Renew", time.Minute(1), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	h.newRequest("Logout", time.Minute(2), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: nil
	// credential: expires: "2020-01-01T00:05:00Z", roles: [role]
	//
	// Renew
	// request: "2020-01-01T00:01:00Z"
	// err: nil
	// credential: expires: "2020-01-01T00:06:00Z", roles: [role]
	//
	// Logout
	// request: "2020-01-01T00:02:00Z"
	// err: nil
	// credential: nil
	//
}

func ExamplePasswordLogin_loginLog() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/AuthByPassword", audit
	// log: "Ticket/Issue/TryToIssue", debug
	// log: "Ticket/Issue/Issue", info
	// log: "ApiToken/IssueApiToken/TryToIssue", debug
	// log: "ApiToken/IssueApiToken/Issue", info
	// log: "ApiToken/IssueContentToken/TryToIssue", debug
	// log: "ApiToken/IssueContentToken/Issue", info
	//
}

func ExamplePasswordLogin_renewLog() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Renew", time.Minute(1), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Renew
	// err: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Ticket/Extend/TryToExtend", debug
	// log: "Ticket/Extend/Extend", info
	// log: "ApiToken/IssueApiToken/TryToIssue", debug
	// log: "ApiToken/IssueApiToken/Issue", info
	// log: "ApiToken/IssueContentToken/TryToIssue", debug
	// log: "ApiToken/IssueContentToken/Issue", info
	//
}

func ExamplePasswordLogin_logoutLog() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Renew", time.Minute(1), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Logout", time.Minute(2), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Renew
	// err: nil
	//
	// Logout
	// err: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Ticket/Shrink/TryToShrink", debug
	// log: "Ticket/Shrink/Shrink", info
	//
}

type (
	passwordLoginTestHelper struct {
		*testBackend
	}

	passwordLoginTestHandler struct {
		*commonTestHandler

		login    user.Login
		password password.RawPassword
	}
)

func newPasswordLoginTestHelper() passwordLoginTestHelper {
	return passwordLoginTestHelper{
		testBackend: newTestBackend(),
	}
}

func newPasswordLoginHandler(handler *commonTestHandler, loginID user.LoginID, rawPassword password.RawPassword) passwordLoginTestHandler {
	return passwordLoginTestHandler{
		commonTestHandler: handler,

		login:    user.NewLogin(loginID),
		password: rawPassword,
	}
}

func (handler passwordLoginTestHandler) handler() PasswordLoginHandler {
	return handler
}
func (handler passwordLoginTestHandler) LoginRequest() (request.Request, user.Login, password.RawPassword, error) {
	return handler.newRequest(), handler.login, handler.password, nil
}
func (handler passwordLoginTestHandler) LoginResponse(err error) {
	handler.setError(err)
}

/*
// パスワードを取得
func Example_getLogin_fail_LoginNotFound() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()
	//h.registerLogin(passwords) // ログインID を登録しない

	request, user := h.context()

	registerer := newRegisterer(logger, passwords, gen)
	login, err := registerer.getLogin(request, user)

	fmt.Println(formatError(err, errPasswordNotFoundLogin))
	fmt.Println(h.formatLogin(login))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug, nil))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info, errPasswordNotFoundLogin))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit, nil))

	// Output:
	// err
	// login: {}
	// debug: ["Password/Register/TryToGetLogin", req, user, err: nil]
	// info: ["Password/Register/FailedToGetLogin", req, user, err]
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

	fmt.Println(formatError(err, nil))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug, nil))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info, nil))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit, nil))
	fmt.Println(h.formatDB(passwords, &raw))

	// Output:
	// err: nil
	// debug: ["Password/Register/TryToRegister", req, user, err: nil]
	// info: []
	// audit: ["Password/Register/Registered", req, user, err: nil]
	// db
}

// 空のパスワードは保存できない
func Example_register_fail_EmptyPassword() {
	h := newRegisterTestHelper()
	logger, passwords, gen, testLogger := h.setup()

	request, user := h.context()
	raw := password.RawPassword("") // 空のパスワード

	registerer := newRegisterer(logger, passwords, gen)
	err := registerer.register(request, user, raw)

	fmt.Println(formatError(err, errPasswordEmpty))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug, nil))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info, errPasswordEmpty))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit, nil))
	fmt.Println(h.formatDB(passwords, nil))

	// Output:
	// err
	// debug: ["Password/Register/TryToRegister", req, user, err: nil]
	// info: ["Password/Register/FailedToRegister", req, user, err]
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

	fmt.Println(formatError(err, errPasswordTooLong))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug, nil))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info, errPasswordTooLong))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit, nil))
	fmt.Println(h.formatDB(passwords, nil))

	// Output:
	// err
	// debug: ["Password/Register/TryToRegister", req, user, err: nil]
	// info: ["Password/Register/FailedToRegister", req, user, err]
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

	fmt.Println(formatError(err, nil))
	fmt.Printf("debug: %s\n", h.formatLog(testLogger.debug, nil))
	fmt.Printf("info: %s\n", h.formatLog(testLogger.info, nil))
	fmt.Printf("audit: %s\n", h.formatLog(testLogger.audit, nil))
	fmt.Println(h.formatDB(passwords, &raw))

	// Output:
	// err: nil
	// debug: ["Password/Register/TryToRegister", req, user, err: nil]
	// info: []
	// audit: ["Password/Register/Registered", req, user, err: nil]
	// db
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

func (h registerTestHelper) formatLogin(login password.Login) string {
	return formatLogin(&login, &h.login)
}

func (h registerTestHelper) formatLog(entry event_log.Entry, err error) string {
	if entry.Message == "" {
		return "[]"
	}

	return fmt.Sprintf(
		"[\"%s\", %s, %s, %s]",
		entry.Message,
		formatRequest(entry.Request, h.request),
		formatUser(entry.User, &h.user),
		formatError(entry.Error, err),
	)
}

func (h registerTestHelper) formatDB(passwords *repository_password.MemoryStore, expected *password.RawPassword) string {
	password, ok := passwords.GetUserPassword(h.user)
	if !ok {
		return "db: nil"
	}

	if expected == nil || string(password) != string(*expected) {
		return fmt.Sprintf("db: \"%s\"", password)
	}

	return "db"
}
*/
