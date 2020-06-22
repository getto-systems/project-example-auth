package serializer

import (
	"log"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/getto-systems/project-example-id/data"

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

func (serializer TicketSerializer) DecodeToken(token data.Token) (data.Ticket, error) {
	claims, err := serializer.jwt.Parse(string(token))
	if err != nil {
		return data.Ticket{}, err
	}

	return data.Ticket{
		Profile: data.Profile{
			UserID: parseUserID(claims["sub"]),
			Roles:  parseRoles(claims["aud"]),
		},
		AuthenticatedAt: data.AuthenticatedAt(parseTime(claims["iat"])),
		Expires:         data.Expires(parseTime(claims["exp"])),
	}, nil
}
func parseUserID(raw interface{}) data.UserID {
	userID, ok := raw.(string)
	if !ok {
		return data.UserID("")
	}

	return data.UserID(userID)
}
func parseRoles(raw interface{}) data.Roles {
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
	timeString, ok := raw.(string)
	if !ok {
		var defaultTime time.Time
		return defaultTime
	}

	unix, err := strconv.Atoi(timeString)
	if err != nil {
		var defaultTime time.Time
		return defaultTime
	}

	return time.Unix(int64(unix), 0)
}

func (serializer TicketSerializer) Serialize(ticket data.Ticket) (data.Token, error) {
	token, err := serializer.jwt.Serialize(jwt.MapClaims{
		"sub": ticket.Profile.UserID,
		"aud": ticket.Profile.Roles,
		"iat": strconv.Itoa(int(time.Time(ticket.AuthenticatedAt).Unix())),
		"exp": strconv.Itoa(int(time.Time(ticket.Expires).Unix())),
	})
	if err != nil {
		log.Println("ticket token sign error")
		return nil, err
	}

	return data.Token(token), nil
}
