package static

type (
	Action interface {
		GetEnv() (Env, error)
		GetSecret(string) (Secret, error)
	}
)
