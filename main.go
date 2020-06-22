package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/infra/db"
	"github.com/getto-systems/project-example-id/infra/logger"
	"github.com/getto-systems/project-example-id/infra/password"
	"github.com/getto-systems/project-example-id/infra/policy"
	"github.com/getto-systems/project-example-id/infra/pubsub"
	"github.com/getto-systems/project-example-id/infra/serializer"

	"github.com/getto-systems/project-example-id/http_handler/auth_handler"

	"github.com/getto-systems/project-example-id/user/authenticate"
	"github.com/getto-systems/project-example-id/user/authorize"

	"github.com/getto-systems/project-example-id/user"
	"github.com/getto-systems/project-example-id/user/subscriber"
)

type Server struct {
	port string

	cors cors.Options
	tls  Tls

	logger logger.Logger

	cookieDomain auth_handler.CookieDomain

	issuer Issuer

	authorizerFactory    authorize.AuthorizerFactory
	authenticatorFactory AuthenticatorFactory
}

type Tls struct {
	cert string
	key  string
}

type Issuer struct {
	awsCloudFront auth_handler.AwsCloudFrontIssuer
	app           auth_handler.AppIssuer
}

type AuthenticatorFactory struct {
	password authenticate.PasswordAuthenticatorFactory
	renew    authenticate.RenewAuthenticatorFactory
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/healthz", healthz).Methods("GET")
	router.HandleFunc("/auth/renew", server.authRenewHandler()).Methods("POST")
	router.HandleFunc("/auth/password", server.authPasswordHandler()).Methods("POST")

	handler := cors.New(server.cors).Handler(router)

	log.Fatal(listen(server, handler))
}
func listen(server Server, handler http.Handler) error {
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
			handler,
		)
	}
}
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := "\"OK\""

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}
func (server Server) authRenewHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := auth_handler.RenewHandler{
			AuthHandler: server.authHandler(w, r),

			AuthenticatorFactory: server.authenticatorFactory.renew,
		}
		handler.Handle()
	}
}
func (server Server) authPasswordHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := auth_handler.PasswordHandler{
			AuthHandler: server.authHandler(w, r),

			AuthenticatorFactory: server.authenticatorFactory.password,
		}
		handler.Handle()
	}
}
func (server Server) authHandler(w http.ResponseWriter, r *http.Request) auth_handler.AuthHandler {
	return auth_handler.AuthHandler{
		Logger: server.logger,

		HttpResponseWriter: w,
		HttpRequest:        r,

		CookieDomain: server.cookieDomain,

		AwsCloudFrontIssuer: server.issuer.awsCloudFront,
		AppIssuer:           server.issuer.app,

		Request: auth_handler.Request(r),

		AuthorizerFactory: server.authorizerFactory,
	}
}

