package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/getto-systems/project-example-id"

	"github.com/getto-systems/project-example-id/x_http_server/http_handler"
	"github.com/getto-systems/project-example-id/x_http_server/logger"
	"github.com/getto-systems/project-example-id/x_http_server/message"
	"github.com/getto-systems/project-example-id/x_http_server/nonce_generator"
	"github.com/getto-systems/project-example-id/x_http_server/password_encrypter"
	"github.com/getto-systems/project-example-id/x_http_server/reset_session_generator"
	"github.com/getto-systems/project-example-id/x_http_server/signer"

	"github.com/getto-systems/project-example-id/credential/log"
	"github.com/getto-systems/project-example-id/credential/repository/api_user"
	"github.com/getto-systems/project-example-id/password/log"
	"github.com/getto-systems/project-example-id/password/repository/password"
	"github.com/getto-systems/project-example-id/password_reset/job_queue/send_token"
	"github.com/getto-systems/project-example-id/password_reset/log"
	"github.com/getto-systems/project-example-id/password_reset/repository/destination"
	"github.com/getto-systems/project-example-id/password_reset/repository/session"
	"github.com/getto-systems/project-example-id/password_reset/sender"
	"github.com/getto-systems/project-example-id/ticket/log"
	"github.com/getto-systems/project-example-id/ticket/repository/ticket"
	"github.com/getto-systems/project-example-id/user/log"
	"github.com/getto-systems/project-example-id/user/repository/user"

	credential_infra "github.com/getto-systems/project-example-id/credential/infra"
	password_infra "github.com/getto-systems/project-example-id/password/infra"
	password_reset_infra "github.com/getto-systems/project-example-id/password_reset/infra"
	user_infra "github.com/getto-systems/project-example-id/user/infra"

	"github.com/getto-systems/project-example-id/credential"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/password_reset"
	"github.com/getto-systems/project-example-id/ticket"
	"github.com/getto-systems/project-example-id/user"

	"github.com/getto-systems/project-example-id/credential/core"
	"github.com/getto-systems/project-example-id/password/core"
	"github.com/getto-systems/project-example-id/password_reset/core"
	"github.com/getto-systems/project-example-id/ticket/core"
	"github.com/getto-systems/project-example-id/user/core"
)

const (
	HEADER_HANDLER = "X-Getto-Example-ID-Handler"
)

type (
	server struct {
		port string

		cookieDomain http_handler.CookieDomain
		backend      _usecase.Backend
	}

	infra struct {
		logger logger.LeveledLogger

		extend extendSecond
	}
	extendSecond struct {
		password credential.TicketExtendSecond
	}
)

func main() {
	log.Fatal(newServer().listen())
}
func (server server) listen() error {
	return http.ListenAndServe(
		server.port,
		server.mux(),
	)
}
func (server server) mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handle)
	return mux
}

func (server server) handle(w http.ResponseWriter, r *http.Request) {
	h := http_handler.NewHandler(w, r)
	u := _usecase.NewUsecase(server.backend, http_handler.NewCredentialHandler(server.cookieDomain, w, r))

	switch r.Header.Get(HEADER_HANDLER) {
	case "Renew":
		_usecase.NewRenew(u).Renew(http_handler.NewRenew(h))
	case "Logout":
		_usecase.NewLogout(u).Logout(http_handler.NewLogout(h))

	case "PasswordLogin":
		_usecase.NewPasswordLogin(u).Login(http_handler.NewPasswordLogin(h))

	case "PasswordChange/GetLogin":
		_usecase.NewPasswordChange(u).GetLogin(http_handler.NewPasswordChange(h))
	case "PasswordChange/Change":
		_usecase.NewPasswordChange(u).Change(http_handler.NewPasswordChange(h))

	case "PasswordReset/CreateSession":
		_usecase.NewPasswordReset(u).CreateSession(http_handler.NewPasswordReset(h))
	case "PasswordReset/SendToken":
		_usecase.NewPasswordReset(u).SendToken(http_handler.NewPasswordReset(h))
	case "PasswordReset/GetStatus":
		_usecase.NewPasswordReset(u).GetStatus(http_handler.NewPasswordReset(h))
	case "PasswordReset/Reset":
		_usecase.NewPasswordReset(u).Reset(http_handler.NewPasswordReset(h))

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "\"OK\"")
	}
}

func newServer() server {
	return server{
		port: ":8080",

		cookieDomain: http_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN")),
		backend:      newBackend(),
	}
}

func newBackend() _usecase.Backend {
	infra := newInfra()

	return _usecase.NewBackend(
		infra.newTicketAction(),
		infra.newCredentialAction(),
		infra.newUserAction(),
		infra.newPasswordAction(),
		infra.newPasswordResetAction(),
	)
}

