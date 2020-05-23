package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/infra/repository/memory"
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

	router.HandleFunc("/renew", renewHandler(server).Handle).Methods("POST")
	router.HandleFunc("/auth/password", authPasswordHandler(server).Handle).Methods("POST")

	corsOptions := cors.New(server.cors)

	handler := corsOptions.Handler(router)

	log.Fatal(http.ListenAndServeTLS(":8080", server.tls.cert, server.tls.key, handler))
}
func renewHandler(server *server) auth_renew.Handler {
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
	cloudfrontPem, err := ioutil.ReadFile(os.Getenv("CLOUDFRONT_PEM"))
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

		ticketTokener: tokener.NewTicketJsonTokener(),
		awsCloudFrontTokener: tokener.NewAwsCloudFrontTokener(
			tokener.AwsCloudFrontPem(cloudfrontPem),
			tokener.AwsCloudFrontBaseURL(os.Getenv("CLOUDFRONT_BASE_URL")),
			token.AwsCloudFrontKeyPairID(os.Getenv("CLOUDFRONT_KEY_PAIR_ID")),
		),

		db: memory.NewMemoryStore(),
	}, nil
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
