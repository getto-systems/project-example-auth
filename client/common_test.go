package client

import (
	"encoding/json"
	"errors"
	"fmt"
	golog "log"
	gotime "time"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"

	ticket_log "github.com/getto-systems/project-example-id/ticket/log"
	ticket_repository_ticket "github.com/getto-systems/project-example-id/ticket/repository/ticket"

	api_token_log "github.com/getto-systems/project-example-id/api_token/log"
	api_token_repository_api_user "github.com/getto-systems/project-example-id/api_token/repository/api_user"

	user_log "github.com/getto-systems/project-example-id/user/log"
	user_repository_user "github.com/getto-systems/project-example-id/user/repository/user"

	password_log "github.com/getto-systems/project-example-id/password/log"
	password_repository_password "github.com/getto-systems/project-example-id/password/repository/password"

	password_reset_job_queue_send_token "github.com/getto-systems/project-example-id/password_reset/job_queue/send_token"
	password_reset_log "github.com/getto-systems/project-example-id/password_reset/log"
	password_reset_repository_destination "github.com/getto-systems/project-example-id/password_reset/repository/destination"
	password_reset_repository_session "github.com/getto-systems/project-example-id/password_reset/repository/session"
	password_reset_sender "github.com/getto-systems/project-example-id/password_reset/sender"
)

type (
	testBackend struct {
		logger  *testLogger
		message *testMessage

		nowSecond  time.Second
		credential *data.Credential
		nonce      *ticket.Nonce

		ticket        ticketTestBackend
		apiToken      apiTokenTestBackend
		user          userTestBackend
		password      passwordTestBackend
		passwordReset passwordResetTestBackend
	}

	ticketTestBackend struct {
		sign ticket.TicketSign
		gen  ticket.NonceGenerator

		tickets ticket.TicketRepository
	}
	apiTokenTestBackend struct {
		apiUsers api_token.ApiUserRepository
	}
	userTestBackend struct {
		users user.UserRepository
	}
	passwordTestBackend struct {
		enc password.PasswordEncrypter

		passwords password.PasswordRepository
	}
	passwordResetTestBackend struct {
		sessions     password_reset.SessionRepository
		destinations password_reset.DestinationRepository

		sendTokenQueue password_reset.SendTokenJobQueue
	}

	testLogger struct {
		log []testLogEntry
	}
	testLogEntry struct {
		level string
		entry log.Entry
	}

	testMessage struct {
		log []string
	}

	ticketTestSign           struct{}
	ticketTestNonceGenerator struct{}

	apiTokenTestApiSigner     struct{}
	apiTokenTestContentSigner struct{}

	passwordTestEncrypter struct{}

	passwordResetTestSessionGenerator struct{}
)

type (
	testHandler interface {
		testContext() testContext
	}

	commonTestHandler struct {
		*testBackend
		context testContext
	}

	testFormatter struct {
		context    testContext
		credential *data.Credential

		logger  *testLogger
		message *testMessage
	}

	testContext struct {
		request request.Request
		err     error
	}
)

func newTestBackend() *testBackend {
	return &testBackend{
		logger:  newTestLogger(),
		message: newTestMessage(),

		nowSecond: time.Second(0),

		ticket: ticketTestBackend{
			sign: ticketTestSign{},
			gen:  ticketTestNonceGenerator{},

			tickets: ticket_repository_ticket.NewMemoryStore(),
		},
		apiToken: apiTokenTestBackend{
			apiUsers: api_token_repository_api_user.NewMemoryStore(),
		},
		user: userTestBackend{
			users: user_repository_user.NewMemoryStore(),
		},
		password: passwordTestBackend{
			enc: passwordTestEncrypter{},

			passwords: password_repository_password.NewMemoryStore(),
		},
		passwordReset: passwordResetTestBackend{
			sessions:     password_reset_repository_session.NewMemoryStore(),
			destinations: password_reset_repository_destination.NewMemoryStore(),

			sendTokenQueue: password_reset_job_queue_send_token.NewMemoryQueue(),
		},
	}
}

