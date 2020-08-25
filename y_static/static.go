package static

type (
	Env struct {
		LogLevel   string
		SecretName string
	}

	Secret struct {
		Admin      AdminSecret
		Cookie     CookieSecret
		Ticket     TicketSecret
		Api        ApiSecret
		Cloudfront CloudfrontSecret
	}
	AdminSecret struct {
		UserID   string
		LoginID  string
		Password string
	}
	CookieSecret struct {
		Domain string
	}
	TicketSecret struct {
		PrivateKey []byte
		PublicKey  []byte
	}
	ApiSecret struct {
		PrivateKey []byte
	}
	CloudfrontSecret struct {
		KeyPairID   string
		PrivateKey  []byte
		ResourceURL string
	}
)
