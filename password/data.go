package password

type (
	LoginID string
	Login   struct {
		id LoginID
	}

	RawPassword    string
	HashedPassword []byte
)

func NewLogin(loginID LoginID) Login {
	return Login{
		id: loginID,
	}
}

func (login Login) ID() LoginID {
	return login.id
}
