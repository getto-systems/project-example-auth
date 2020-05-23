package tokener

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/getto-systems/project-example-id/auth"
)

type JsonTokener struct {
}

type TicketTokenJson struct {
	UserID  string   `json:"user_id"`
	Roles   []string `json:"roles"`
	Expires int64    `json:"expires"`
}

type FullTokenJson struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Token  string   `json:"token"`
}

func NewJsonTokener() JsonTokener {
	return JsonTokener{}
}

func (JsonTokener) Parse(raw auth.TicketToken, path auth.Path) (*auth.Ticket, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return nil, err
	}

	var data TicketTokenJson
	err = json.NewDecoder(strings.NewReader(string(decoded))).Decode(&data)
	if err != nil {
		return nil, err
	}

	user, err := auth.NewUser(
		auth.UserID(data.UserID),
		auth.Roles(data.Roles),
		path,
	)
	if err != nil {
		return nil, err
	}

	return user.NewTicket(time.Unix(data.Expires, 0)), nil
}

func (JsonTokener) TicketToken(ticket *auth.Ticket) (auth.TicketToken, error) {
	data, err := json.Marshal(TicketTokenJson{
		UserID:  string(ticket.User().UserID()),
		Roles:   []string(ticket.User().Roles()),
		Expires: ticket.Expires().Unix(),
	})
	if err != nil {
		return auth.TicketToken(""), err
	}

	return auth.TicketToken(base64.StdEncoding.EncodeToString(data)), nil
}

func (tokener JsonTokener) FullToken(ticket *auth.Ticket) ([]byte, error) {
	token, err := tokener.FullToken(ticket)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(FullTokenJson{
		UserID: string(ticket.User().UserID()),
		Roles:  []string(ticket.User().Roles()),
		Token:  string(token),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
