package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/infra/db/memory"
	"github.com/getto-systems/project-example-id/infra/serializer"
	"github.com/getto-systems/project-example-id/infra/simple_logger"

	auth_handler "github.com/getto-systems/project-example-id/http_handler/auth"

	"github.com/getto-systems/project-example-id/logger"

	"github.com/getto-systems/project-example-id/auth"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"time"
)

type Server struct {
	authCookieDomain auth_handler.CookieDomain

	cors cors.Options
	tls  Tls

	ticketSerializer        serializer.TicketJsonSerializer
	awsCloudFrontSerializer serializer.AwsCloudFrontSerializer

	db memory.MemoryStore
}

type Tls struct {
	cert string
	key  string
}

func main() {
	server, err := initServer()
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
		AuthenticatorFactory: func(r *http.Request) (auth.RenewAuthenticator, error) { return server.initHandler(r) },
	}
}
func authPasswordHandler(server *Server) auth_handler.PasswordHandler {
	return auth_handler.PasswordHandler{
		CookieDomain:         server.authCookieDomain,
		AuthenticatorFactory: func(r *http.Request) (auth.PasswordAuthenticator, error) { return server.initHandler(r) },
	}
}

func initServer() (*Server, error) {
	ticketSerializer, err := initTicketSerializer()
	if err != nil {
		return nil, err
	}

	awsCloudFrontSerializer, err := initAwsCloudFrontSerializer()
	if err != nil {
		return nil, err
	}

	db, err := initDB()
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

		db: db,
	}, nil
}
func initTicketSerializer() (serializer.TicketJsonSerializer, error) {
	return serializer.NewTicketJsonSerializer(), nil
}
func initAwsCloudFrontSerializer() (serializer.AwsCloudFrontSerializer, error) {
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
func initDB() (memory.MemoryStore, error) {
	return memory.NewMemoryStore(), nil
}

// interface methods (auth/renew:Authenticator, auth/password:Authenticator)
type Handler struct {
	server *Server
	logger logger.Logger
}

func (server *Server) initHandler(r *http.Request) (Handler, error) {
	var nullHandler Handler

	logger, err := initLogger(r)
	if err != nil {
		return nullHandler, err
	}

	return Handler{
		server: server,
		logger: logger,
	}, nil
}

type RequestLogEntry struct {
	RequestID string
	RemoteIP  string `json:"remote_ip"`
}

func initLogger(r *http.Request) (logger.Logger, error) {
	requestID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	request := RequestLogEntry{
		RequestID: requestID.String(),
		RemoteIP:  r.RemoteAddr,
	}

	return leveledLogger(os.Getenv("LOG_LEVEL"), request), nil
}
func leveledLogger(level string, request RequestLogEntry) logger.Logger {
	switch level {
	case "DEBUG":
		return simple_logger.NewDebugLogger(request)
	case "INFO":
		return simple_logger.NewInfoLogger(request)
	case "WARNING":
		return simple_logger.NewWarningLogger(request)
	default:
		return simple_logger.NewErrorLogger(request)
	}
}

func (handler Handler) Logger() logger.Logger {
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
