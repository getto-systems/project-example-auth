package message

type (
	LogMessage interface {
		Send(message string) error
	}
)
