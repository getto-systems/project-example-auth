package ticket

import (
	"errors"
	"testing"
	"time"

	"github.com/getto-systems/project-example-id/data"
)

// 適切な有効期限でチケット、nonce を発行して認証済み情報を登録
func TestIssue(t *testing.T) {
	expireSecond := data.Second(10)
	limitSecond := data.Second(100)
	nonce := Nonce("nonce")
	ticket := Ticket("ticket")

	user := data.NewUser("user-id")
	request, expires, limit := issueRequestAndExpectedExpiration(expireSecond, limitSecond)

	pub := newIssueTestEventPublisher()
	expiration := NewExpiration(ExpirationParam{
		Expires:     expireSecond,
		ExtendLimit: limitSecond,
	})

	db := newIssueTestDB()
	signer := newIssueTestSigner(ticket)
	gen := newIssueTestNonceGenerator(nonce)

	// issue
	issuer := NewIssuer(pub, db, signer, expiration, gen)
	issuedTicket, issuedNonce, issuedExpires, err := issuer.issue(request, user)

	h := newIssueTestHelper(t, pub, db, signer, issuedTicket, issuedNonce, issuedExpires, err)
	h.checkIssueError(nil)
	h.checkTicket(ticket)
	h.checkNonce(nonce)
	h.checkExpires(expires)

	h.checkSignedNonce(nonce)
	h.checkSignedUser(user)
	h.checkSignedExpires(expires)

	h.checkRegisteredTicket(issueTestTicketData{
		nonce:   nonce,
		user:    user,
		expires: expires,
		limit:   limit,
	})

	h.checkIssueTicketEvent("fired")
	h.checkIssueTicketFailedEvent(nil)
}

// チケットの署名に失敗した場合は IssueTicketFailed イベントを発行。チケット情報は登録される
func TestIssueFailedWhenSignedError(t *testing.T) {
	expireSecond := data.Second(10)
	limitSecond := data.Second(100)
	nonce := Nonce("nonce")
	ticket := Ticket("ticket")

	user := data.NewUser("user-id")
	request, expires, limit := issueRequestAndExpectedExpiration(expireSecond, limitSecond)

	pub := newIssueTestEventPublisher()
	expiration := NewExpiration(ExpirationParam{
		Expires:     expireSecond,
		ExtendLimit: limitSecond,
	})

	db := newIssueTestDB()
	signer := newIssueTestInvalidSigner() // 失敗する signer
	gen := newIssueTestNonceGenerator(nonce)

	// issue
	issuer := NewIssuer(pub, db, signer, expiration, gen)
	issuedTicket, issuedNonce, issuedExpires, err := issuer.issue(request, user)

	h := newIssueTestHelper(t, pub, db, signer, issuedTicket, issuedNonce, issuedExpires, err)
	h.checkIssueError(errors.New("ticket sign error"))
	h.checkInvalidTicket(ticket)
	h.checkInvalidNonce(nonce)
	h.checkInvalidExpires(expires)

	h.checkSignedNonce(nonce)
	h.checkSignedUser(user)
	h.checkSignedExpires(expires)

	h.checkRegisteredTicket(issueTestTicketData{
		nonce:   nonce,
		user:    user,
		expires: expires,
		limit:   limit,
	})

	h.checkIssueTicketEvent("fired")
	h.checkIssueTicketFailedEvent(errors.New("ticket sign error"))
}

