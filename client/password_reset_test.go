package client

import (
	"errors"
	"log"
	"strings"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"
)

func ExamplePasswordReset_reset() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printResetSession(passwordResetHandler.session)
	})

	h.newRequest("PasswordReset/GetStatus", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printResetDestination(passwordResetHandler.dest)
		f.printResetStatus(passwordResetHandler.status)
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printMessage()
	})

	h.newRequest("PasswordReset/GetStatus", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printResetStatus(passwordResetHandler.status)
	})

	// メッセージからトークンを取得
	passwordResetHandler.fetchToken()

	h.newRequest("PasswordReset/Reset", time.Minute(4), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printResetToken(passwordResetHandler.token)
		f.printError()
		f.printCredential()
	})

	// Output:
	// PasswordReset/CreateSession
	// request: "2020-01-01T00:00:00Z"
	// err: nil
	// session: {reset-session-id}
	//
	// PasswordReset/GetStatus
	// request: "2020-01-01T00:01:00Z"
	// err: nil
	// dest: {Log}
	// status: {sending: {since: "2020-01-01T00:00:00Z"}}
	//
	// PasswordReset/SendToken
	// err: nil
	// message: password reset token: reset-token
	//
	// PasswordReset/GetStatus
	// request: "2020-01-01T00:03:00Z"
	// err: nil
	// status: {complete: {at: "2020-01-01T00:00:00Z"}}
	//
	// PasswordReset/Reset
	// request: "2020-01-01T00:04:00Z"
	// token: "reset-token"
	// err: nil
	// credential: expires: "2020-01-01T00:09:00Z", roles: [role]
	//
}

func ExamplePasswordReset_createSessionLog() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "PasswordReset/CreateSession/TryToCreateSession", debug
	// log: "PasswordReset/CreateSession/CreateSession", info
	// log: "PasswordReset/PushSendTokenJob/TryToPushSendTokenJob", debug
	// log: "PasswordReset/PushSendTokenJob/PushSendTokenJob", info
	//
}

func ExamplePasswordReset_sendTokenLog() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	// log: "PasswordReset/SendToken/TryToSendToken", debug
	// log: "PasswordReset/SendToken/SendToken", info
	//
}

func ExamplePasswordReset_getStatusLog() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: nil
	// log: "PasswordReset/GetStatus/TryToGetStatus", debug
	// log: "PasswordReset/GetStatus/GetStatus", info
	//
}

func ExamplePasswordReset_resetLog() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// メッセージからトークンを取得
	passwordResetHandler.fetchToken()

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: nil
	//
	// PasswordReset/Reset
	// err: nil
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/AuthByToken", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/Change", audit
	// log: "PasswordReset/CloseSession/TryToCloseSession", debug
	// log: "PasswordReset/CloseSession/CloseSession", info
	// log: "Ticket/Register/TryToRegister", debug
	// log: "Ticket/Register/Register", info
	// log: "Credential/IssueTicket/TryToIssueTicket", debug
	// log: "Credential/IssueTicket/IssueTicket", info
	// log: "Credential/IssueApiToken/TryToIssue", debug
	// log: "Credential/IssueApiToken/Issue", info
	// log: "Credential/IssueContentToken/TryToIssue", debug
	// log: "Credential/IssueContentToken/Issue", info
	//
}

func ExamplePasswordReset_disableOldPassword() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")
	// 古いパスワードでログイン
	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// メッセージからトークンを取得
	passwordResetHandler.fetchToken()

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 前のパスワードでログインを試みる
	h.newRequest("PasswordLogin", time.Minute(2), passwordLoginHandler, func() {
		NewPasswordLogin(client).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: nil
	//
	// PasswordLogin
	// err: "Password.Validate/MatchFailed"
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordMatchFailed", audit
	//
}

func ExamplePasswordReset_disableResetSession() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// メッセージからトークンを取得
	passwordResetHandler.fetchToken()

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// もう一度同じトークンを使用してリセットを試みる
	h.newRequest("PasswordReset/Reset", time.Minute(4), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: nil
	//
	// PasswordReset/Reset
	// err: "PasswordReset.Validate/AlreadyClosed"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/FailedToValidateTokenBecauseSessionClosed", info
	//
}

func ExamplePasswordReset_createSessionFailedBecauseLoginNotFound() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録ていないログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "unknown-login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: "User.GetUser/NotFound.User"
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/FailedToGetUserBecauseUserNotFound", info
	//
}

func ExamplePasswordReset_createSessionFailedBecauseDestinationNotFound() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	//h.registerResetDestination("user-id")                                 // 宛先を登録しない

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録したログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: "PasswordReset.CreateSession/NotFound.Destination"
	// log: "User/GetUser/TryToGetUser", debug
	// log: "User/GetUser/GetUser", info
	// log: "PasswordReset/CreateSession/TryToCreateSession", debug
	// log: "PasswordReset/CreateSession/FailedToCreateSessionBecauseDestinationNotFound", info
	//
}