func (backend *testBackend) newBackend() Backend {
	return NewBackend(
		NewTicketAction(
			ticket_log.NewLogger(backend.logger),

			backend.ticket.sign,
			ticket.ExpirationParam{
				Expires:     time.Minute(5),
				ExtendLimit: time.Minute(8),
			},
			backend.ticket.gen,

			backend.ticket.tickets,
		),
		NewApiTokenAction(
			api_token_log.NewLogger(backend.logger),

			apiTokenTestApiSigner{},
			apiTokenTestContentSigner{},

			backend.apiToken.apiUsers,
		),
		NewUserAction(
			user_log.NewLogger(backend.logger),

			backend.user.users,
		),
		NewPasswordAction(
			password_log.NewLogger(backend.logger),

			backend.password.enc,

			backend.password.passwords,
		),
		NewPasswordResetAction(
			password_reset_log.NewLogger(backend.logger),

			time.Minute(30),
			passwordResetTestSessionGenerator{},

			backend.passwordReset.sessions,
			backend.passwordReset.destinations,

			backend.passwordReset.sendTokenQueue,
			password_reset_sender.NewTokenSender(backend.message),
		),
	)
}

func (backend *testBackend) now() time.RequestedAt {
	now, err := gotime.Parse(gotime.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		golog.Fatalf("failed to initialize 'now': %s", err)
	}
	return time.RequestedAt(now.Add(gotime.Duration(backend.nowSecond * 1_000_000_000)))
}

func newTestLogger() *testLogger {
	return &testLogger{}
}
func (logger *testLogger) Audit(entry log.Entry) {
	logger.log = append(logger.log, logger.entry("audit", entry))
}
func (logger *testLogger) Error(entry log.Entry) {
	logger.log = append(logger.log, logger.entry("error", entry))
}
func (logger *testLogger) Info(entry log.Entry) {
	logger.log = append(logger.log, logger.entry("info", entry))
}
func (logger *testLogger) Debug(entry log.Entry) {
	logger.log = append(logger.log, logger.entry("debug", entry))
}
func (logger *testLogger) entry(level string, entry log.Entry) testLogEntry {
	return testLogEntry{level: level, entry: entry}
}
func (logger *testLogger) clear() {
	logger.log = []testLogEntry{}
}

func newTestMessage() *testMessage {
	return &testMessage{}
}
func (log *testMessage) Send(message string) error {
	log.log = append(log.log, message)
	return nil
}
func (log *testMessage) fetch() (_ string, found bool) {
	if len(log.log) == 0 {
		return
	}

	return log.log[0], true
}
func (log *testMessage) clear() {
	log.log = []string{}
}

type ticketTestSignToken struct {
	UserID  string `json:"user_id"`
	Nonce   string `json:"nonce"`
	Expires int64  `json:"expires"`
}

