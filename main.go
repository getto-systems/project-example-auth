package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/adapter/http_handler"
	"github.com/getto-systems/project-example-id/adapter/logger"
	"github.com/getto-systems/project-example-id/adapter/message"
	"github.com/getto-systems/project-example-id/adapter/nonce_generator"
	"github.com/getto-systems/project-example-id/adapter/password_encrypter"
	"github.com/getto-systems/project-example-id/adapter/reset_session_generator"
	"github.com/getto-systems/project-example-id/adapter/signer"

	"github.com/getto-systems/project-example-id/client"

	"github.com/getto-systems/project-example-id/data/api_token"
	"github.com/getto-systems/project-example-id/data/password"
	"github.com/getto-systems/project-example-id/data/password_reset"
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

const (
	HEADER_HANDLER = "X-Getto-Example-ID-Handler"
)

type (
	server struct {
		port string
		cors cors.Options
		tls  tls

		cookieDomain http_handler.CookieDomain
		backend      client.Backend
	}

	tls struct {
		cert string
		key  string
	}
)

func main() {
	log.Fatal(newServer().listen())
}
func (server server) listen() error {
	mux := server.mux()

	if os.Getenv("SERVER_MODE") == "backend" {
		return http.ListenAndServe(
			server.port,
			mux,
		)
	} else {
		return http.ListenAndServeTLS(
			server.port,
			server.tls.cert,
			server.tls.key,
			cors.New(server.cors).Handler(mux),
		)
	}
}
func (server server) mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handle)
	return mux
}

func (server server) handle(w http.ResponseWriter, r *http.Request) {
	h := http_handler.NewHandler(w, r)
	c := client.NewClient(server.backend, http_handler.NewCredentialHandler(server.cookieDomain, w, r))

	switch r.Header.Get(HEADER_HANDLER) {
	case "Renew":
		client.NewRenew(c).Renew(http_handler.NewRenew(h))
	case "Logout":
		client.NewLogout(c).Logout(http_handler.NewLogout(h))

	case "PasswordLogin":
		client.NewPasswordLogin(c).Login(http_handler.NewPasswordLogin(h))

	case "PasswordChange/GetLogin":
		client.NewPasswordChange(c).GetLogin(http_handler.NewPasswordChange(h))
	case "PasswordChange/Change":
		client.NewPasswordChange(c).Change(http_handler.NewPasswordChange(h))

	case "PasswordReset/CreateSession":
		client.NewPasswordReset(c).CreateSession(http_handler.NewPasswordReset(h))
	case "PasswordReset/SendToken":
		client.NewPasswordReset(c).SendToken(http_handler.NewPasswordReset(h))
	case "PasswordReset/GetStatus":
		client.NewPasswordReset(c).GetStatus(http_handler.NewPasswordReset(h))
	case "PasswordReset/Reset":
		client.NewPasswordReset(c).Reset(http_handler.NewPasswordReset(h))

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "\"OK\"")
	}
}

func newServer() server {
	return server{
		port: ":8080",

		cors: cors.Options{
			AllowedOrigins:   []string{os.Getenv("ORIGIN")},
			AllowedMethods:   []string{"POST"},
			AllowedHeaders:   []string{HEADER_HANDLER},
			AllowCredentials: true,
		},

		tls: tls{
			cert: os.Getenv("TLS_CERT"),
			key:  os.Getenv("TLS_KEY"),
		},

		cookieDomain: http_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN")),
		backend:      newBackend(),
	}
}

func newBackend() client.Backend {
	appLogger := newAppLogger()

	return client.NewBackend(
		newTicketAction(appLogger),
		newApiTokenAction(appLogger),
		newUserAction(appLogger),
		newPasswordAction(appLogger),
		newPasswordResetAction(appLogger),
	)
}
func newTicketAction(appLogger logger.Logger) client.TicketAction {
	return client.NewTicketAction(
		ticket_log.NewLogger(appLogger),

		newTicketSigner(),
		ticket.ExpirationParam{
			Expires:     time.Minute(5),
			ExtendLimit: time.Day(14),
		},
		nonce_generator.NewNonceGenerator(),

		ticket_repository_ticket.NewMemoryStore(),
	)
}
func newApiTokenAction(appLogger logger.Logger) client.ApiTokenAction {
	api_users := api_token_repository_api_user.NewMemoryStore()

	initApiUserRepository(api_users)

	return client.NewApiTokenAction(
		api_token_log.NewLogger(appLogger),

		newApiTokenSigner(),
		newContentTokenSigner(),

		api_users,
	)
}
func newUserAction(appLogger logger.Logger) client.UserAction {
	users := user_repository_user.NewMemoryStore()

	initUserRepository(users)

	return client.NewUserAction(
		user_log.NewLogger(appLogger),

		users,
	)
}
func newPasswordAction(appLogger logger.Logger) client.PasswordAction {
	enc := password_encrypter.NewEncrypter(10) // bcrypt.DefaultCost
	passwords := password_repository_password.NewMemoryStore()

	initPasswordRepository(passwords, enc)

	return client.NewPasswordAction(
		password_log.NewLogger(appLogger),

		enc,

		passwords,
	)
}
func newPasswordResetAction(appLogger logger.Logger) client.PasswordResetAction {
	destinations := password_reset_repository_destination.NewMemoryStore()

	initPasswordResetDestinationRepository(destinations)

	return client.NewPasswordResetAction(
		password_reset_log.NewLogger(appLogger),

		time.Minute(30),
		reset_session_generator.NewGenerator(),

		password_reset_repository_session.NewMemoryStore(),
		destinations,

		password_reset_job_queue_send_token.NewMemoryQueue(),
		password_reset_sender.NewTokenSender(message.NewLogMessage()),
	)
}

func initUserRepository(users user.UserRepository) {
	login := user.NewLogin(user.LoginID(os.Getenv("ADMIN_LOGIN_ID")))

	err := users.RegisterUser(adminUser(), login)
	if err != nil {
		log.Fatalf("failed to register admin user: %s", err)
	}
}
func initApiUserRepository(api_users api_token.ApiUserRepository) {
	err := api_users.RegisterApiRoles(adminUser(), api_token.ApiRoles([]string{"admin"}))
	if err != nil {
		log.Fatalf("failed to register admin user api roles: %s", err)
	}
}
func initPasswordRepository(passwords password.PasswordRepository, gen password.PasswordGenerator) {
	raw := password.RawPassword(os.Getenv("ADMIN_PASSWORD"))

	hashed, err := gen.GeneratePassword(raw)
	if err != nil {
		log.Fatalf("failed to generate admin user password: %s", err)
	}

	passwords.ChangePassword(adminUser(), hashed)
}
func initPasswordResetDestinationRepository(destinations password_reset.DestinationRepository) {
	err := destinations.RegisterDestination(adminUser(), password_reset.NewLogDestination())
	if err != nil {
		log.Fatalf("failed to register admin user destination: %s", err)
	}
}

func adminUser() user.User {
	return user.NewUser(user.UserID(os.Getenv("ADMIN_USER_ID")))
}

func newAppLogger() logger.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"))
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
	pem, err := ioutil.ReadFile(os.Getenv("APP_PRIVATE_KEY"))
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
		api_token.ContentKeyID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)
}
