package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pkg_error "go.cognotif/pkg/error"
	pkg_util "go.cognotif/pkg/util"
	"golang.org/x/crypto/nacl/box"
)

var (
	Token token
	b64   = base64.URLEncoding
)

const (
	token_header_key = "Authorization"
)

type token struct{}

func (token) Create(opts ...Token_Options) (t string, err error) {
	o := apply(new(Token_Option), opts...)

	{ // check token lifetime
		if o.ExpiryTime == nil || o.ExpiryTime.IsZero() {
			err = fmt.Errorf("no expiry time")
			return
		}

		if o.CreateTime == nil || o.CreateTime.IsZero() {
			err = fmt.Errorf("no create time")
			return
		}

		now := time.Now()
		if now.After(*o.ExpiryTime) {
			err = fmt.Errorf("token expired %s", o.ExpiryTime)
			return
		}

		if now.Before(*o.CreateTime) {
			err = fmt.Errorf("token predate %s", o.CreateTime)
			return
		}
	}

	p, err := json.Marshal(o)
	if err != nil {
		return
	}

	nonce := pkg_util.GenerateNonce()
	t = b64.EncodeToString(box.Seal(nonce[:], p, &nonce, &o.publicKey, &o.privateKey))

	return
}

func (token) Verify(t string, opts ...Token_Options) (oo Token_Option, err error) {
	o := apply(new(Token_Option), opts...)

	if len(t) == 0 {
		err = fmt.Errorf(pkg_error.MISSING_HEADER_AUTHORIZATION)
		return
	}

	p, err := b64.DecodeString(t)
	if err != nil {
		err = fmt.Errorf(pkg_error.TOKEN_INVALID)
		return
	}

	rest, nonce := pkg_util.ExtractNonce(p)
	p, ok := box.Open(nil, rest, &nonce, &o.publicKey, &o.privateKey)
	if !ok {
		err = fmt.Errorf(pkg_error.TOKEN_INVALID)
		return
	}

	err = json.Unmarshal(p, &oo)
	if err != nil {
		err = fmt.Errorf(pkg_error.TOKEN_INVALID)
		return
	}

	{ // check token lifetime
		if oo.ExpiryTime == nil || oo.ExpiryTime.IsZero() {
			err = fmt.Errorf(pkg_error.TOKEN_INVALID)
			return
		}

		if oo.CreateTime == nil || oo.CreateTime.IsZero() {
			err = fmt.Errorf(pkg_error.TOKEN_INVALID)
			return
		}

		now := time.Now()
		if now.After(*oo.ExpiryTime) {
			err = fmt.Errorf(pkg_error.TOKEN_EXPIRED)
			return
		}

		if now.Before(*oo.CreateTime) {
			err = fmt.Errorf(pkg_error.TOKEN_INVALID)
			return
		}
	}

	{ // check out of scopes
		var outOfScopes []string
		for _, scope := range o.Scopes {
			var found bool
			for _, v := range oo.Scopes {
				if scope == v {
					found = true
					break
				}
			}
			if !found {
				outOfScopes = append(outOfScopes, scope)
			}
		}
	}

	return
}

// CreateResponse is a http helper that will add a mapping function from the
// given *http.Response to inject as a http header
func (token) CreateResponse(opts ...Token_Options) (addToken func(*http.Response) *http.Response, err error) {
	t, err := Token.Create(opts...)
	return func(r *http.Response) *http.Response {
		if err == nil {
			r.Header.Set(token_header_key, t)
		}
		return r
	}, err
}

// VerifyRequest is a http helper that will verify the given *http.Request
// that previously injected as a http header
func (token) VerifyRequest(r *http.Request, opts ...Token_Options) (oo Token_Option, err error) {
	return Token.Verify(r.Header.Get(token_header_key), opts...)
}

func apply[T1 comparable](opt T1, opts ...func(T1)) T1 {
	for _, fn := range opts {
		if fn != nil {
			fn(opt)
		}
	}
	return opt
}

type Token_Options = func(x *Token_Option)

type Token_Option struct {
	publicKey  [32]byte `json:"-"`
	privateKey [32]byte `json:"-"`

	Scopes     []string   `json:"s,omitempty"`
	CreateTime *time.Time `json:"c,omitempty"`
	ExpiryTime *time.Time `json:"e,omitempty"`
	IDCustomer string     `json:"ic,omitempty"`
	IsAdmin    bool       `json:"ia,omitempty"`
}

func (token) With(publicKey, privateKey [32]byte, d time.Duration, scopes ...string) Token_Options {
	return func(x *Token_Option) {
		Token.WithKeypair(publicKey, privateKey)(x)
		Token.WithDuration(d)(x)
		Token.WithScopes(scopes...)
	}
}

func (token) WithScopes(scopes ...string) Token_Options {
	return func(x *Token_Option) {
		x.Scopes = scopes
	}
}

func (token) WithKeypair(publicKey, privateKey [32]byte) Token_Options {
	return func(x *Token_Option) {
		x.publicKey = publicKey
		x.privateKey = privateKey
	}
}

func (token) WithCreateTime(t time.Time) Token_Options {
	return func(x *Token_Option) {
		x.CreateTime = &t
	}
}

func (token) WithExpiryTime(t time.Time) Token_Options {
	return func(x *Token_Option) {
		x.ExpiryTime = &t
	}
}

func (token) WithDuration(d time.Duration) Token_Options {
	return func(x *Token_Option) {
		n := time.Now()
		Token.WithCreateTime(n)(x)
		Token.WithExpiryTime(n.Add(d))(x)
	}
}

func (token) WithIDCustomer(id string) Token_Options {
	return func(x *Token_Option) {
		x.IDCustomer = id
	}
}

func (token) WithIsAdmin(is_admin bool) Token_Options {
	return func(x *Token_Option) {
		x.IsAdmin = is_admin
	}
}
