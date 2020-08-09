package client

import (
	"encoding/json"
	"errors"
	"fmt"
	golog "log"
	gotime "time"

	"github.com/getto-systems/project-example-id/log"

	"github.com/getto-systems/project-example-id/credential/log"
	"github.com/getto-systems/project-example-id/credential/repository/api_user"
	password_log "github.com/getto-systems/project-example-id/password/log"
	password_repository_password "github.com/getto-systems/project-example-id/password/repository/password"
	password_reset_job_queue_send_token "github.com/getto-systems/project-example-id/password_reset/job_queue/send_token"
	password_reset_log "github.com/getto-systems/project-example-id/password_reset/log"
	password_reset_repository_destination "github.com/getto-systems/project-example-id/password_reset/repository/destination"
	password_reset_repository_session "github.com/getto-systems/project-example-id/password_reset/repository/session"
	password_reset_sender "github.com/getto-systems/project-example-id/password_reset/sender"
	ticket_log "github.com/getto-systems/project-example-id/ticket/log"
	ticket_repository_ticket "github.com/getto-systems/project-example-id/ticket/repository/ticket"
	user_log "github.com/getto-systems/project-example-id/user/log"
	user_repository_user "github.com/getto-systems/project-example-id/user/repository/user"

	credential_infra "github.com/getto-systems/project-example-id/credential/infra"
	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"
	ticket_infra "github.com/getto-systems/project-example-id/infra/ticket"
	user_infra "github.com/getto-systems/project-example-id/infra/user"
	password_infra "github.com/getto-systems/project-example-id/password/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/data/request"
	"github.com/getto-systems/project-example-id/data/ticket"
	"github.com/getto-systems/project-example-id/data/time"
	"github.com/getto-systems/project-example-id/data/user"
	"github.com/getto-systems/project-example-id/password"

	"github.com/getto-systems/project-example-id/credential/core"
	"github.com/getto-systems/project-example-id/password/core"
	password_reset_core "github.com/getto-systems/project-example-id/password_reset"
	ticket_core "github.com/getto-systems/project-example-id/ticket"
	user_core "github.com/getto-systems/project-example-id/user"
)

