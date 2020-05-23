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
	"github.com/getto-systems/project-example-id/infra/serializer"

	"github.com/getto-systems/project-example-id/http_handler"
	"github.com/getto-systems/project-example-id/http_handler/auth"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type server struct {
	cookieDomain http_handler.CookieDomain

	cors cors.Options
	tls  tls

	ticketSerializer        serializer.TicketJsonSerializer
	awsCloudFrontSerializer serializer.AwsCloudFrontSerializer

	db memory.MemoryStore
}

type tls struct {
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
func authRenewHandler(server *server) auth.RenewHandler {
	return auth.RenewHandler{
		CookieDomain:  server.cookieDomain,
		Authenticator: server,
	}
}
func authPasswordHandler(server *server) auth.PasswordHandler {
	return auth.PasswordHandler{
		CookieDomain:  server.cookieDomain,
		Authenticator: server,
	}
}

func initServer() (*server, error) {
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

	return &server{
		cookieDomain: http_handler.CookieDomain(os.Getenv("DOMAIN")),

		cors: cors.Options{
			AllowedOrigins:   []string{os.Getenv("ORIGIN")},
			AllowedMethods:   []string{"POST"},
			AllowCredentials: true,
		},
		tls: tls{
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
func (server *server) UserFactory() user.UserFactory {
	return user.NewUserFactory(server.db)
}

func (server *server) UserPasswordFactory() user.UserPasswordFactory {
	return user.NewUserPasswordFactory(server.db)
}

func (server *server) TicketSerializer() token.TicketSerializer {
	return server.ticketSerializer
}

func (server *server) AwsCloudFrontSerializer() token.AwsCloudFrontSerializer {
	return server.awsCloudFrontSerializer
}

func (server *server) Now() time.Time {
	return time.Now().UTC()
}
