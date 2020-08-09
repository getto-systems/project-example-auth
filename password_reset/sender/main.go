package sender

import (
	"errors"
	"fmt"

	password_reset_infra "github.com/getto-systems/project-example-id/infra/password_reset"

	"github.com/getto-systems/project-example-id/data/password_reset"
	"github.com/getto-systems/project-example-id/message"
)

type (
	TokenSender struct {
		log message.LogMessage
	}
)

func NewTokenSender(log message.LogMessage) TokenSender {
	return TokenSender{
		log: log,
	}
}

func (sender TokenSender) tokenSender() password_reset_infra.TokenSender {
	return sender
}

func (sender TokenSender) SendToken(dest password_reset.Destination, token password_reset.Token) (err error) {
	switch dest.Type {
	case "Log":
		return sender.log.Send(fmt.Sprintf("password reset token: %s", token))
	}

	return errors.New("unknown destination")
}