// nonce の生成にした場合は IssueTicketFailed イベントを発行
func TestIssueFailedWhenGenerateNonceError(t *testing.T) {
	expireSecond := data.Second(10)
	limitSecond := data.Second(100)
	nonce := Nonce("nonce")
	ticket := Ticket("ticket")

	user := data.NewUser("user-id")
	request, expires, _ := issueRequestAndExpectedExpiration(expireSecond, limitSecond)

	pub := newIssueTestEventPublisher()
	expiration := NewExpiration(ExpirationParam{
		Expires:     expireSecond,
		ExtendLimit: limitSecond,
	})

	db := newIssueTestDB()
	signer := newIssueTestSigner(ticket)
	gen := newIssueTestInvalidNonceGenerator() // 失敗する nonce generator

	// issue
	issuer := NewIssuer(pub, db, signer, expiration, gen)
	issuedTicket, issuedNonce, issuedExpires, err := issuer.issue(request, user)

	h := newIssueTestHelper(t, pub, db, signer, issuedTicket, issuedNonce, issuedExpires, err)
	h.checkIssueError(errors.New("generate nonce error"))
	h.checkInvalidTicket(ticket)
	h.checkInvalidNonce(nonce)
	h.checkInvalidExpires(expires)

	h.checkTicketNotRegister()

	h.checkIssueTicketEvent("fired")
	h.checkIssueTicketFailedEvent(errors.New("generate nonce error"))
}

// nonce がすべて埋まっていた場合は IssueTicketFailed イベントを発行
func TestIssueFailedWhenAllGenerateNonceTryFailed(t *testing.T) {
	expireSecond := data.Second(10)
	limitSecond := data.Second(100)
	nonce := Nonce("nonce")
	ticket := Ticket("ticket")

	user := data.NewUser("user-id")
	request, expires, _ := issueRequestAndExpectedExpiration(expireSecond, limitSecond)

	pub := newIssueTestEventPublisher()
	expiration := NewExpiration(ExpirationParam{
		Expires:     expireSecond,
		ExtendLimit: limitSecond,
	})

	db := newIssueTestFilledDB() // nonce が埋まっている db
	signer := newIssueTestSigner(ticket)
	gen := newIssueTestNonceGenerator(nonce)

	// issue
	issuer := NewIssuer(pub, db, signer, expiration, gen)
	issuedTicket, issuedNonce, issuedExpires, err := issuer.issue(request, user)

	h := newIssueTestHelper(t, pub, db, signer, issuedTicket, issuedNonce, issuedExpires, err)
	h.checkIssueError(errors.New("generate nonce try failed"))
	h.checkInvalidTicket(ticket)
	h.checkInvalidNonce(nonce)
	h.checkInvalidExpires(expires)

	h.checkTicketNotRegister()

	h.checkIssueTicketEvent("fired")
	h.checkIssueTicketFailedEvent(errors.New("generate nonce try failed"))
}

// 認証済み情報の登録にした場合は IssueTicketFailed イベントを発行
func TestIssueFailedWhenRegisterError(t *testing.T) {
	expireSecond := data.Second(10)
	limitSecond := data.Second(100)
	nonce := Nonce("nonce")
	ticket := Ticket("ticket")

	user := data.NewUser("user-id")
	request, expires, _ := issueRequestAndExpectedExpiration(expireSecond, limitSecond)

	pub := newIssueTestEventPublisher()
	expiration := NewExpiration(ExpirationParam{
		Expires:     expireSecond,
		ExtendLimit: limitSecond,
	})

	db := newIssueTestInvalidDB() // register が失敗する db
	signer := newIssueTestSigner(ticket)
	gen := newIssueTestNonceGenerator(nonce)

	// issue
	issuer := NewIssuer(pub, db, signer, expiration, gen)
	issuedTicket, issuedNonce, issuedExpires, err := issuer.issue(request, user)

	h := newIssueTestHelper(t, pub, db, signer, issuedTicket, issuedNonce, issuedExpires, err)
	h.checkIssueError(errors.New("register error"))
	h.checkInvalidTicket(ticket)
	h.checkInvalidNonce(nonce)
	h.checkInvalidExpires(expires)

	h.checkIssueTicketEvent("fired")
	h.checkIssueTicketFailedEvent(errors.New("register error"))
}

