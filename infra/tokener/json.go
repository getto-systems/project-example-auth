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

type JsonCredential struct {
	UserID  string   `json:"user_id"`
	Roles   []string `json:"roles"`
	Expires int64    `json:"expires"`
}

type JsonToken struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Token  string   `json:"token"`
}

func NewJsonTokener() JsonTokener {
	return JsonTokener{}
}

func (JsonTokener) Parse(raw auth.Credential, path auth.Path) (*auth.Ticket, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return nil, err
	}

	var data JsonCredential
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

func (JsonTokener) Credential(ticket *auth.Ticket) (auth.Credential, error) {
	data, err := json.Marshal(JsonCredential{
		UserID:  string(ticket.User().UserID()),
		Roles:   []string(ticket.User().Roles()),
		Expires: ticket.Expires().Unix(),
	})
	if err != nil {
		return auth.Credential(""), err
	}

	return auth.Credential(base64.StdEncoding.EncodeToString(data)), nil
}

func (tokener JsonTokener) Token(ticket *auth.Ticket) ([]byte, error) {
	token, err := tokener.Credential(ticket)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(JsonToken{
		UserID: string(ticket.User().UserID()),
		Roles:  []string(ticket.User().Roles()),
		Token:  string(token),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
