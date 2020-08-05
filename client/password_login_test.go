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

func ExamplePasswordLogin_log() {
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

func ExampleRenew_log() {
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

func ExampleLogout_log() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Logout", time.Minute(1), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// Logout
	// err: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Ticket/Deactivate/TryToDeactivate", debug
	// log: "Ticket/Deactivate/Deactivate", info
	//
}

func ExamplePasswordLogin_failedBecausePasswordNotFound() {
	h := newPasswordLoginTestHelper()
	h.registerOnlyUserAndLogin("user-id", "login-id") // ユーザーとログインだけ登録してパスワードは登録しない

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// パスワードの登録なしでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

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
	// err: "Password.Validate/NotFound.Password"
	// credential: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordNotFound", audit
	//
}

func ExamplePasswordLogin_failedBecauseDifferentPassword() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{}) // ユーザー登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 違うパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "different-password")

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
	// err: "Password.Validate/MatchFailed"
	// credential: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordMatchFailed", audit
	//
}

func ExamplePasswordLogin_failedBecauseEmptyPassword() {
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

func ExamplePasswordLogin_failedBecauseTooLongPassword() {
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

func ExamplePasswordLogin_successWithLongPassword() {
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

func ExampleRenew_limited() {
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

func ExampleRenew_failedBecauseAlreadyExpired() {
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

func ExampleRenew_failedBecauseAlreadyLogout() {
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

	// Renew する前にログアウトしてしまう
	h.newRequest("Logout", time.Minute(1), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 取っておいた認証情報でログインを試みる
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

func ExampleRenew_failedBecauseDifferentNonce() {
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

	// 違う Nonce を使用して認証を試みる
	h.setNonce("different-nonce")

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
	// Renew
	// request: "2020-01-01T00:02:00Z"
	// err: "Ticket.Validate/DifferentNonce"
	// credential: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseDifferentInfo", audit
	//
}

func ExampleRenew_failedBecauseTicketNotFound() {
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

	// Nonce を別なものにして認証を試みる
	h.setNonce("another-nonce")
	h.setCredentialNonce("another-nonce")

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
	// Renew
	// request: "2020-01-01T00:02:00Z"
	// err: "Ticket.Validate/NotFound.Ticket"
	// credential: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseDifferentInfo", audit
	//
}

func ExampleRenew_failedBecauseDifferentUser() {
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

	// User を別なものにして認証を試みる
	h.setCredentialUser(user.NewUser("another-user-id"))

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
	// Renew
	// request: "2020-01-01T00:02:00Z"
	// err: "Ticket.Validate/DifferentUser"
	// credential: nil
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseDifferentInfo", audit
	//
}

func ExampleLogout_failedBecauseAlreadyExpired() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", time.Minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Logout", time.Minute(10), logoutHandler, func() {
		NewLogout(client).Logout(logoutHandler)
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
	// request: "2020-01-01T00:10:00Z"
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