func ExamplePasswordReset_getStatus() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/GetStatus", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printResetStatus(passwordResetHandler.status)
	})

	// 疑似的にステータスを変更
	h.passwordReset.sessions.UpdateStatusToFailed(passwordResetHandler.session, h.now(), errors.New("send token error"))

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printResetStatus(passwordResetHandler.status)
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: nil
	// status: {sending: {since: "2020-01-01T00:00:00Z"}}
	//
	// PasswordReset/GetStatus
	// err: nil
	// status: {failed: {at: "2020-01-01T00:01:00Z", reason: "send token error"}}
	//
}

func ExamplePasswordReset_getStatusFailedBecauseSessionNotFound() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 存在しないセッションステータスの取得を試みる
	passwordResetHandler.session = password_reset.NewSession("unknown-session")

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: "PasswordReset.GetStatus/NotFound.Session"
	// log: "PasswordReset/GetStatus/TryToGetStatus", debug
	// log: "PasswordReset/GetStatus/FailedToGetStatusBecauseSessionNotFound", audit
	//
}

func ExamplePasswordReset_getStatusFailedBecauseUnknownLogin() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 存在しないログインID でステータスの取得を試みる
	passwordResetHandler.login = user.NewLogin("unknown-login-id")

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: "PasswordReset.GetStatus/MatchFailed.Login"
	// log: "PasswordReset/GetStatus/TryToGetStatus", debug
	// log: "PasswordReset/GetStatus/FailedToGetStatusBecauseLoginMatchFailed", audit
	//
}

