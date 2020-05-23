package serializer

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type TicketJsonSerializer struct {
}

type TicketTokenJson struct {
	UserID     string   `json:"user_id"`
	Roles      []string `json:"roles"`
	Authorized int64    `json:"authorized"`
	Expires    int64    `json:"expires"`
}

type TicketInfoJson struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Token  string   `json:"token"`
}

func NewTicketJsonSerializer() TicketJsonSerializer {
	return TicketJsonSerializer{}
}

func (TicketJsonSerializer) Parse(raw token.TicketToken, path user.Path) (user.Ticket, error) {
	var nullTicket user.Ticket

	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return nullTicket, err
	}

	var data TicketTokenJson
	err = json.NewDecoder(strings.NewReader(string(decoded))).Decode(&data)
	if err != nil {
		return nullTicket, err
	}

	return user.RestrictTicket(path, user.TicketData{
		UserID:     user.UserID(data.UserID),
		Roles:      data.Roles,
		Authorized: time.Unix(data.Authorized, 0),
		Expires:    time.Unix(data.Expires, 0),
	})
}

func (TicketJsonSerializer) Token(ticket user.Ticket) (token.TicketToken, error) {
	data, err := json.Marshal(TicketTokenJson{
		UserID:     string(ticket.UserID()),
		Roles:      []string(ticket.Roles()),
		Authorized: ticket.Authorized().Unix(),
		Expires:    ticket.Expires().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return token.TicketToken(base64.StdEncoding.EncodeToString(data)), nil
}

func (serializer TicketJsonSerializer) Info(ticket user.Ticket) (token.TicketInfo, error) {
	token, err := serializer.Token(ticket)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(TicketInfoJson{
		UserID: string(ticket.UserID()),
		Roles:  []string(ticket.Roles()),
		Token:  string(token),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
