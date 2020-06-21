package auth_handler

import (
	"github.com/getto-systems/project-example-id/basic"
)

type (
	AppToken  []byte
	AppTicket struct {
		UserID string   `json:"user_id"`
		Roles  []string `json:"roles"`
		Token  string   `json:"token"`
	}
)

type AppIssuer struct {
	serializer AppTicketSerializer
}

type AppTicketSerializer interface {
	Serialize(basic.Ticket) (AppToken, error)
}

func (iss AppIssuer) Authorized(ticket basic.Ticket) (AppTicket, error) {
	token, err := iss.serializer.Serialize(ticket)
	if err != nil {
		return AppTicket{}, err
	}

	return AppTicket{
		UserID: string(ticket.Profile.UserID),
		Roles:  ticket.Profile.Roles,
		Token:  string(token),
	}, nil
}

func NewAppIssuer(serializer AppTicketSerializer) AppIssuer {
	return AppIssuer{
		serializer: serializer,
	}
}