func (ticketTestSign) Sign(user user.User, nonce ticket.Nonce, expires time.Expires) (_ ticket.Token, err error) {
	data, err := json.Marshal(ticketTestSignToken{
		UserID:  string(user.ID()),
		Nonce:   string(nonce),
		Expires: gotime.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return ticket.Token(data), nil
}
func (ticketTestSign) Parse(token ticket.Token) (_ user.User, _ ticket.Nonce, err error) {
	var data ticketTestSignToken

	err = json.Unmarshal(token, &data)
	if err != nil {
		return
	}

	return user.NewUser(user.UserID(data.UserID)), ticket.Nonce(data.Nonce), nil
}

func (ticketTestNonceGenerator) GenerateNonce() (_ ticket.Nonce, err error) {
	return "ticket-nonce", nil
}

func (apiTokenTestApiSigner) Sign(user user.User, roles api_token.ApiRoles, expires time.Expires) (_ api_token.ApiToken, err error) {
	return api_token.NewApiToken(roles, []byte("api-token")), nil
}

func (apiTokenTestContentSigner) Sign(expires time.Expires) (_ api_token.ContentToken, err error) {
	return api_token.NewContentToken(
		api_token.ContentKeyID("content-key"),
		api_token.ContentPolicy([]byte("content-policy")),
		api_token.ContentSignature([]byte("content-signature")),
	), nil
}

func (passwordTestEncrypter) GeneratePassword(raw password.RawPassword) (password.HashedPassword, error) {
	return password.HashedPassword(raw), nil
}
func (passwordTestEncrypter) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) (bool, error) {
	return string(hashed) == string(raw), nil
}

func (passwordResetTestSessionGenerator) GenerateSession() (password_reset.SessionID, password_reset.Token, error) {
	return "reset-session-id", "reset-token", nil
}

func (backend *testBackend) credentialHandler() CredentialHandler {
	return backend
}
func (backend *testBackend) GetTicket() (_ ticket.Ticket, err error) {
	if backend.credential == nil {
		err = errors.New("credential not set")
		return
	}

	if backend.nonce == nil {
		return backend.credential.Ticket(), nil
	}

	return ticket.NewTicket(backend.credential.Ticket().Token(), *backend.nonce), nil
}
func (backend *testBackend) SetCredential(credential data.Credential) {
	backend.credential = &credential
}
func (backend *testBackend) ClearCredential() {
	backend.credential = nil
}

func (backend *testBackend) setNonce(nonce ticket.Nonce) {
	backend.nonce = &nonce
}
func (backend *testBackend) setCredentialNonce(nonce ticket.Nonce) {
	if backend.credential != nil {
		user, _, _ := backend.ticket.sign.Parse(backend.credential.Ticket().Token())
		token, _ := backend.ticket.sign.Sign(user, nonce, backend.credential.Expires())

		credential := data.NewCredential(
			ticket.NewTicket(token, nonce),
			backend.credential.ApiToken(),
			backend.credential.ContentToken(),
			backend.credential.Expires(),
		)
		backend.credential = &credential
	}
}
func (backend *testBackend) setCredentialUser(user user.User) {
	if backend.credential != nil {
		_, nonce, _ := backend.ticket.sign.Parse(backend.credential.Ticket().Token())
		token, _ := backend.ticket.sign.Sign(user, nonce, backend.credential.Expires())

		credential := data.NewCredential(
			ticket.NewTicket(token, nonce),
			backend.credential.ApiToken(),
			backend.credential.ContentToken(),
			backend.credential.Expires(),
		)
		backend.credential = &credential
	}
}

func (backend *testBackend) registerUserData(userID user.UserID, loginID user.LoginID, rawPassword password.RawPassword, apiRoles api_token.ApiRoles) {
	testUser := user.NewUser(userID)

	err := backend.user.users.RegisterUser(testUser, user.NewLogin(loginID))
	if err != nil {
		golog.Fatalf("register user error: %s", err)
	}

	hashed, err := backend.password.enc.GeneratePassword(rawPassword)
	if err != nil {
		golog.Fatalf("generate password error: %s", err)
	}

	err = backend.password.passwords.ChangePassword(testUser, hashed)
	if err != nil {
		golog.Fatalf("change password error: %s", err)
	}

	err = backend.apiToken.apiUsers.RegisterApiRoles(testUser, apiRoles)
	if err != nil {
		golog.Fatalf("register api roles error: %s", err)
	}
}

func (backend *testBackend) newRequest(label string, nowSecond time.Second, handler testHandler, exec func(), format func(f testFormatter)) {
	backend.nowSecond = nowSecond

	fmt.Println(label)
	exec()
	format(testFormatter{
		context:    handler.testContext(),
		credential: backend.credential,
		logger:     backend.logger,
		message:    backend.message,
	})
	fmt.Println("")
	backend.logger.clear()
}

func (backend *testBackend) newHandler() *commonTestHandler {
	return &commonTestHandler{
		testBackend: backend,
	}
}
func (handler *commonTestHandler) handler() testHandler {
	return handler
}
func (handler *commonTestHandler) testContext() testContext {
	return handler.context
}
func (handler *commonTestHandler) newRequest() request.Request {
	req := request.NewRequest(handler.now(), "test-remote")
	handler.context.request = req
	return req
}
func (handler *commonTestHandler) setError(err error) {
	handler.context.err = err
}

func (f testFormatter) printError() {
	if f.context.err == nil {
		fmt.Println("err: nil")
	} else {
		fmt.Printf("err: \"%s\"\n", f.context.err)
	}
}
func (f testFormatter) printRequest() {
	fmt.Printf(
		"request: \"%s\"\n",
		gotime.Time(f.context.request.RequestedAt()).Format(gotime.RFC3339),
	)
}
func (f testFormatter) printCredential() {
	if f.credential == nil {
		fmt.Println("credential: nil")
	} else {
		fmt.Printf(
			"credential: expires: \"%s\", roles: %s\n",
			gotime.Time(f.credential.Expires()).Format(gotime.RFC3339),
			f.credential.ApiToken().ApiRoles(),
		)
	}
}
func (f testFormatter) printLogin(login user.Login) {
	fmt.Printf("login: {%s}\n", login.ID())
}
func (f testFormatter) printResetSession(session password_reset.Session) {
	fmt.Printf("session: {%s}\n", session.ID())
}
func (f testFormatter) printResetStatus(status password_reset.Status) {
	fmt.Printf("status: %v\n", status) // TODO ちゃんとする
}
func (f testFormatter) printResetToken(token password_reset.Token) {
	fmt.Printf("token: \"%s\"\n", token)
}

func (f testFormatter) printLog() {
	for _, entry := range f.logger.log {
		fmt.Printf("log: \"%s\", %s\n", entry.entry.Message, entry.level)
	}
}
func (f testFormatter) printMessage() {
	for _, message := range f.message.log {
		fmt.Printf("message: %s\n", message)
	}
}
