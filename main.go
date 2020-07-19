package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/adapter/logger"
	"github.com/getto-systems/project-example-id/adapter/nonce_generator"
	"github.com/getto-systems/project-example-id/adapter/password_encrypter"
	"github.com/getto-systems/project-example-id/adapter/signer"

	"github.com/getto-systems/project-example-id/http_handler"
	"github.com/getto-systems/project-example-id/http_handler/password_handler"
	"github.com/getto-systems/project-example-id/http_handler/ticket_handler"

	password_core "github.com/getto-systems/project-example-id/password/core"
	password_db "github.com/getto-systems/project-example-id/password/db"
	password_event_log "github.com/getto-systems/project-example-id/password/event_log"
	password_pubsub "github.com/getto-systems/project-example-id/password/pubsub"

	ticket_core "github.com/getto-systems/project-example-id/ticket/core"
	ticket_db "github.com/getto-systems/project-example-id/ticket/db"
	ticket_event_log "github.com/getto-systems/project-example-id/ticket/event_log"
	ticket_pubsub "github.com/getto-systems/project-example-id/ticket/pubsub"

	"github.com/getto-systems/project-example-id/data"
	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"
)

const (
	HEADER_HANDLER = "X-Getto-Example-ID-Handler"
)

type (
	server struct {
		port string

		cors cors.Options
		tls  tls

		handler handler
	}

	tls struct {
		cert string
		key  string
	}

	handler struct {
		ticket   ticket_handler.Handler
		password password_handler.Handler
	}
)

func main() {
	log.Fatal(newServer().listen())
}
func (server server) listen() error {
	handler := server.http_handler()

	if os.Getenv("SERVER_MODE") == "backend" {
		return http.ListenAndServe(
			server.port,
			handler,
		)
	} else {
		return http.ListenAndServeTLS(
			server.port,
			server.tls.cert,
			server.tls.key,
			cors.New(server.cors).Handler(handler),
		)
	}
}
func (server server) http_handler() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", server.handle)
	return router
}

func (server server) handle(w http.ResponseWriter, r *http.Request) {
	handler := r.Header.Get(HEADER_HANDLER)

	switch handler {
	case "password/validate":
		server.handler.password.Validate(w, r)
	case "password/register":
		server.handler.password.Register(w, r)

	case "ticket/extend":
		server.handler.ticket.Extend(w, r)
	case "ticket/shrink":
		server.handler.ticket.Shrink(w, r)

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

		handler: newHandler(),
	}
}

func newHandler() handler {
	appLogger := newAppLogger()

	response := newResponse()

	ticket := newTicketUsecase(appLogger)
	password := newPasswordUsecase(appLogger, ticket)

	return handler{
		ticket:   ticket_handler.NewHandler(appLogger, response, ticket),
		password: password_handler.NewHandler(appLogger, response, password),
	}
}

func newTicketUsecase(appLogger logger.Logger) ticket.Usecase {
	log := ticket_event_log.NewEventLogger(appLogger)
	pub := ticket_pubsub.NewPubSub()
	pub.Subscribe(log)
	db := ticket_db.NewMemoryStore()

	exp := newExpiration()

	signer := newTicketSigner()
	apiTokenSigner := newApiTokenSigner()
	contentTokenSigner := newContentTokenSigner()

	gen := nonce_generator.NewNonceGenerator()

	return ticket_core.NewUsecase(
		pub,
		db,
		exp,

		signer,
		apiTokenSigner,
		contentTokenSigner,

		gen,
	)
}

func newPasswordUsecase(appLogger logger.Logger, ticket ticket.Usecase) password.Usecase {
	log := password_event_log.NewEventLogger(appLogger)
	pub := password_pubsub.NewPubSub()
	pub.Subscribe(log)
	db := password_db.NewMemoryStore()

	encrypter := password_encrypter.NewPasswordEncrypter(10) // bcrypt.DefaultCost

	initAdminPassword(db, encrypter)

	return password_core.NewUsecase(
		pub,
		db,

		encrypter,
		encrypter,

		ticket,
	)
}
func initAdminPassword(db password.DB, gen password.Generator) {
	admin_user_id := os.Getenv("ADMIN_ID")
	admin_password := os.Getenv("ADMIN_PASSWORD")

	p, err := gen.GeneratePassword(password.RawPassword(admin_password))
	if err == nil {
		db.RegisterUserPassword(data.NewUser(data.UserID(admin_user_id)), p)
	}
}

func newAppLogger() logger.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"), log.New(os.Stdout, "", 0))
}

func newResponse() http_handler.Response {
	cookie := http_handler.NewCookie(
		http_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN")),
		http_handler.ContentTokenID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
	)
	return http_handler.NewResponse(cookie)
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
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)
}

func newExpiration() ticket.Expiration {
	return ticket.NewExpiration(ticket.ExpirationParam{
		Expires:     data.Minute(5),
		ExtendLimit: data.Day(14),
	})
}
