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

	"github.com/getto-systems/project-example-id/event_log"

	"github.com/getto-systems/project-example-id/password"
	"github.com/getto-systems/project-example-id/ticket"

	password_db "github.com/getto-systems/project-example-id/password/db"
	ticket_db "github.com/getto-systems/project-example-id/ticket/db"

	password_pubsub "github.com/getto-systems/project-example-id/password/pubsub"
	ticket_pubsub "github.com/getto-systems/project-example-id/ticket/pubsub"

	"github.com/getto-systems/project-example-id/data"
)

const (
	HEADER_HANDLER = "X-Getto-Example-ID-Handler"
)

type (
	server struct {
		port string

		cors cors.Options
		tls  tls

		logger logger.Logger

		response http_handler.Response

		password passwordUsecase
		ticket   ticketUsecase
	}

	tls struct {
		cert string
		key  string
	}

	passwordUsecase struct {
		verifier password.Verifier
		register password.Register
	}

	ticketUsecase struct {
		verifier ticket.Verifier
		extender ticket.Extender
		shrinker ticket.Shrinker

		issuer             ticket.Issuer
		apiTokenIssuer     ticket.ApiTokenIssuer
		contentTokenIssuer ticket.ContentTokenIssuer
	}
)

func main() {
	log.Fatal(NewServer().listen())
}
func (server server) listen() error {
	handler := server.handler()

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
func (server server) handler() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", server.http_handler)
	return router
}

func (server server) http_handler(w http.ResponseWriter, r *http.Request) {
	handler := r.Header.Get(HEADER_HANDLER)

	switch handler {
	case "password/verify":
		server.password_verify().Handle(w, r)
	case "password/register":
		server.password_register().Handle(w, r)

	case "ticket/extend":
		server.ticket_extend().Handle(w, r)

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "\"OK\"")
	}
}
func (server server) password_verify() password_handler.Verify {
	return password_handler.NewVerify(
		server.logger,
		server.response,
		password.NewPasswordVerifier(
			server.password.verifier,
			ticket.NewTicketIssuer(server.ticket.issuer),
			server.ticket.apiTokenIssuer,
			server.ticket.contentTokenIssuer,
		),
	)
}
func (server server) password_register() password_handler.Register {
	return password_handler.NewRegister(
		server.logger,
		server.response,
		password.NewPasswordRegister(
			ticket.NewTicketVerifier(server.ticket.verifier),
			server.password.verifier,
			server.password.register,
		),
	)
}

func (server server) ticket_extend() ticket_handler.Extend {
	return ticket_handler.NewExtend(
		server.logger,
		server.response,
		ticket.NewTicketExtender(
			server.ticket.verifier,
			server.ticket.extender,
			server.ticket.apiTokenIssuer,
			server.ticket.contentTokenIssuer,
		),
	)
}
func (server server) ticket_shrink() ticket_handler.Shrink {
	return ticket_handler.NewShrink(
		server.logger,
		server.response,
		ticket.NewTicketShrinker(
			server.ticket.verifier,
			server.ticket.shrinker,
		),
	)
}

func NewServer() server {
	appLogger := NewAppLogger()

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

		logger: appLogger,

		response: NewResponse(),

		password: NewPasswordUsecase(appLogger),
		ticket:   NewTicketUsecase(appLogger),
	}
}

func NewAppLogger() logger.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"), log.New(os.Stdout, "", 0))
}

func NewResponse() http_handler.Response {
	cookie := http_handler.NewCookie(
		http_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN")),
		http_handler.ContentTokenID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
	)
	return http_handler.NewResponse(cookie)
}

func NewPasswordUsecase(appLogger logger.Logger) passwordUsecase {
	log := event_log.NewPasswordEventLogger(appLogger)
	pub := password_pubsub.NewPasswordPubSub()
	pub.Subscribe(log)
	db := password_db.NewPasswordStore()
	encrypter := password_encrypter.NewPasswordEncrypter(10) // bcrypt.DefaultCost

	// test
	p, err := encrypter.GeneratePassword("password")
	if err == nil {
		db.RegisterUserPassword(data.NewUser("admin"), p)
	}

	return passwordUsecase{
		verifier: password.NewVerifier(pub, db, encrypter),
		register: password.NewRegister(pub, db, encrypter),
	}
}

func NewTicketUsecase(appLogger logger.Logger) ticketUsecase {
	log := event_log.NewTicketEventLogger(appLogger)
	pub := ticket_pubsub.NewTicketPubSub()
	pub.Subscribe(log)
	db := ticket_db.NewTicketStore()

	signer := NewTicketSigner()
	apiTokenSigner := NewApiTokenSigner()
	contentTokenSigner := NewContentTokenSigner()

	expiration := NewExpiration()

	gen := NewNonceGenerator()

	return ticketUsecase{
		verifier: ticket.NewVerifier(pub, signer),
		extender: ticket.NewExtender(pub, db, signer, expiration),
		shrinker: ticket.NewShrinker(pub, db),

		issuer:             ticket.NewIssuer(pub, db, signer, expiration, gen),
		apiTokenIssuer:     ticket.NewApiTokenIssuer(pub, db, apiTokenSigner),
		contentTokenIssuer: ticket.NewContentTokenIssuer(pub, contentTokenSigner),
	}
}

func NewTicketSigner() signer.TicketSigner {
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
func NewApiTokenSigner() signer.ApiTokenSigner {
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
func NewContentTokenSigner() signer.ContentTokenSigner {
	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		log.Fatalf("aws cloudfront private key read failed: %s", err)
	}

	return signer.NewContentTokenSigner(
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)
}

func NewExpiration() ticket.Expiration {
	return ticket.NewExpiration(ticket.ExpirationParam{
		Expires:     data.Minute(5),
		ExtendLimit: data.Day(14),
	})
}

func NewNonceGenerator() nonce_generator.NonceGenerator {
	return nonce_generator.NewNonceGenerator()
}