func newInfra() infra {
	return infra{
		logger: newLeveledLogger(),
		extend: newExtend(),
	}
}
func (infra infra) newTicketAction() ticket.Action {
	return ticket_core.NewAction(
		ticket_log.NewLogger(infra.logger),

		credential.TicketExpireWeek(1),
		credential.TokenExpireMinute(5),

		nonce_generator.NewNonceGenerator(),

		ticket_repository_ticket.NewMemoryStore(),
	)
}
func (infra infra) newCredentialAction() credential.Action {
	apiUsers := credential_repository_apiUser.NewMemoryStore()

	initApiUserRepository(apiUsers)

	return credential_core.NewAction(
		credential_log.NewLogger(infra.logger),

		newTicketSigner(),
		newApiTokenSigner(),
		newContentTokenSigner(),

		apiUsers,
	)
}
func (infra infra) newUserAction() user.Action {
	users := user_repository_user.NewMemoryStore()

	initUserRepository(users)

	return user_core.NewAction(
		user_log.NewLogger(infra.logger),

		users,
	)
}
func (infra infra) newPasswordAction() password.Action {
	encrypter := password_encrypter.NewEncrypter(10) // bcrypt.DefaultCost
	passwords := password_repository_password.NewMemoryStore()

	initPasswordRepository(passwords, encrypter)

	return password_core.NewAction(
		password_log.NewLogger(infra.logger),

		infra.extend.password,
		encrypter,

		passwords,
	)
}
func (infra infra) newPasswordResetAction() password_reset.Action {
	destinations := password_reset_repository_destination.NewMemoryStore()

	initPasswordResetDestinationRepository(destinations)

	return password_reset_core.NewAction(
		password_reset_log.NewLogger(infra.logger),

		// パスワードリセットはパスワード認証と同等なので、最大延長期間はパスワード認証と同じ
		infra.extend.password,
		password_reset.ExpireMinute(30),

		reset_session_generator.NewGenerator(),

		password_reset_repository_session.NewMemoryStore(),
		destinations,

		password_reset_job_queue_sendToken.NewMemoryQueue(),
		password_reset_sender.NewTokenSender(message.NewLogMessage()),
	)
}

func initUserRepository(users user_infra.UserRepository) {
	login := user.NewLogin(user.LoginID(os.Getenv("ADMIN_LOGIN_ID")))

	err := users.RegisterUser(adminUser(), login)
	if err != nil {
		log.Fatalf("failed to register admin user: %s", err)
	}
}
func initApiUserRepository(apiUsers credential_infra.ApiUserRepository) {
	err := apiUsers.RegisterApiRoles(adminUser(), credential.ApiRoles([]string{"admin"}))
	if err != nil {
		log.Fatalf("failed to register admin user api roles: %s", err)
	}
}
func initPasswordRepository(passwords password_infra.PasswordRepository, gen password_infra.PasswordGenerator) {
	raw := password.RawPassword(os.Getenv("ADMIN_PASSWORD"))

	hashed, err := gen.GeneratePassword(raw)
	if err != nil {
		log.Fatalf("failed to generate admin user password: %s", err)
	}

	passwords.ChangePassword(adminUser(), hashed)
}
func initPasswordResetDestinationRepository(destinations password_reset_infra.DestinationRepository) {
	err := destinations.RegisterDestination(adminUser(), password_reset.NewLogDestination())
	if err != nil {
		log.Fatalf("failed to register admin user destination: %s", err)
	}
}

func adminUser() user.User {
	return user.NewUser(user.UserID(os.Getenv("ADMIN_USER_ID")))
}

func newLeveledLogger() logger.LeveledLogger {
	return logger.NewLeveledLogger(os.Getenv("LOG_LEVEL"))
}

func newExtend() extendSecond {
	// パスワードで認証した場合、有効期限 5分、最大延長 14日
	return extendSecond{
		password: credential.TicketExtendWeek(2),
	}
}

func newTicketSigner() signer.TicketSigner {
	privateKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("ticket private key read failed: %s", err)
	}

	publicKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PUBLIC_KEY"))
	if err != nil {
		log.Fatalf("ticket public key read failed: %s", err)
	}

	key, err := signer.NewJWT_ES_512_Key(signer.JWT_Pem{
		PrivateKey: privateKeyPem,
		PublicKey:  publicKeyPem,
	})
	if err != nil {
		log.Fatalf("ticket key parse failed: %s", err)
	}

	jwt := signer.NewJWTSigner(key)
	return signer.NewTicketSigner(jwt)
}
func newApiTokenSigner() signer.ApiTokenSigner {
	pem, err := ioutil.ReadFile(os.Getenv("API_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("app private key read failed: %s", err)
	}

	key, err := signer.NewJWT_ES_512_Key(signer.JWT_Pem{
		PrivateKey: pem,
	})
	if err != nil {
		log.Fatalf("app key parse failed: %s", err)
	}

	jwt := signer.NewJWTSigner(key)
	return signer.NewApiTokenSigner(jwt)
}
func newContentTokenSigner() signer.ContentTokenSigner {
	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		log.Fatalf("aws cloudfront private key read failed: %s", err)
	}

	return signer.NewContentTokenSigner(
		credential.ContentKeyID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)
}
