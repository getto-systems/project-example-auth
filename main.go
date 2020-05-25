package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/infra/db/memory"

	"github.com/getto-systems/project-example-id/adapter/logger"
	"github.com/getto-systems/project-example-id/adapter/serializer"

	auth_handler "github.com/getto-systems/project-example-id/http_handler/auth"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/journal"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type Server struct {
	authCookieDomain auth_handler.CookieDomain

	cors cors.Options
	tls  Tls

	ticketSerializer        serializer.TicketJsonSerializer
	awsCloudFrontSerializer serializer.AwsCloudFrontSerializer

	log Log

	db memory.MemoryStore
}

type Tls struct {
	cert string
	key  string
}

type Log struct {
	level  string
	logger *log.Logger
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/auth/renew", authRenewHandler(server).Handle).Methods("POST")
	router.HandleFunc("/auth/password", authPasswordHandler(server).Handle).Methods("POST")

	handler := cors.New(server.cors).Handler(router)

	log.Fatal(http.ListenAndServeTLS(
		":8080",
		server.tls.cert,
		server.tls.key,
		handler,
	))
}
func authRenewHandler(server *Server) auth_handler.RenewHandler {
	return auth_handler.RenewHandler{
		CookieDomain:         server.authCookieDomain,
		AuthenticatorFactory: func(r *http.Request) (auth.RenewAuthenticator, error) { return server.NewHandler(r) },
	}
}
func authPasswordHandler(server *Server) auth_handler.PasswordHandler {
	return auth_handler.PasswordHandler{
		CookieDomain:         server.authCookieDomain,
		AuthenticatorFactory: func(r *http.Request) (auth.PasswordAuthenticator, error) { return server.NewHandler(r) },
	}
}

func NewServer() (*Server, error) {
	ticketSerializer, err := NewTicketSerializer()
	if err != nil {
		return nil, err
	}

	awsCloudFrontSerializer, err := NewAwsCloudFrontSerializer()
	if err != nil {
		return nil, err
	}

	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	return &Server{
		authCookieDomain: auth_handler.CookieDomain(os.Getenv("DOMAIN")),

		cors: cors.Options{
			AllowedOrigins:   []string{os.Getenv("ORIGIN")},
			AllowedMethods:   []string{"POST"},
			AllowCredentials: true,
		},
		tls: Tls{
			cert: os.Getenv("TLS_CERT"),
			key:  os.Getenv("TLS_KEY"),
		},

		ticketSerializer:        ticketSerializer,
		awsCloudFrontSerializer: awsCloudFrontSerializer,

		log: Log{
			level:  os.Getenv("LOG_LEVEL"),
			logger: log.New(os.Stdout, "", 0),
		},

		db: db,
	}, nil
}
func NewTicketSerializer() (serializer.TicketJsonSerializer, error) {
	return serializer.NewTicketJsonSerializer(), nil
}
func NewAwsCloudFrontSerializer() (serializer.AwsCloudFrontSerializer, error) {
	var nullSerializer serializer.AwsCloudFrontSerializer

	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		return nullSerializer, err
	}

	return serializer.NewAwsCloudFrontSerializer(
		serializer.AwsCloudFrontPem(pem),
		serializer.AwsCloudFrontBaseURL(os.Getenv("AWS_CLOUDFRONT_BASE_URL")),
		token.AwsCloudFrontKeyPairID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
	), nil
}
func NewDB() (memory.MemoryStore, error) {
	return memory.NewMemoryStore(), nil
}

// interface methods (auth/renew:Authenticator, auth/password:Authenticator)
type Handler struct {
	server *Server
	logger journal.Logger
}

func (server *Server) NewHandler(r *http.Request) (Handler, error) {
	var nullHandler Handler

	logger, err := logger.NewLogger(server.log.level, server.log.logger, r)
	if err != nil {
		return nullHandler, err
	}

	return Handler{
		server: server,
		logger: logger,
	}, nil
}

func (handler Handler) Logger() journal.Logger {
	return handler.logger
}

func (handler Handler) TicketSerializer() token.TicketSerializer {
	return handler.server.ticketSerializer
}

func (handler Handler) AwsCloudFrontSerializer() token.AwsCloudFrontSerializer {
	return handler.server.awsCloudFrontSerializer
}

func (handler Handler) UserFactory() user.UserFactory {
	return user.NewUserFactory(handler.server.db)
}

func (handler Handler) UserPasswordFactory() user.UserPasswordFactory {
	return user.NewUserPasswordFactory(handler.server.db)
}

func (handler Handler) Now() time.Time {
	return time.Now().UTC()
}
