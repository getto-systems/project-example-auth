package _usecase

import (
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/request"
	"github.com/getto-systems/project-example-id/user"
)

func ExamplePasswordChange_getLogin_change() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printCredential()
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLogin(passwordChangeHandler.login)
	})

	h.newRequest("PasswordChange/Change", minute(2), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
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

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).GetLogin(passwordChangeHandler)
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
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "User/GetLogin/TryToGetLogin", debug
	// log: "User/GetLogin/GetLogin", info
	//
}

func ExamplePasswordChange_changeLog() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/GetLogin", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/Change", minute(2), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
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
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/AuthByPassword", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/Change", audit
	//
}

func ExamplePasswordChange_disableOldPassword() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	h.newRequest("PasswordChange/Change", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 前のパスワードでログインを試みる
	h.newRequest("PasswordLogin", minute(2), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/Change
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

func ExamplePasswordChange_getLoginFailedBecauseTicketExpired() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 有効期限切れでログイン情報の取得を試みる
	h.newRequest("PasswordChange/GetLogin", minute(10), passwordChangeHandler, func() {
		NewPasswordChange(u).GetLogin(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLogin(passwordChangeHandler.login)
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/GetLogin
	// request: "2020-01-01T00:10:00Z"
	// err: "Ticket.Validate/AlreadyExpired"
	// login: {}
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

func ExamplePasswordChange_changeFailedBecauseTicketExpired() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、新パスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 有効期限切れでパスワードの変更を試みる
	h.newRequest("PasswordChange/Change", minute(10), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/Change
	// request: "2020-01-01T00:10:00Z"
	// err: "Ticket.Validate/AlreadyExpired"
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/FailedToValidateBecauseExpired", info
	//
}

func ExamplePasswordChange_changeFailedBecauseDifferentPassword() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと異なるパスワードで確認、空のパスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "different-password", "new-password")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 登録されているものと違うパスワードで変更を試みる
	h.newRequest("PasswordChange/Change", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/Change
	// request: "2020-01-01T00:01:00Z"
	// err: "Password.Validate/MatchFailed"
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/FailedToValidateBecausePasswordMatchFailed", audit
	//
}

func ExamplePasswordChange_changeFailedBecauseEmptyPassword() {
	h := newPasswordChangeTestHelper()
	h.registerUserData("user-id", "login-id", "password", []string{"role"}) // ユーザーを登録

	u := NewUsecase(h.newBackend(), h.credentialHandler())

	handler := h.newHandler()

	passwordLoginHandler := newPasswordLoginHandler(handler, "login-id", "password")
	// 登録済みデータと同じパスワードで確認、空のパスワードに変更
	passwordChangeHandler := newPasswordChangeHandler(handler, "password", "")

	h.newRequest("PasswordLogin", minute(0), passwordLoginHandler, func() {
		NewPasswordLogin(u).Login(passwordLoginHandler)
	}, func(f testFormatter) {
		f.printError()
	})

	// 空のパスワードへの変更を試みる
	h.newRequest("PasswordChange/Change", minute(1), passwordChangeHandler, func() {
		NewPasswordChange(u).Change(passwordChangeHandler)
	}, func(f testFormatter) {
		f.printRequest()
		f.printError()
		f.printLog()
	})

	// Output:
	// PasswordLogin
	// err: nil
	//
	// PasswordChange/Change
	// request: "2020-01-01T00:01:00Z"
	// err: "Password.Check/Length.Empty"
	// log: "Credential/ParseTicket/TryToParseTicket", debug
	// log: "Credential/ParseTicket/ParseTicket", info
	// log: "Ticket/Validate/TryToValidate", debug
	// log: "Ticket/Validate/AuthByTicket", info
	// log: "Password/Validate/TryToValidate", debug
	// log: "Password/Validate/AuthByPassword", audit
	// log: "Password/Change/TryToChange", debug
	// log: "Password/Change/FailedToChangeBecausePasswordCheckFailed", info
	//
}

type (
	passwordChangeTestHelper struct {
		*testInfra
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
		testInfra: newTestInfra(),
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
