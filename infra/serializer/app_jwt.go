package serializer

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/http_handler/auth_handler"

	"github.com/getto-systems/project-example-id/data"
)

type AppSerializer struct {
	jwt JWTSerializer
}

func NewAppSerializer(jwt JWTSerializer) AppSerializer {
	return AppSerializer{
		jwt: jwt,
	}
}

func (serializer AppSerializer) Serialize(ticket data.Ticket) (auth_handler.AppToken, error) {
	token, err := serializer.jwt.Serialize(jwt.MapClaims{
		"sub": ticket.Profile.UserID,
		"aud": ticket.Profile.Roles,
		"iat": time.Time(ticket.AuthenticatedAt).Unix(),
		"exp": time.Time(ticket.Expires).Unix(),
	})
	if err != nil {
		log.Println("app token sign error")
		return auth_handler.AppToken{}, err
	}

	return auth_handler.AppToken(token), nil
}
