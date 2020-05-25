package serializer

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"

	"time"
)

type TicketJsonSerializer struct {
}

type RenewTokenJson struct {
	RenewTokenID string   `json:"renew_token_id"`
	UserID       string   `json:"user_id"`
	Roles        []string `json:"roles"`
	Authorized   int64    `json:"authorized"`
	Expires      int64    `json:"expires"`
}

type AppTokenJson struct {
	AppTokenID string   `json:"app_token_id"`
	UserID     string   `json:"user_id"`
	Roles      []string `json:"roles"`
	Authorized int64    `json:"authorized"`
	Expires    int64    `json:"expires"`
}

func NewTicketJsonSerializer() TicketJsonSerializer {
	return TicketJsonSerializer{}
}

func (TicketJsonSerializer) Parse(raw token.RenewToken, path user.Path) (user.Ticket, error) {
	var nullTicket user.Ticket

	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return nullTicket, err
	}

	var data RenewTokenJson
	err = json.NewDecoder(strings.NewReader(string(decoded))).Decode(&data)
	if err != nil {
		return nullTicket, err
	}

	if data.RenewTokenID == "" {
		return nullTicket, err
	}

	return user.RestrictTicket(path, user.TicketData{
		UserID:     user.UserID(data.UserID),
		Roles:      data.Roles,
		Authorized: time.Unix(data.Authorized, 0),
		Expires:    time.Unix(data.Expires, 0),
	})
}

func (TicketJsonSerializer) RenewToken(ticket user.Ticket) (token.RenewToken, error) {
	data, err := json.Marshal(RenewTokenJson{
		RenewTokenID: "RenewTokenID",
		UserID:       string(ticket.UserID()),
		Roles:        []string(ticket.Roles()),
		Authorized:   ticket.Authorized().Unix(),
		Expires:      ticket.Expires().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return token.RenewToken(base64.StdEncoding.EncodeToString(data)), nil
}

func (TicketJsonSerializer) AppToken(ticket user.Ticket) (token.AppToken, error) {
	var nullToken token.AppToken

	data, err := json.Marshal(AppTokenJson{
		AppTokenID: "AppTokenID",
		UserID:     string(ticket.UserID()),
		Roles:      []string(ticket.Roles()),
		Authorized: ticket.Authorized().Unix(),
		Expires:    ticket.Expires().Unix(),
	})
	if err != nil {
		return nullToken, err
	}

	return token.AppToken{
		Token:  base64.StdEncoding.EncodeToString(data),
		UserID: ticket.UserID(),
		Roles:  ticket.Roles(),
	}, nil
}
