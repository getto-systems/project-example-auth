package password

type (
	RawPassword    string
	HashedPassword []byte

	ChangeParam struct {
		OldPassword RawPassword
		NewPassword RawPassword
	}
)
