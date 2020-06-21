package serializer

import (
	"log"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/basic"

	"time"
)

type TicketSerializer struct {
	jwt JWTSerializer
}

func NewTicketSerializer(jwt JWTSerializer) TicketSerializer {
	return TicketSerializer{
		jwt: jwt,
	}
}

func (serializer TicketSerializer) DecodeToken(token basic.Token) (basic.Ticket, error) {
	claims, err := serializer.jwt.Parse(string(token))
	if err != nil {
		return basic.Ticket{}, err
	}

	return basic.Ticket{
		Profile: basic.Profile{
			UserID: parseUserID(claims["sub"]),
			Roles:  parseRoles(claims["aud"]),
		},
		AuthenticatedAt: basic.AuthenticatedAt(parseTime(claims["iat"])),
		Expires:         basic.Expires(parseTime(claims["exp"])),
	}, nil
}
func parseUserID(raw interface{}) basic.UserID {
	userID, ok := raw.(string)
	if !ok {
		return basic.UserID("")
	}

	return basic.UserID(userID)
}
func parseRoles(raw interface{}) basic.Roles {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	var roles []string

	for _, val := range arr {
		role, ok := val.(string)
		if !ok {
			return nil
		}

		roles = append(roles, role)
	}

	return roles
}
func parseTime(raw interface{}) time.Time {
	unixSecond, ok := raw.(int64)
	if !ok {
		var defaultTime time.Time
		return defaultTime
	}

	return time.Unix(unixSecond, 0)
}

func (serializer TicketSerializer) Serialize(ticket basic.Ticket) (basic.Token, error) {
	token, err := serializer.jwt.Serialize(jwt.MapClaims{
		"sub": ticket.Profile.UserID,
		"aud": ticket.Profile.Roles,
		"iat": time.Time(ticket.AuthenticatedAt).Unix(),
		"exp": time.Time(ticket.Expires).Unix(),
	})
	if err != nil {
		log.Println("ticket token sign error")
		return nil, err
	}

	return basic.Token(token), nil
}
