package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/getto-systems/project-example-auth"

	"github.com/getto-systems/project-example-auth/x_http_server/http_handler"
	"github.com/getto-systems/project-example-auth/x_http_server/logger"
	"github.com/getto-systems/project-example-auth/x_http_server/message"
	"github.com/getto-systems/project-example-auth/x_http_server/nonce_generator"
	"github.com/getto-systems/project-example-auth/x_http_server/password_encrypter"
	"github.com/getto-systems/project-example-auth/x_http_server/reset_session_generator"
	"github.com/getto-systems/project-example-auth/x_http_server/signer"

	"github.com/getto-systems/project-example-auth/credential/log"
	"github.com/getto-systems/project-example-auth/credential/repository/api_user"
	"github.com/getto-systems/project-example-auth/password/log"
	"github.com/getto-systems/project-example-auth/password/repository/password"
	"github.com/getto-systems/project-example-auth/password_reset/job_queue/send_token"
	"github.com/getto-systems/project-example-auth/password_reset/log"
	"github.com/getto-systems/project-example-auth/password_reset/repository/destination"
	"github.com/getto-systems/project-example-auth/password_reset/repository/session"
	"github.com/getto-systems/project-example-auth/password_reset/sender"
	"github.com/getto-systems/project-example-auth/ticket/log"
	"github.com/getto-systems/project-example-auth/ticket/repository/ticket"
	"github.com/getto-systems/project-example-auth/user/log"
	"github.com/getto-systems/project-example-auth/user/repository/user"
	"github.com/getto-systems/project-example-auth/y_static/repository/env"
	"github.com/getto-systems/project-example-auth/y_static/repository/secret"

	credential_infra "github.com/getto-systems/project-example-auth/credential/infra"
	password_infra "github.com/getto-systems/project-example-auth/password/infra"
	password_reset_infra "github.com/getto-systems/project-example-auth/password_reset/infra"
	user_infra "github.com/getto-systems/project-example-auth/user/infra"

	"github.com/getto-systems/project-example-auth/credential"
	"github.com/getto-systems/project-example-auth/password"
	"github.com/getto-systems/project-example-auth/password_reset"
	"github.com/getto-systems/project-example-auth/ticket"
	"github.com/getto-systems/project-example-auth/user"
	"github.com/getto-systems/project-example-auth/y_static"

	"github.com/getto-systems/project-example-auth/credential/core"
	"github.com/getto-systems/project-example-auth/password/core"
	"github.com/getto-systems/project-example-auth/password_reset/core"
	"github.com/getto-systems/project-example-auth/ticket/core"
	"github.com/getto-systems/project-example-auth/user/core"
	"github.com/getto-systems/project-example-auth/y_static/core"
)

const (
	HEADER_HANDLER = "X-Getto-Example-ID-Handler"
)

type (
	server struct {
		port string

		cookieDomain http_handler.CookieDomain
		backend      auth.Backend
	}

	env struct {
		logLevel   string
		secretName string
	}
	secret struct {
		admin      adminSecret
		cookie     cookieSecret
		ticket     ticketSecret
		api        apiSecret
		cloudfront cloudfrontSecret
	}
	adminSecret struct {
		userID   string
		loginID  string
		password string
	}
	cookieSecret struct {
		domain string
	}
	ticketSecret struct {
		privateKey []byte
		publicKey  []byte
	}
	apiSecret struct {
		privateKey []byte
	}
	cloudfrontSecret struct {
		keyPairID   string
		privateKey  []byte
		resourceURL string
	}

	infra struct {
		logger logger.LeveledLogger
		extend extendSecond

		env    static.Env
		secret static.Secret
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
	u := auth.NewUsecase(server.backend, http_handler.NewCredentialHandler(server.cookieDomain, w, r))

	switch r.Header.Get(HEADER_HANDLER) {
	case "Renew":
		auth.NewRenew(u).Renew(http_handler.NewRenew(h))
	case "Logout":
		auth.NewLogout(u).Logout(http_handler.NewLogout(h))

	case "PasswordLogin":
		auth.NewPasswordLogin(u).Login(http_handler.NewPasswordLogin(h))

	case "PasswordChange/GetLogin":
		auth.NewPasswordChange(u).GetLogin(http_handler.NewPasswordChange(h))
	case "PasswordChange/Change":
		auth.NewPasswordChange(u).Change(http_handler.NewPasswordChange(h))

	case "PasswordReset/CreateSession":
		auth.NewPasswordReset(u).CreateSession(http_handler.NewPasswordReset(h))
	case "PasswordReset/SendToken":
		auth.NewPasswordReset(u).SendToken(http_handler.NewPasswordReset(h))
	case "PasswordReset/GetStatus":
		auth.NewPasswordReset(u).GetStatus(http_handler.NewPasswordReset(h))
	case "PasswordReset/Reset":
		auth.NewPasswordReset(u).Reset(http_handler.NewPasswordReset(h))

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "\"OK\"")
	}
}

