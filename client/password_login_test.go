package client

import (
	"strings"

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

func ExamplePasswordLogin_loginFailed_emptyPassword() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "", []string{}) // 空のパスワードで登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 空のパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "")

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "Password.Check/Length.Empty"
	// credential: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordCheckFailed", info
	//
}

func ExamplePasswordLogin_loginFailed_tooLongPassword() {
	longPassword := password.RawPassword(strings.Repeat("a", 73)) // 長すぎるパスワード(72 バイトまで)

	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", longPassword, []string{}) // 長いパスワードで登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 長いパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", longPassword)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "Password.Check/Length.TooLong"
	// credential: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordCheckFailed", info
	//
}

func ExamplePasswordLogin_loginSuccess_longPassword() {
	longPassword := password.RawPassword(strings.Repeat("a", 72)) // 72 バイトのパスワード

	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", longPassword, []string{}) // 長いパスワードで登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 長いパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", longPassword)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
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

func ExamplePasswordLogin_renew_limited() {
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

	// 最大延長期間まで延長可能
	h.newRequest("Renew", time.Minute(5), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Renew
	// request: "2020-01-01T00:05:00Z"
	// err: nil
	// credential: expires: "2020-01-01T00:08:00Z", roles: [role]
	//
}

func ExamplePasswordLogin_renew_failed_alreadyExpired() {
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

	// 有効期限切れ
	h.newRequest("Renew", time.Minute(6), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Renew
	// request: "2020-01-01T00:06:00Z"
	// err: "Ticket.Validate/AlreadyExpired"
	// credential: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

func ExamplePasswordLogin_renew_failed_alreadyLogout() {
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

	// 認証情報を取っておく
	credential := h.credential

	// ログアウト
	h.newRequest("Logout", time.Minute(1), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.credential = credential

	// ログアウト済み
	h.newRequest("Renew", time.Minute(2), renewHandler, func() {
		NewRenew(client).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Logout
	// err: nil
	//
	// Renew
	// request: "2020-01-01T00:02:00Z"
	// err: "Ticket.Validate/AlreadyExpired"
	// credential: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
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
