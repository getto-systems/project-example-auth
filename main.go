package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/infra/db/memory"
	"github.com/getto-systems/project-example-id/infra/tokener"

	"github.com/getto-systems/project-example-id/http_handler"

	auth_password "github.com/getto-systems/project-example-id/http_handler/auth/password"
	auth_renew "github.com/getto-systems/project-example-id/http_handler/auth/renew"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
	user_password "github.com/getto-systems/project-example-id/user/password"
)

type server struct {
	cookieDomain http_handler.CookieDomain

	cors cors.Options
	tls  tls

	ticketTokener        tokener.TicketJsonTokener
	awsCloudFrontTokener tokener.AwsCloudFrontTokener

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
func authRenewHandler(server *server) auth_renew.Handler {
	return auth_renew.Handler{
		CookieDomain:  server.cookieDomain,
		Authenticator: server,
	}
}
func authPasswordHandler(server *server) auth_password.Handler {
	return auth_password.Handler{
		CookieDomain:  server.cookieDomain,
		Authenticator: server,
	}
}

func initServer() (*server, error) {
	ticketTokener, err := initTicketTokener()
	if err != nil {
		return nil, err
	}

	awsCloudFrontTokener, err := initAwsCloudFrontTokener()
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

		ticketTokener:        ticketTokener,
		awsCloudFrontTokener: awsCloudFrontTokener,

		db: db,
	}, nil
}
func initTicketTokener() (tokener.TicketJsonTokener, error) {
	return tokener.NewTicketJsonTokener(), nil
}
func initAwsCloudFrontTokener() (tokener.AwsCloudFrontTokener, error) {
	var nullTokener tokener.AwsCloudFrontTokener

	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		return nullTokener, err
	}

	return tokener.NewAwsCloudFrontTokener(
		tokener.AwsCloudFrontPem(pem),
		tokener.AwsCloudFrontBaseURL(os.Getenv("AWS_CLOUDFRONT_BASE_URL")),
		token.AwsCloudFrontKeyPairID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
	), nil
}
func initDB() (memory.MemoryStore, error) {
	return memory.NewMemoryStore(), nil
}

func (server *server) UserRepository() user.UserRepository {
	return server.db
}

func (server *server) UserPasswordRepository() user_password.UserPasswordRepository {
	return server.db
}

func (server *server) TicketTokener() token.TicketTokener {
	return server.ticketTokener
}

func (server *server) AwsCloudFrontTokener() token.AwsCloudFrontTokener {
	return server.awsCloudFrontTokener
}