type (
	testInfra struct {
		logger  *testLogger
		message *testMessage

		exp testExpiration

		session testSession

		ticket        ticketTestInfra
		credential    credentialTestInfra
		user          userTestInfra
		password      passwordTestInfra
		passwordReset passwordResetTestInfra
	}

	testExpiration struct {
		password      ticket.Expiration
		passwordReset password_reset.Expiration
	}

	testSession struct {
		nowSecond  time.Second
		credential *credential.Credential
		nonce      *credential.TicketNonce
	}

	ticketTestInfra struct {
		gen ticket_infra.TicketNonceGenerator

		tickets ticket_infra.TicketRepository
	}
	credentialTestInfra struct {
		ticketSign credential_infra.TicketSign
		apiUsers   credential_infra.ApiUserRepository
	}
	userTestInfra struct {
		users user_infra.UserRepository
	}
	passwordTestInfra struct {
		enc password_infra.PasswordEncrypter

		passwords password_infra.PasswordRepository
	}
	passwordResetTestInfra struct {
		gen *passwordResetTestSessionGenerator

		sessions     password_reset_infra.SessionRepository
		destinations password_reset_infra.DestinationRepository

		sendTokenQueue password_reset_infra.SendTokenJobQueue
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

	credentialTestApiSigner     struct{}
	credentialTestContentSigner struct{}

	passwordTestEncrypter struct{}

	passwordResetTestSessionGenerator struct {
		id    password_reset.SessionID
		token password_reset.Token
	}
)

type (
	testHandler interface {
		testContext() testContext
	}

	commonTestHandler struct {
		*testInfra
		context testContext
	}

	testFormatter struct {
		context    testContext
		credential *credential.Credential

		logger  *testLogger
		message *testMessage
	}

	testContext struct {
		request request.Request
		err     error
	}
)

func newTestInfra() *testInfra {
	return &testInfra{
		logger:  newTestLogger(),
		message: newTestMessage(),

		session: testSession{
			nowSecond: time.Second(0),
		},

		exp: testExpiration{
			password: ticket.NewExpiration(ticket.ExpirationParam{
				Expires:     time.Minute(5),
				ExtendLimit: time.Minute(8),
			}),
			passwordReset: password_reset.NewExpiration(time.Minute(30)),
		},

		ticket: ticketTestInfra{
			gen: ticketTestNonceGenerator{},

			tickets: ticket_repository_ticket.NewMemoryStore(),
		},
		credential: credentialTestInfra{
			ticketSign: ticketTestSign{},
			apiUsers:   credential_repository_apiUser.NewMemoryStore(),
		},
		user: userTestInfra{
			users: user_repository_user.NewMemoryStore(),
		},
		password: passwordTestInfra{
			enc: passwordTestEncrypter{},

			passwords: password_repository_password.NewMemoryStore(),
		},
		passwordReset: passwordResetTestInfra{
			gen: &passwordResetTestSessionGenerator{
				id:    "reset-session-id",
				token: "reset-token",
			},

			sessions:     password_reset_repository_session.NewMemoryStore(),
			destinations: password_reset_repository_destination.NewMemoryStore(),

			sendTokenQueue: password_reset_job_queue_send_token.NewMemoryQueue(),
		},
	}
}

func (backend *testInfra) newBackend() Backend {
	return NewBackend(
		ticket_core.NewAction(
			ticket_log.NewLogger(backend.logger),

			backend.ticket.gen,

			backend.ticket.tickets,
		),
		credential_core.NewAction(
			credential_log.NewLogger(backend.logger),

			backend.credential.ticketSign,
			credentialTestApiSigner{},
			credentialTestContentSigner{},

			backend.credential.apiUsers,
		),
		user_core.NewAction(
			user_log.NewLogger(backend.logger),

			backend.user.users,
		),
		password_core.NewAction(
			password_log.NewLogger(backend.logger),

			backend.exp.password,
			backend.password.enc,

			backend.password.passwords,
		),
		password_reset_core.NewAction(
			password_reset_log.NewLogger(backend.logger),

			backend.exp.password,
			backend.exp.passwordReset,
			backend.passwordReset.gen,

			backend.passwordReset.sessions,
			backend.passwordReset.destinations,

			backend.passwordReset.sendTokenQueue,
			password_reset_sender.NewTokenSender(backend.message),
		),
	)
}

func (backend *testInfra) now() time.RequestedAt {
	now, err := gotime.Parse(gotime.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		golog.Fatalf("failed to initialize 'now': %s", err)
	}
	return time.RequestedAt(now.Add(gotime.Duration(backend.session.nowSecond * 1_000_000_000)))
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

func (ticketTestSign) Sign(user user.User, nonce credential.TicketNonce, expires time.Expires) (_ credential.TicketSignature, err error) {
	data, err := json.Marshal(ticketTestSignToken{
		UserID:  string(user.ID()),
		Nonce:   string(nonce),
		Expires: gotime.Time(expires).Unix(),
	})
	if err != nil {
		return
	}

	return credential.TicketSignature(data), nil
}
func (ticketTestSign) Parse(signature credential.TicketSignature) (_ user.User, _ credential.TicketNonce, err error) {
	var data ticketTestSignToken

	err = json.Unmarshal(signature, &data)
	if err != nil {
		return
	}

	return user.NewUser(user.UserID(data.UserID)), credential.TicketNonce(data.Nonce), nil
}

func (ticketTestNonceGenerator) GenerateTicketNonce() (_ credential.TicketNonce, err error) {
	return "ticket-nonce", nil
}

func (credentialTestApiSigner) Sign(user user.User, roles credential.ApiRoles, expires time.Expires) (_ credential.ApiToken, err error) {
	return credential.NewApiToken(roles, []byte("api-token")), nil
}

func (credentialTestContentSigner) Sign(expires time.Expires) (_ credential.ContentToken, err error) {
	return credential.NewContentToken(
		credential.ContentKeyID("content-key"),
		credential.ContentPolicy([]byte("content-policy")),
		credential.ContentSignature([]byte("content-signature")),
	), nil
}

func (passwordTestEncrypter) GeneratePassword(raw password.RawPassword) (password.HashedPassword, error) {
	return password.HashedPassword(raw), nil
}
func (passwordTestEncrypter) MatchPassword(hashed password.HashedPassword, raw password.RawPassword) (bool, error) {
	return string(hashed) == string(raw), nil
}

func (gen *passwordResetTestSessionGenerator) GenerateSession() (password_reset.SessionID, password_reset.Token, error) {
	return gen.id, gen.token, nil
}
func (gen *passwordResetTestSessionGenerator) another() {
	gen.id = "another-reset-session-id"
	gen.token = "another-reset-token"
}

func (backend *testInfra) credentialHandler() CredentialHandler {
	return backend
}
func (backend *testInfra) GetTicket() (_ credential.Ticket, err error) {
	if backend.session.credential == nil {
		err = errors.New("credential not set")
		return
	}

	if backend.session.nonce == nil {
		return backend.session.credential.Ticket(), nil
	}

	return credential.NewTicket(backend.session.credential.Ticket().Signature(), *backend.session.nonce), nil
}
func (backend *testInfra) SetCredential(credential credential.Credential) {
	backend.session.credential = &credential
}
func (backend *testInfra) ClearCredential() {
	backend.session.credential = nil
}

func (backend *testInfra) setNonce(nonce credential.TicketNonce) {
	backend.session.nonce = &nonce
}
func (backend *testInfra) setCredentialNonce(nonce credential.TicketNonce) {
	if backend.session.credential != nil {
		user, _, _ := backend.credential.ticketSign.Parse(backend.session.credential.Ticket().Signature())
		signature, _ := backend.credential.ticketSign.Sign(user, nonce, backend.session.credential.Expires())

		credential := credential.NewCredential(
			credential.NewTicket(signature, nonce),
			backend.session.credential.ApiToken(),
			backend.session.credential.ContentToken(),
			backend.session.credential.Expires(),
		)
		backend.session.credential = &credential
	}
}
func (backend *testInfra) setCredentialUser(user user.User) {
	if backend.session.credential != nil {
		_, nonce, _ := backend.credential.ticketSign.Parse(backend.session.credential.Ticket().Signature())
		signature, _ := backend.credential.ticketSign.Sign(user, nonce, backend.session.credential.Expires())

		credential := credential.NewCredential(
			credential.NewTicket(signature, nonce),
			backend.session.credential.ApiToken(),
			backend.session.credential.ContentToken(),
			backend.session.credential.Expires(),
		)
		backend.session.credential = &credential
	}
}

func (backend *testInfra) registerUserData(userID user.UserID, loginID user.LoginID, rawPassword password.RawPassword, apiRoles credential.ApiRoles) {
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

	err = backend.credential.apiUsers.RegisterApiRoles(testUser, apiRoles)
	if err != nil {
		golog.Fatalf("register api roles error: %s", err)
	}
}
func (backend *testInfra) registerUserDataWithoutApiRoles(userID user.UserID, loginID user.LoginID, rawPassword password.RawPassword) {
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
}
func (backend *testInfra) registerOnlyUserAndLogin(userID user.UserID, loginID user.LoginID) {
	testUser := user.NewUser(userID)

	err := backend.user.users.RegisterUser(testUser, user.NewLogin(loginID))
	if err != nil {
		golog.Fatalf("register user error: %s", err)
	}
}

func (backend *testInfra) newRequest(label string, nowSecond time.Second, handler testHandler, exec func(), format func(f testFormatter)) {
	backend.session.nowSecond = nowSecond

	fmt.Println(label)
	exec()
	format(testFormatter{
		context:    handler.testContext(),
		credential: backend.session.credential,
		logger:     backend.logger,
		message:    backend.message,
	})
	fmt.Println("")
	backend.logger.clear()
}

func (backend *testInfra) newHandler() *commonTestHandler {
	return &commonTestHandler{
		testInfra: backend,
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
	fmt.Printf("request: \"%s\"\n", formatTime(f.context.request.RequestedAt().Time()))
}
func (f testFormatter) printCredential() {
	if f.credential == nil {
		fmt.Println("credential: nil")
	} else {
		fmt.Printf(
			"credential: expires: \"%s\", roles: %s\n",
			formatTime(f.credential.Expires().Time()),
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
func (f testFormatter) printResetDestination(dest password_reset.Destination) {
	fmt.Printf("dest: %v\n", dest) // TODO ちゃんとする
}
func (f testFormatter) printResetStatus(status password_reset.Status) {
	if status.Waiting() {
		fmt.Printf("status: {waiting: {since: \"%s\"}}\n", formatTime(status.WaitingSince()))
	} else if status.Sending() {
		fmt.Printf("status: {sending: {since: \"%s\"}}\n", formatTime(status.SendingSince()))
	} else if status.Complete() {
		fmt.Printf("status: {complete: {at: \"%s\"}}\n", formatTime(status.CompleteAt()))
	} else if status.Failed() {
		at, reason := status.FailedAtAndReason()
		fmt.Printf("status: {failed: {at: \"%s\", reason: \"%s\"}}\n", formatTime(at), reason)
	} else {
		fmt.Println("status: {EMPTY}")
	}
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

func formatTime(t time.Time) string {
	return gotime.Time(t).Format(gotime.RFC3339)
}