func NewServer() (Server, error) {
	appLogger := NewAppLogger()

	awsCloudFrontIssuer, err := NewAwsCloudFrontIssuer(appLogger)
	if err != nil {
		return Server{}, err
	}

	appIssuer, err := NewAppIssuer(appLogger)
	if err != nil {
		return Server{}, err
	}

	db, err := NewDB()
	if err != nil {
		return Server{}, err
	}

	pubsub, err := NewPubSub()
	if err != nil {
		return Server{}, err
	}
	pubsub.Subscribe(subscriber.NewUserLogger(appLogger))

	ticketSerializer, err := NewTicketSerializer(appLogger)
	if err != nil {
		return Server{}, err
	}

	passwordEncrypter := NewPasswordEncrypter()

	passwordRepository := NewPasswordRepository(db, passwordEncrypter)
	issuerRepository := NewIssuerRepository(db, ticketSerializer)

	ticketAuthorizer := NewTicketAuthorizer(ticketSerializer)

	userFactory := user.NewUserFactory(pubsub)

	return Server{
		port: ":8080",

		cors: cors.Options{
			AllowedOrigins:   []string{os.Getenv("ORIGIN")},
			AllowedMethods:   []string{"POST"},
			AllowCredentials: true,
		},
		tls: Tls{
			cert: os.Getenv("TLS_CERT"),
			key:  os.Getenv("TLS_KEY"),
		},

		logger: appLogger,

		cookieDomain: auth_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN")),

		issuer: Issuer{
			awsCloudFront: awsCloudFrontIssuer,
			app:           appIssuer,
		},

		authorizerFactory: authorize.NewAuthorizerFactory(ticketAuthorizer, userFactory),
		authenticatorFactory: AuthenticatorFactory{
			password: authenticate.NewPasswordAuthenticatorFactory(passwordRepository, issuerRepository, userFactory),
			renew:    authenticate.NewRenewAuthenticatorFactory(issuerRepository, userFactory),
		},
	}, nil
}
func NewAppLogger() logger.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"), log.New(os.Stdout, "", 0))
}
func NewAwsCloudFrontIssuer(appLogger logger.Logger) (auth_handler.AwsCloudFrontIssuer, error) {
	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		appLogger.Debugf(nil, "aws cloudfront private key read failed: %s", err)
		return auth_handler.AwsCloudFrontIssuer{}, err
	}

	serializer := serializer.NewAwsCloudFrontSerializer(
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)

	return auth_handler.NewAwsCloudFrontIssuer(
		auth_handler.AwsCloudFrontKeyPairID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
		serializer,
	), nil
}
func NewAppIssuer(appLogger logger.Logger) (auth_handler.AppIssuer, error) {
	pem, err := ioutil.ReadFile(os.Getenv("APP_PRIVATE_KEY"))
	if err != nil {
		appLogger.Debugf(nil, "app private key read failed: %s", err)
		return auth_handler.AppIssuer{}, err
	}

	key, err := serializer.NewJWT_ES_512_Key(serializer.JWT_Pem{
		PrivateKey: pem,
	})
	if err != nil {
		appLogger.Debugf(nil, "app key parse failed: %s", err)
		return auth_handler.AppIssuer{}, err
	}

	jwt := serializer.NewJWTSerializer(key)
	app := serializer.NewAppSerializer(jwt)

	return auth_handler.NewAppIssuer(app), nil
}
func NewDB() (*db.MemoryStore, error) {
	return db.NewMemoryStore(), nil
}
func NewPubSub() (*pubsub.SyncPubSub, error) {
	return pubsub.NewSyncPubSub(), nil
}
func NewTicketSerializer(appLogger logger.Logger) (serializer.TicketSerializer, error) {
	privateKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PRIVATE_KEY"))
	if err != nil {
		appLogger.Debugf(nil, "ticket private key read failed: %s", err)
		return serializer.TicketSerializer{}, err
	}

	publicKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PUBLIC_KEY"))
	if err != nil {
		appLogger.Debugf(nil, "ticket public key read failed: %s", err)
		return serializer.TicketSerializer{}, err
	}

	key, err := serializer.NewJWT_ES_512_Key(serializer.JWT_Pem{
		PrivateKey: privateKeyPem,
		PublicKey:  publicKeyPem,
	})
	if err != nil {
		appLogger.Debugf(nil, "ticket key parse failed: %s", err)
		return serializer.TicketSerializer{}, err
	}

	jwt := serializer.NewJWTSerializer(key)
	return serializer.NewTicketSerializer(jwt), nil
}
func NewPasswordEncrypter() password.PasswordEncrypter {
	return password.NewPasswordEncrypter(10) // bcrypt.DefaultCost
}
func NewPasswordRepository(db *db.MemoryStore, passwordEncrypter password.PasswordEncrypter) authenticate.PasswordRepository {
	passwordMatcherFactory := user.NewPasswordMatcherFactory(passwordEncrypter)
	return authenticate.NewPasswordRepository(db, passwordMatcherFactory)
}
func NewIssuerRepository(db *db.MemoryStore, ticketSerializer user.TicketSerializer) authenticate.IssuerRepository {
	issuerFactory := user.NewIssuerFactory(ticketSerializer)
	return authenticate.NewIssuerRepository(db, issuerFactory)
}
func NewTicketAuthorizer(ticketSerializer serializer.TicketSerializer) user.TicketAuthorizer {
	checker := policy.NewPolicyChecker()
	return user.NewTicketAuthorizer(ticketSerializer, checker)
}