type (
	issueTestEventPublisher struct {
		issueTicket           string
		issueTicketFailed     error
		authenticatedByTicket string
	}

	issueTestDB struct {
		invalid          bool
		filled           bool
		ticketRegistered bool
		ticketData       issueTestTicketData
	}

	issueTestTicketData struct {
		nonce   Nonce
		user    data.User
		expires data.Expires
		limit   data.ExtendLimit
	}

	issueTestSigner struct {
		ticket        Ticket
		signedNonce   Nonce
		signedUser    data.User
		signedExpires data.Expires
	}

	issueTestNonceGenerator struct {
		nonce Nonce
	}
	issueTestInvalidNonceGenerator struct{}

	issueTestHelper struct {
		t       *testing.T
		pub     *issueTestEventPublisher
		db      *issueTestDB
		signer  *issueTestSigner
		ticket  Ticket
		nonce   Nonce
		expires data.Expires
		err     error
	}
)

func issueRequestAndExpectedExpiration(expireSecond data.Second, limitSecond data.Second) (data.Request, data.Expires, data.ExtendLimit) {
	requestedAt := data.RequestedAt(time.Now())
	expires := requestedAt.Expires(expireSecond)
	limit := data.ExtendLimit(requestedAt.Expires(limitSecond))
	request := data.NewRequest(requestedAt, "")

	return request, expires, limit
}

func newIssueTestEventPublisher() *issueTestEventPublisher {
	return &issueTestEventPublisher{
		issueTicket: "never",
	}
}

func (pub *issueTestEventPublisher) IssueTicket(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit) {
	pub.issueTicket = "fired"
}
func (pub *issueTestEventPublisher) IssueTicketFailed(request data.Request, user data.User, expires data.Expires, limit data.ExtendLimit, err error) {
	pub.issueTicketFailed = err
}

func newIssueTestDB() *issueTestDB {
	return &issueTestDB{}
}

func newIssueTestInvalidDB() *issueTestDB {
	return &issueTestDB{
		invalid: true,
	}
}

func newIssueTestFilledDB() *issueTestDB {
	return &issueTestDB{
		filled: true,
	}
}

func (*issueTestDB) RegisterTransaction(nonce Nonce, cb func(Nonce) error) (Nonce, error) {
	err := cb(nonce)
	if err != nil {
		return Nonce(""), err
	}
	return nonce, err
}

func (db *issueTestDB) RegisterTicket(nonce Nonce, user data.User, expires data.Expires, limit data.ExtendLimit) error {
	if db.invalid {
		return errors.New("register error")
	}

	db.ticketRegistered = true
	db.ticketData = issueTestTicketData{
		nonce:   nonce,
		user:    user,
		expires: expires,
		limit:   limit,
	}
	return nil
}

func (db *issueTestDB) NonceExists(nonce Nonce) bool {
	return db.filled
}

func newIssueTestSigner(ticket Ticket) *issueTestSigner {
	return &issueTestSigner{
		ticket: ticket,
	}
}

func newIssueTestInvalidSigner() *issueTestSigner {
	return &issueTestSigner{}
}

func (*issueTestSigner) Parse(Ticket) (Nonce, data.User, data.Expires, error) {
	return Nonce(""), data.User{}, data.Expires{}, nil
}

func (signer *issueTestSigner) Sign(nonce Nonce, user data.User, expires data.Expires) (Ticket, error) {
	signer.signedNonce = nonce
	signer.signedUser = user
	signer.signedExpires = expires

	if signer.ticket == nil {
		return nil, errors.New("ticket sign error")
	}
	return signer.ticket, nil
}

func newIssueTestNonceGenerator(nonce Nonce) issueTestNonceGenerator {
	return issueTestNonceGenerator{
		nonce: nonce,
	}
}

func (gen issueTestNonceGenerator) GenerateNonce() (Nonce, error) {
	return gen.nonce, nil
}

func newIssueTestInvalidNonceGenerator() issueTestInvalidNonceGenerator {
	return issueTestInvalidNonceGenerator{}
}

func (issueTestInvalidNonceGenerator) GenerateNonce() (Nonce, error) {
	return Nonce(""), errors.New("generate nonce error")
}

