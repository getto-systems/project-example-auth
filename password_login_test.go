package _usecase

import (
	"strings"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func ExamplePasswordLogin_login_renew_logout() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	h.newRequest("Renew", minute(1), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	h.newRequest("Logout", minute(2), logoutHandler, func() {
		NewLogout(u).Logout(logoutHandler)
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

func ExamplePasswordLogin_loginWithNoApiRoles() {
	h := newPasswordLoginTestHelper()
	h.registerUserDataWithoutApiRoles("user-id", "login-id", "password") // ユーザーを登録(ApiRoles なし)

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: nil
	// credential: expires: "2020-01-01T00:05:00Z", roles: []
	//
}

func ExamplePasswordLogin_log() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
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
	// log: "Ticket/Register/TryToRegister", debug
	// log: "Ticket/Register/Register", info
	// log: "Credential/IssueTicketToken/TryToIssue", debug
	// log: "Credential/IssueTicketToken/Issue", info
	// log: "Credential/IssueApiToken/TryToIssue", debug
	// log: "Credential/IssueApiToken/Issue", info
	// log: "Credential/IssueContentToken/TryToIssue", debug
	// log: "Credential/IssueContentToken/Issue", info
	//
}

func ExampleRenew_log() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Renew", minute(1), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Ticket/Extend/TryToExtend", debug
	// log: "Ticket/Extend/Extend", info
	// log: "Credential/IssueTicketToken/TryToIssue", debug
	// log: "Credential/IssueTicketToken/Issue", info
	// log: "Credential/IssueApiToken/TryToIssue", debug
	// log: "Credential/IssueApiToken/Issue", info
	// log: "Credential/IssueContentToken/TryToIssue", debug
	// log: "Credential/IssueContentToken/Issue", info
	//
}

func ExampleLogout_log() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Logout", minute(1), logoutHandler, func() {
		NewLogout(u).Logout(logoutHandler)
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
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Ticket/Deactivate/TryToDeactivate", debug
	// log: "Ticket/Deactivate/Deactivate", info
	//
}

func ExamplePasswordLogin_failedBecausePasswordNotFound() {
	h := newPasswordLoginTestHelper()
	h.registerOnlyUserAndLogin("user-id", "login-id") // ユーザーとログインだけ登録してパスワードは登録しない

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// パスワードの登録なしでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "invalid-password-login"
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 違うパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "different-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "invalid-password-login"
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 空のパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "invalid-password-login"
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 長いパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", longPassword)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printCredential()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// request: "2020-01-01T00:00:00Z"
	// err: "invalid-password-login"
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 長いパスワードでログインを試みる
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", longPassword)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
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
	// log: "Ticket/Register/TryToRegister", debug
	// log: "Ticket/Register/Register", info
	// log: "Credential/IssueTicketToken/TryToIssue", debug
	// log: "Credential/IssueTicketToken/Issue", info
	// log: "Credential/IssueApiToken/TryToIssue", debug
	// log: "Credential/IssueApiToken/Issue", info
	// log: "Credential/IssueContentToken/TryToIssue", debug
	// log: "Credential/IssueContentToken/Issue", info
	//
}

func ExampleRenew_limited() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 最大延長期間まで延長可能
	h.newRequest("Renew", minute(5), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 有効期限切れ
	h.newRequest("Renew", minute(6), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

func ExampleRenew_failedBecauseAlreadyLogout() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 認証情報を取っておく
	credential := h.session.credential

	// Renew する前にログアウトしてしまう
	h.newRequest("Logout", minute(1), logoutHandler, func() {
		NewLogout(u).Logout(logoutHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 取っておいた認証情報でログインを試みる
	h.session.credential = credential

	// ログアウト済み
	h.newRequest("Renew", minute(2), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

func ExampleRenew_failedBecauseDifferentNonce() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 違う Nonce を使用して認証を試みる
	h.setNonce("different-nonce")

	h.newRequest("Renew", minute(2), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/FailedToParseBecauseNonceMatchFailed", audit
	//
}

func ExampleRenew_failedBecauseTicketNotFound() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// Nonce を別なものにして認証を試みる
	h.setNonce("another-nonce")
	h.setCredentialNonce("another-nonce")

	h.newRequest("Renew", minute(2), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseTicketNotFound", audit
	//
}

func ExampleRenew_failedBecauseDifferentUser() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	renewHandler := newRenewHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// User を別なものにして認証を試みる
	h.setCredentialUser(user.NewUser("another-user-id"))

	h.newRequest("Renew", minute(2), renewHandler, func() {
		NewRenew(u).Renew(renewHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseUserMatchFailed", audit
	//
}

func ExampleLogout_failedBecauseAlreadyExpired() {
	h := newPasswordLoginTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録済みデータと同じログインID・パスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	logoutHandler := newLogoutHandler(handler)

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("Logout", minute(10), logoutHandler, func() {
		NewLogout(u).Logout(logoutHandler)
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
	// err: "invalid-ticket"
	// credential: nil
	// log: "Credential/ParseTicketSignature/TryToParse", debug
	// log: "Credential/ParseTicketSignature/Parse", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

type (
	passwordLoginTestHelper struct {
		*testInfra
	}

	passwordLoginTestHandler struct {
		*commonTestHandler

		login    user.Login
		password password.RawPassword
	}
)

func newPasswordLoginTestHelper() passwordLoginTestHelper {
	return passwordLoginTestHelper{
		testInfra: newTestInfra(),
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
