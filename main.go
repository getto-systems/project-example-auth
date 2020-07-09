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

	"github.com/getto-systems/project-example-id/authenticate"

	"github.com/getto-systems/project-example-id/user"
	"github.com/getto-systems/project-example-id/user/subscriber"
)

type Server struct {
	port string

	cors cors.Options
	tls  Tls

	logger logger.Logger

	authResponse auth_handler.AuthResponse
	userFactory  UserFactory
}

type Tls struct {
	cert string
	key  string
}

type UserFactory struct {
	authenticated user.UserAuthenticatedFactory
	ticketAuth    user.UserTicketAuthFactory
	passwordAuth  user.UserPasswordAuthFactory
}

func main() {
	log.Fatal(NewServer().listen())
}
func (server Server) listen() error {
	router := server.routes()
	handler := cors.New(server.cors).Handler(router)

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
func (server Server) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/healthz", healthz).Methods("GET")

	router.HandleFunc("/auth/ticket", server.AuthByTicket().Handle).Methods("POST")
	router.HandleFunc("/auth/password", server.AuthByPassword().Handle).Methods("POST")

	return router
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := "\"OK\""

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func (server Server) AuthByTicket() auth_handler.AuthByTicketHandler {
	auth := authenticate.NewAuthByTicket(server.userFactory.authenticated, server.userFactory.ticketAuth)
	return auth_handler.NewAuthByTicketHandler(server.logger, server.authResponse, auth)
}
func (server Server) AuthByPassword() auth_handler.AuthByPasswordHandler {
	auth := authenticate.NewAuthByPassword(server.userFactory.authenticated, server.userFactory.passwordAuth)
	return auth_handler.NewAuthByPasswordHandler(server.logger, server.authResponse, auth)
}

func NewServer() Server {
	appLogger := NewAppLogger()

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

		authResponse: NewAuthResponse(),
		userFactory:  NewUserFactory(appLogger),
	}
}

func NewAppLogger() logger.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"), log.New(os.Stdout, "", 0))
}

func NewAuthResponse() auth_handler.AuthResponse {
	cookieDomain := auth_handler.CookieDomain(os.Getenv("COOKIE_DOMAIN"))
	return auth_handler.NewAuthResponse(cookieDomain, NewAwsCloudFrontIssuer(), NewAppIssuer())
}
func NewAwsCloudFrontIssuer() auth_handler.AwsCloudFrontIssuer {
	pem, err := ioutil.ReadFile(os.Getenv("AWS_CLOUDFRONT_PEM"))
	if err != nil {
		log.Fatalf("aws cloudfront private key read failed: %s", err)
	}

	serializer := serializer.NewAwsCloudFrontSerializer(
		pem,
		os.Getenv("AWS_CLOUDFRONT_SECURE_URL"),
	)

	return auth_handler.NewAwsCloudFrontIssuer(
		auth_handler.AwsCloudFrontKeyPairID(os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")),
		serializer,
	)
}

func NewAppIssuer() auth_handler.AppIssuer {
	pem, err := ioutil.ReadFile(os.Getenv("APP_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("app private key read failed: %s", err)
	}

	key, err := serializer.NewJWT_ES_512_Key(serializer.JWT_Pem{
		PrivateKey: pem,
	})
	if err != nil {
		log.Fatalf("app key parse failed: %s", err)
	}

	jwt := serializer.NewJWTSerializer(key)
	app := serializer.NewAppSerializer(jwt)

	return auth_handler.NewAppIssuer(app)
}

func NewUserFactory(appLogger logger.Logger) UserFactory {
	db := NewDB()
	pubsub := NewPubSub()
	policy := NewPolicy()

	userLogger := subscriber.NewUserLogger(appLogger)
	pubsub.SubscribeAuthenticated(userLogger)
	pubsub.SubscribeTicketAuth(userLogger)
	pubsub.SubscribePasswordAuth(userLogger)

	ticketSerializer := NewTicketSerializer()
	passwordEncrypter := NewPasswordEncrypter()

	return UserFactory{
		authenticated: user.NewUserAuthenticatedFactory(pubsub, db, policy, ticketSerializer),
		ticketAuth:    user.NewUserTicketAuthFactory(pubsub, ticketSerializer),
		passwordAuth:  user.NewUserPasswordAuthFactory(pubsub, db, passwordEncrypter),
	}
}
func NewDB() *db.MemoryStore {
	return db.NewMemoryStore()
}
func NewPubSub() *pubsub.SyncPubSub {
	return pubsub.NewSyncPubSub()
}
func NewPolicy() policy.Policy {
	return policy.NewPolicy()
}
func NewTicketSerializer() serializer.TicketSerializer {
	privateKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("ticket private key read failed: %s", err)
	}

	publicKeyPem, err := ioutil.ReadFile(os.Getenv("TICKET_PUBLIC_KEY"))
	if err != nil {
		log.Fatalf("ticket public key read failed: %s", err)
	}

	key, err := serializer.NewJWT_ES_512_Key(serializer.JWT_Pem{
		PrivateKey: privateKeyPem,
		PublicKey:  publicKeyPem,
	})
	if err != nil {
		log.Fatalf("ticket key parse failed: %s", err)
	}

	jwt := serializer.NewJWTSerializer(key)
	return serializer.NewTicketSerializer(jwt)
}
func NewPasswordEncrypter() password.PasswordEncrypter {
	return password.NewPasswordEncrypter(10) // bcrypt.DefaultCost
}
