package client

import (
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
)

func ExamplePasswordChange_getLogin_change() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printCredential()
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", time.Minute(1), passwordChangeHandler, func() {
		NewPasswordChange(client).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLogin(passwordChangeHandler.login)
	})

	h.newRequest("PasswordChange/Change", time.Minute(2), passwordChangeHandler, func() {
		NewPasswordChange(client).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// credential: expires: "2020-01-01T00:05:00Z", roles: [role]
	// err: nil
	//
	// PasswordChange/GetLogin
	// request: "2020-01-01T00:01:00Z"
	// err: nil
	// login: {login-id}
	//
	// PasswordChange/Change
	// request: "2020-01-01T00:02:00Z"
	// err: nil
	//
}

func ExamplePasswordChange_getLoginLog() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", time.Minute(1), passwordChangeHandler, func() {
		NewPasswordChange(client).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/GetLogin
	// err: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "User/GetLogin/TryToGetLogin", debug
	// log: "User/GetLogin/GetLogin", info
	//
}

func ExamplePasswordChange_changeLog() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", time.Minute(1), passwordChangeHandler, func() {
		NewPasswordChange(client).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/Change", time.Minute(2), passwordChangeHandler, func() {
		NewPasswordChange(client).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/GetLogin
	// err: nil
	//
	// PasswordChange/Change
	// err: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/AuthByPassword", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/Change", audit
	//
}

/*
func ExamplePasswordChange_disableOldPassword() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", time.Minute(1), passwordChangeHandler, func() {
		NewPasswordChange(client).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/Change", time.Minute(2), passwordChangeHandler, func() {
		NewPasswordChange(client).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 前のパスワードでログインを試みる
	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// credential: expires: "2020-01-01T00:05:00Z", roles: [role]
	// err: nil
	//
	// PasswordChange/GetLogin
	// request: "2020-01-01T00:01:00Z"
	// err: nil
	// login: {login-id}
	//
	// PasswordChange/Change
	// request: "2020-01-01T00:02:00Z"
	// err: nil
	//
}
*/

type (
	passwordChangeTestHelper struct {
		*testBackend
	}

	passwordChangeTestHandler struct {
		*commonTestHandler

		oldPassword password.RawPassword
		newPassword password.RawPassword

		login user.Login
	}
)

func newPasswordChangeTestHelper() passwordChangeTestHelper {
	return passwordChangeTestHelper{
		testBackend: newTestBackend(),
	}
}

func newPasswordChangeHandler(handler *commonTestHandler, oldPassword password.RawPassword, newPassword password.RawPassword) *passwordChangeTestHandler {
	return &passwordChangeTestHandler{
		commonTestHandler: handler,

		oldPassword: oldPassword,
		newPassword: newPassword,
	}
}

func (handler *passwordChangeTestHandler) handler() PasswordChangeHandler {
	return handler
}
func (handler *passwordChangeTestHandler) GetLoginRequest() (request.Request, error) {
	return handler.newRequest(), nil
}
func (handler *passwordChangeTestHandler) GetLoginResponse(login user.Login, err error) {
	handler.setError(err)
	handler.login = login
}
func (handler *passwordChangeTestHandler) ChangeRequest() (request.Request, password.ChangeParam, error) {
	return handler.newRequest(), password.ChangeParam{
		OldPassword: handler.oldPassword,
		NewPassword: handler.newPassword,
	}, nil
}
func (handler *passwordChangeTestHandler) ChangeResponse(err error) {
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
