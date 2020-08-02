package password_reset

type (
	SessionGenerator interface {
		GenerateSession() (SessionID, Token, error)
	}
)
