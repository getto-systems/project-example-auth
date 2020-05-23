package tokener

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/getto-systems/project-example-id/token"
	"github.com/getto-systems/project-example-id/user"
)

type TicketJsonTokener struct {
}

type TicketTokenJson struct {
	UserID  string   `json:"user_id"`
	Roles   []string `json:"roles"`
	Expires int64    `json:"expires"`
}

type TicketInfoJson struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Token  string   `json:"token"`
}

func NewTicketJsonTokener() TicketJsonTokener {
	return TicketJsonTokener{}
}

func (TicketJsonTokener) Parse(raw token.TicketToken, path user.Path) (user.Ticket, error) {
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

	user, err := user.NewUser(
		user.UserID(data.UserID),
		user.Roles(data.Roles),
		path,
	)
	if err != nil {
		return nullTicket, err
	}

	return user.NewTicket(time.Unix(data.Expires, 0)), nil
}

func (TicketJsonTokener) Token(ticket user.Ticket) (token.TicketToken, error) {
	data, err := json.Marshal(TicketTokenJson{
		UserID:  string(ticket.User().UserID()),
		Roles:   []string(ticket.User().Roles()),
		Expires: ticket.Expires().Unix(),
	})
	if err != nil {
		return nil, err
	}

	return token.TicketToken(base64.StdEncoding.EncodeToString(data)), nil
}

func (tokener TicketJsonTokener) Info(ticket user.Ticket) (token.TicketInfo, error) {
	token, err := tokener.Token(ticket)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(TicketInfoJson{
		UserID: string(ticket.User().UserID()),
		Roles:  []string(ticket.User().Roles()),
		Token:  string(token),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