func newIssueTestHelper(
	t *testing.T,
	pub *issueTestEventPublisher,
	db *issueTestDB,
	signer *issueTestSigner,
	ticket Ticket,
	nonce Nonce,
	expires data.Expires,
	err error,
) issueTestHelper {
	return issueTestHelper{
		t:       t,
		pub:     pub,
		db:      db,
		signer:  signer,
		ticket:  ticket,
		nonce:   nonce,
		expires: expires,
		err:     err,
	}
}

func (h issueTestHelper) checkTicket(ticket Ticket) {
	if string(h.ticket) != string(ticket) {
		h.t.Errorf("different ticket: %s (expected: %s)", h.ticket, ticket)
	}
}
func (h issueTestHelper) checkNonce(nonce Nonce) {
	if h.nonce != nonce {
		h.t.Errorf("different nonce: %s (expected: %s)", h.nonce, nonce)
	}
}
func (h issueTestHelper) checkExpires(expires data.Expires) {
	if h.expires != expires {
		h.t.Errorf("different expires: %v (expected: %v)", h.expires, expires)
	}
}

func (h issueTestHelper) checkInvalidTicket(ticket Ticket) {
	if string(h.ticket) == string(ticket) {
		h.t.Errorf("receive valid ticket: %v", h.ticket)
	}
}
func (h issueTestHelper) checkInvalidNonce(nonce Nonce) {
	if h.nonce == nonce {
		h.t.Errorf("receive valid nonce: %s", h.nonce)
	}
}
func (h issueTestHelper) checkInvalidExpires(expires data.Expires) {
	if h.expires == expires {
		h.t.Errorf("receive valid expires: %v", h.expires)
	}
}

func (h issueTestHelper) checkSignedNonce(nonce Nonce) {
	if h.signer.signedNonce != nonce {
		h.t.Errorf("different signed nonce: %s (expected: %s)", h.signer.signedNonce, nonce)
	}
}
func (h issueTestHelper) checkSignedUser(user data.User) {
	if h.signer.signedUser != user {
		h.t.Errorf("different signed user: %s (expected: %s)", h.signer.signedUser, user)
	}
}
func (h issueTestHelper) checkSignedExpires(expires data.Expires) {
	if h.signer.signedExpires != expires {
		h.t.Errorf("different signed expires: %v (expected: %v)", h.signer.signedExpires, expires)
	}
}

func (h issueTestHelper) checkRegisteredTicket(ticketData issueTestTicketData) {
	if h.db.ticketData != ticketData {
		h.t.Errorf("different ticket data: %v (expected: %v)", h.db.ticketData, ticketData)
	}
}
func (h issueTestHelper) checkTicketNotRegister() {
	if h.db.ticketRegistered {
		h.t.Error("ticket data registered")
	}
}

func (h issueTestHelper) checkIssueError(err error) {
	if err == nil {
		if h.err != nil {
			h.t.Errorf("issue fired: %s", h.err)
		}
	} else {
		if h.err == nil {
			h.t.Error("issue success")
		} else {
			if h.err.Error() != err.Error() {
				h.t.Errorf("issue error message is not matched: %s (expected: %s)", h.err, err)
			}
		}
	}
}
func (h issueTestHelper) checkIssueTicketEvent(event string) {
	if h.pub.issueTicket != event {
		h.t.Errorf("IssueTicket event not match: %s (expected: %s)", h.pub.issueTicket, event)
	}
}
func (h issueTestHelper) checkIssueTicketFailedEvent(err error) {
	if err == nil {
		if h.pub.issueTicketFailed != nil {
			h.t.Errorf("IssueTicketFailed event fired: %s", h.pub.issueTicketFailed)
		}
	} else {
		if h.pub.issueTicketFailed == nil {
			h.t.Error("IssueTicketFailed event not fired")
		} else {
			if h.pub.issueTicketFailed.Error() != err.Error() {
				h.t.Errorf("IssueTicketFailed error message is not matched: %s (expected: %s)", h.pub.issueTicketFailed, err)
			}
		}
	}
}