func ExamplePasswordReset_getStatusFailedBecauseDifferentLogin() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"})                         // ユーザーを登録
	h.registerUserData("another-user-id", "another-login-id", "another-password", []string{"role"}) // 別なユーザーを登録
	h.registerResetDestination("user-id")                                                           // 宛先を登録
	h.registerResetDestination("another-user-id")                                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")
	// 別なユーザーでも並行してリセット
	another_passwordResetHandler := newPasswordResetHandler(handler, "another-login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 並行してリセット
	h.passwordReset.gen.another()
	h.newRequest("PasswordReset/CreateSession", time.Minute(0), another_passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(another_passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 別なログインID のステータス取得を試みる
	passwordResetHandler.login = user.NewLogin("another-login-id")

	h.newRequest("PasswordReset/GetStatus", time.Minute(2), passwordResetHandler, func() {
		NewPasswordReset(client).GetStatus(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/GetStatus
	// err: "PasswordReset.GetStatus/MatchFailed.Login"
	// log: "PasswordReset/GetStatus/TryToGetStatus", debug
	// log: "PasswordReset/GetStatus/FailedToGetStatusBecauseLoginMatchFailed", audit
	//
}

func ExamplePasswordReset_resetFailedBecauseSessionNotFound() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 違うトークンでリセットを試みる
	passwordResetHandler.token = "unknown-token"

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: "PasswordReset.Validate/NotFound.Session"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/FailedToValidateTokenBecauseSessionNotFound", audit
	//
}

func ExamplePasswordReset_resetFailedBecauseSessionExpired() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	passwordResetHandler.fetchToken()

	// 有効期限切れ
	h.newRequest("PasswordReset/Reset", time.Minute(31), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: "PasswordReset.Validate/AlreadyExpired"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/FailedToValidateTokenBecauseSessionExpired", info
	//
}

func ExamplePasswordReset_resetSuccessWithTimeLimit() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	passwordResetHandler.fetchToken()

	// 有効期限ぎりぎり
	h.newRequest("PasswordReset/Reset", time.Minute(30), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: nil
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/AuthByToken", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/Change", audit
	// log: "PasswordReset/CloseSession/TryToCloseSession", debug
	// log: "PasswordReset/CloseSession/CloseSession", info
	// log: "Ticket/Register/TryToRegister", debug
	// log: "Ticket/Register/Register", info
	// log: "Credential/IssueTicket/TryToIssueTicket", debug
	// log: "Credential/IssueTicket/IssueTicket", info
	// log: "Credential/IssueApiToken/TryToIssue", debug
	// log: "Credential/IssueApiToken/Issue", info
	// log: "Credential/IssueContentToken/TryToIssue", debug
	// log: "Credential/IssueContentToken/Issue", info
	//
}

func ExamplePasswordReset_resetFailedBecauseUnknownLogin() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	passwordResetHandler.fetchToken()

	// 存在しないログインID でリセットを試みる
	passwordResetHandler.login = user.NewLogin("unknown-login-id")

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: "PasswordReset.Validate/MatchFailed.Login"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/FailedToValidateTokenBecauseLoginMatchFailed", audit
	//
}

func ExamplePasswordReset_resetFailedBecauseDifferentLogin() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"})                         // ユーザーを登録
	h.registerUserData("another-user-id", "another-login-id", "another-password", []string{"role"}) // 別なユーザーを登録
	h.registerResetDestination("user-id")                                                           // 宛先を登録
	h.registerResetDestination("another-user-id")                                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID でリセット
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "new-password")
	// 別なユーザーでも並行してリセット
	another_passwordResetHandler := newPasswordResetHandler(handler, "another-login-id", "new-password")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 並行してリセット
	h.passwordReset.gen.another()
	h.newRequest("PasswordReset/CreateSession", time.Minute(0), another_passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(another_passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	passwordResetHandler.fetchToken()

	// 別なログインID でリセットを試みる
	passwordResetHandler.login = user.NewLogin("another-login-id")

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: "PasswordReset.Validate/MatchFailed.Login"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/FailedToValidateTokenBecauseLoginMatchFailed", audit
	//
}

func ExamplePasswordReset_resetFailedBecauseEmptyPassword() {
	h := newPasswordResetTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録
	h.registerResetDestination("user-id")                                   // 宛先を登録

	client := NewClient(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	// 登録されたログインID で空のパスワードにリセットを試みる
	passwordResetHandler := newPasswordResetHandler(handler, "login-id", "")

	h.newRequest("PasswordReset/CreateSession", time.Minute(0), passwordResetHandler, func() {
		NewPasswordReset(client).CreateSession(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordReset/SendToken", time.Minute(1), passwordResetHandler, func() {
		NewPasswordReset(client).SendToken(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	passwordResetHandler.fetchToken()

	h.newRequest("PasswordReset/Reset", time.Minute(3), passwordResetHandler, func() {
		NewPasswordReset(client).Reset(passwordResetHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordReset/CreateSession
	// err: nil
	//
	// PasswordReset/SendToken
	// err: nil
	//
	// PasswordReset/Reset
	// err: "Password.Check/Length.Empty"
	// log: "PasswordReset/Validate/TryToValidateToken", debug
	// log: "PasswordReset/Validate/AuthByToken", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/FailedToChangeBecausePasswordCheckFailed", info
	//
}

type (
	passwordResetTestHelper struct {
		*testInfra
	}

	passwordResetTestHandler struct {
		*commonTestHandler

		login       user.Login
		newPassword password.RawPassword

		session password_reset.Session
		dest    password_reset.Destination
		status  password_reset.Status
		token   password_reset.Token
	}
)

func newPasswordResetTestHelper() passwordResetTestHelper {
	return passwordResetTestHelper{
		testInfra: newTestInfra(),
	}
}

func (h passwordResetTestHelper) registerResetDestination(userID user.UserID) {
	testUser := user.NewUser(userID)

	err := h.passwordReset.destinations.RegisterDestination(testUser, password_reset.NewLogDestination())
	if err != nil {
		log.Fatalf("register destination error: %s", err)
	}
}

func newPasswordResetHandler(handler *commonTestHandler, loginID user.LoginID, newPassword password.RawPassword) *passwordResetTestHandler {
	return &passwordResetTestHandler{
		commonTestHandler: handler,

		login:       user.NewLogin(loginID),
		newPassword: newPassword,
	}
}

func (handler *passwordResetTestHandler) handler() PasswordResetHandler {
	return handler
}
func (handler *passwordResetTestHandler) CreateSessionRequest() (request.Request, user.Login, error) {
	return handler.newRequest(), handler.login, nil
}
func (handler *passwordResetTestHandler) CreateSessionResponse(session password_reset.Session, err error) {
	handler.setError(err)
	handler.session = session
}
func (handler *passwordResetTestHandler) SendTokenResponse(err error) {
	handler.setError(err)
}
func (handler *passwordResetTestHandler) GetStatusRequest() (request.Request, user.Login, password_reset.Session, error) {
	return handler.newRequest(), handler.login, handler.session, nil
}
func (handler *passwordResetTestHandler) GetStatusResponse(dest password_reset.Destination, status password_reset.Status, err error) {
	handler.setError(err)
	handler.dest = dest
	handler.status = status
}
func (handler *passwordResetTestHandler) ResetRequest() (request.Request, user.Login, password_reset.Token, password.RawPassword, error) {
	return handler.newRequest(), handler.login, handler.token, handler.newPassword, nil
}
func (handler *passwordResetTestHandler) ResetResponse(err error) {
	handler.setError(err)
}

func (handler *passwordResetTestHandler) fetchToken() {
	message, found := handler.message.fetch()
	if found {
		tips := strings.Split(message, ": ")
		if len(tips) > 1 {
			handler.token = password_reset.Token(tips[1])
		}
	}

	handler.message.clear()
}
