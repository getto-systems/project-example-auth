package password_reset

type (
	TokenSender interface {
		SendToken(Destination, Token) error
	}
)