func newServer() server {
	infra := newInfra()

	return server{
		port: ":8080",

		cookieDomain: http_handler.CookieDomain(infra.secret.Cookie.Domain),
		backend:      newBackend(infra),
	}
}

func newBackend(infra infra) auth.Backend {
	return auth.NewBackend(
		infra.newTicketAction(),
		infra.newCredentialAction(),
		infra.newUserAction(),
		infra.newPasswordAction(),
		infra.newPasswordResetAction(),
	)
}

func newInfra() infra {
	staticAction := static_core.NewAction(
		static_repository_env.NewOSEnv(),
		static_repository_secret.NewSecretManager(),
	)

	env, err := staticAction.GetEnv()
	if err != nil {
		log.Fatalf("failed to get env: %s", err)
	}

	secret, err := staticAction.GetSecret(env.SecretName)
	if err != nil {
		log.Fatalf("failed to get secret: %s", err)
	}

	return infra{
		logger: newLeveledLogger(env.LogLevel),
		extend: newExtend(),

		env:    env,
		secret: secret,
	}
}
func newLeveledLogger(level string) logger.LeveledLogger {
	return logger.NewLeveledLogger(level)
}
func newExtend() extendSecond {
	// パスワードで認証した場合、有効期限 5分、最大延長 14日
	return extendSecond{
		password: credential.TicketExtendWeek(2),
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

	infra.initApiUserRepository(apiUsers)

	return credential_core.NewAction(
		credential_log.NewLogger(infra.logger),

		infra.newTicketSigner(),
		infra.newApiTokenSigner(),
		infra.newContentTokenSigner(),

		apiUsers,
	)
}
func (infra infra) newUserAction() user.Action {
	users := user_repository_user.NewMemoryStore()

	infra.initUserRepository(users)

	return user_core.NewAction(
		user_log.NewLogger(infra.logger),

		users,
	)
}
func (infra infra) newPasswordAction() password.Action {
	encrypter := password_encrypter.NewEncrypter(10) // bcrypt.DefaultCost
	passwords := password_repository_password.NewMemoryStore()

	infra.initPasswordRepository(passwords, encrypter)

	return password_core.NewAction(
		password_log.NewLogger(infra.logger),

		infra.extend.password,
		encrypter,

		passwords,
	)
}
func (infra infra) newPasswordResetAction() password_reset.Action {
	destinations := password_reset_repository_destination.NewMemoryStore()

	infra.initPasswordResetDestinationRepository(destinations)

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

func (infra infra) initUserRepository(users user_infra.UserRepository) {
	login := user.NewLogin(user.LoginID(infra.secret.Admin.LoginID))

	err := users.RegisterUser(infra.adminUser(), login)
	if err != nil {
		log.Fatalf("failed to register admin user: %s", err)
	}
}
func (infra infra) initApiUserRepository(apiUsers credential_infra.ApiUserRepository) {
	err := apiUsers.RegisterApiRoles(infra.adminUser(), credential.ApiRoles([]string{"admin"}))
	if err != nil {
		log.Fatalf("failed to register admin user api roles: %s", err)
	}
}
func (infra infra) initPasswordRepository(passwords password_infra.PasswordRepository, gen password_infra.PasswordGenerator) {
	raw := password.RawPassword(infra.secret.Admin.Password)

	hashed, err := gen.GeneratePassword(raw)
	if err != nil {
		log.Fatalf("failed to generate admin user password: %s", err)
	}

	passwords.ChangePassword(infra.adminUser(), hashed)
}
func (infra infra) initPasswordResetDestinationRepository(destinations password_reset_infra.DestinationRepository) {
	err := destinations.RegisterDestination(infra.adminUser(), password_reset.NewLogDestination())
	if err != nil {
		log.Fatalf("failed to register admin user destination: %s", err)
	}
}
func (infra infra) adminUser() user.User {
	return user.NewUser(user.UserID(infra.secret.Admin.UserID))
}

func (infra infra) newTicketSigner() signer.TicketSigner {
	key, err := signer.NewJWT_ES_512_Key(signer.JWT_Pem{
		PrivateKey: infra.secret.Ticket.PrivateKey,
		PublicKey:  infra.secret.Ticket.PublicKey,
	})
	if err != nil {
		log.Fatalf("ticket key parse failed: %s", err)
	}

	jwt := signer.NewJWTSigner(key)
	return signer.NewTicketSigner(jwt)
}
func (infra infra) newApiTokenSigner() signer.ApiTokenSigner {
	key, err := signer.NewJWT_ES_512_Key(signer.JWT_Pem{
		PrivateKey: infra.secret.Api.PrivateKey,
	})
	if err != nil {
		log.Fatalf("app key parse failed: %s", err)
	}

	jwt := signer.NewJWTSigner(key)
	return signer.NewApiTokenSigner(jwt)
}
func (infra infra) newContentTokenSigner() signer.ContentTokenSigner {
	return signer.NewContentTokenSigner(
		credential.ContentKeyID(infra.secret.Cloudfront.KeyPairID),
		infra.secret.Cloudfront.PrivateKey,
		infra.secret.Cloudfront.ResourceURL,
	)
}
