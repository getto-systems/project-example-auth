package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/getto-systems/project-example-id/handler"
	"github.com/getto-systems/project-example-id/handler/auth/password"
	"github.com/getto-systems/project-example-id/handler/renew"
	"github.com/getto-systems/project-example-id/infra/repository/memory"
	"github.com/getto-systems/project-example-id/infra/signature/cloudfront"
	"github.com/getto-systems/project-example-id/infra/tokener"
	"github.com/getto-systems/project-example-id/signature"
)

type server struct {
	domain      handler.Domain
	allowOrigin []string

	tls tls

	db         memory.MemoryStore
	tokener    tokener.JsonTokener
	cloudfront *cloudfront.Signer
}

type tls struct {
	cert string
	key  string
}

func main() {
	config, err := initServer()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/renew", renewHandler(config).Handle).Methods("POST")
	router.HandleFunc("/auth/password", authPasswordHandler(config).Handle).Methods("POST")

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   config.allowOrigin,
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(router)

	log.Fatal(http.ListenAndServeTLS(":8080", config.tls.cert, config.tls.key, handler))
}
func renewHandler(config *server) renew.Handler {
	return renew.NewHandler(
		config.domain,
		config.tokener,
		config.cloudfront,
		config.db,
	)
}
func authPasswordHandler(config *server) password.Handler {
	return password.NewHandler(
		config.domain,
		config.tokener,
		config.cloudfront,
		config.db,
		config.db,
	)
}

func initServer() (*server, error) {
	cloudfrontPem, err := ioutil.ReadFile(os.Getenv("CLOUDFRONT_PEM"))
	if err != nil {
		return nil, err
	}

	return &server{
		domain:      handler.Domain(os.Getenv("DOMAIN")),
		allowOrigin: []string{os.Getenv("ORIGIN")},

		cloudfront: cloudfront.NewSigner(
			cloudfront.Pem(cloudfrontPem),
			cloudfront.BaseURL(os.Getenv("CLOUDFRONT_BASE_URL")),
			signature.CloudFrontKeyPairID(os.Getenv("CLOUDFRONT_KEY_PAIR_ID")),
		),

		tls: tls{
			cert: os.Getenv("TLS_CERT"),
			key:  os.Getenv("TLS_KEY"),
		},

		db:      memory.NewMemoryStore(),
		tokener: tokener.NewJsonTokener(),
	}, nil
}
