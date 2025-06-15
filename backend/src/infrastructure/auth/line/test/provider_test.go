// provider_test.go
package line_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"word_app/backend/config"
	"word_app/backend/src/infrastructure/auth/line"
	"word_app/backend/src/interfaces/usecase/port/auth"

	"bou.ke/monkey"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func newTestProvider() *line.Provider {
	// 実際の verify を呼ばないため clientID だけ適当で OK
	p, _ := line.NewProvider(config.LineOAuth{
		ClientID:     "cid",
		ClientSecret: "cs",
		RedirectURI:  "https://example.com/callback",
	})
	return p.(*line.Provider)
}

// ----------------------------------------------------------------------
// AuthURL
// ----------------------------------------------------------------------
func TestAuthURL(t *testing.T) {
	p := newTestProvider()
	urlStr := p.AuthURL("st", "nn")
	u, _ := url.Parse(urlStr)

	assert.Equal(t, "st", u.Query().Get("state"))
	assert.Equal(t, "nn", u.Query().Get("nonce"))
	assert.Equal(t, "cid", u.Query().Get("client_id"))
}

// ----------------------------------------------------------------------
// Exchange
// ----------------------------------------------------------------------
func TestExchange(t *testing.T) {
	defer monkey.UnpatchAll()

	okTok := oauth2.Token{AccessToken: "at"}
	okTokPtr := okTok.WithExtra(map[string]interface{}{"id_token": "dummy"})
	okIdTok := &oidc.IDToken{}

	type want struct {
		identity *auth.Identity
		err      bool
	}
	tests := []struct {
		name          string
		patchExchange func()
		patchVerify   func()
		patchClaims   func()
		want          want
	}{
		{
			name: "success",
			patchExchange: func() {
				monkey.PatchInstanceMethod(
					(*oauth2.Config)(nil), "Exchange",
					func(_ *oauth2.Config, _ context.Context, _ string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
						return okTokPtr, nil
					})
			},
			patchVerify: func() {
				monkey.PatchInstanceMethod(
					(*oidc.IDTokenVerifier)(nil), "Verify",
					func(_ *oidc.IDTokenVerifier, _ context.Context, _ string) (*oidc.IDToken, error) {
						return okIdTok, nil
					})
			},
			patchClaims: func() {
				monkey.PatchInstanceMethod(
					(*oidc.IDToken)(nil), "Claims",
					func(_ *oidc.IDToken, v interface{}) error {
						out := v.(*struct {
							Sub   string `json:"sub"`
							Email string `json:"email"`
							Name  string `json:"name"`
						})
						out.Sub = "sub1"
						out.Email = "a@b.com"
						out.Name = "Alice"
						return nil
					})
			},
			want: want{&auth.Identity{
				Provider: "line", Subject: "sub1", Email: "a@b.com", Name: "Alice",
			}, false},
		},
		{
			name: "exchange_error",
			patchExchange: func() {
				monkey.PatchInstanceMethod(
					(*oauth2.Config)(nil), "Exchange",
					func(*oauth2.Config, context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
						return nil, errors.New("exch err")
					})
			},
			want: want{nil, true},
		},
		{
			name: "id_token_missing",
			patchExchange: func() {
				monkey.PatchInstanceMethod(
					(*oauth2.Config)(nil), "Exchange",
					func(*oauth2.Config, context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
						return &oauth2.Token{AccessToken: "at"}, nil
					})
			},
			want: want{nil, true},
		},
		{
			name: "verify_error",
			patchExchange: func() {
				monkey.PatchInstanceMethod(
					(*oauth2.Config)(nil), "Exchange",
					func(*oauth2.Config, context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
						return okTokPtr, nil
					})
			},
			patchVerify: func() {
				monkey.PatchInstanceMethod(
					(*oidc.IDTokenVerifier)(nil), "Verify",
					func(*oidc.IDTokenVerifier, context.Context, string) (*oidc.IDToken, error) {
						return nil, errors.New("verify err")
					})
			},
			want: want{nil, true},
		},
		{
			name: "claims_error",
			patchExchange: func() {
				monkey.PatchInstanceMethod(
					(*oauth2.Config)(nil), "Exchange",
					func(*oauth2.Config, context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
						return okTokPtr, nil
					})
			},
			patchVerify: func() {
				monkey.PatchInstanceMethod(
					(*oidc.IDTokenVerifier)(nil), "Verify",
					func(*oidc.IDTokenVerifier, context.Context, string) (*oidc.IDToken, error) {
						return okIdTok, nil
					})
			},
			patchClaims: func() {
				monkey.PatchInstanceMethod(
					(*oidc.IDToken)(nil), "Claims",
					func(*oidc.IDToken, interface{}) error {
						return errors.New("claims err")
					})
			},
			want: want{nil, true},
		},
	}

	for _, tc := range tests {
		monkey.UnpatchAll()
		if tc.patchExchange != nil {
			tc.patchExchange()
		}
		if tc.patchVerify != nil {
			tc.patchVerify()
		}
		if tc.patchClaims != nil {
			tc.patchClaims()
		}

		p := newTestProvider()
		got, err := p.Exchange(context.Background(), "code123")

		assert.Equal(t, tc.want.identity, got, tc.name)
		if tc.want.err {
			assert.Error(t, err, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
		}
	}
}

// ----------------------------------------------------------------------
// ValidateNonce
// ----------------------------------------------------------------------
func TestValidateNonce(t *testing.T) {
	defer monkey.UnpatchAll()

	okIdTok := &oidc.IDToken{}
	badIdTok := &oidc.IDToken{}

	// Success: Claims returns expected nonce
	monkey.PatchInstanceMethod(
		(*oidc.IDToken)(nil), "Claims",
		func(tok *oidc.IDToken, v interface{}) error {
			out := v.(*struct {
				Nonce string `json:"nonce"`
			})
			if tok == okIdTok {
				out.Nonce = "good"
			} else {
				out.Nonce = "bad"
			}
			return nil
		})

	p := newTestProvider()

	assert.NoError(t, p.ValidateNonce(okIdTok, "good"), "nonce match ok")
	assert.EqualError(t, p.ValidateNonce(badIdTok, "good"), "oidc: nonce mismatch")
	// Claims error
	monkey.UnpatchInstanceMethod((*oidc.IDToken)(nil), "Claims")
	monkey.PatchInstanceMethod(
		(*oidc.IDToken)(nil), "Claims",
		func(*oidc.IDToken, interface{}) error { return errors.New("claims err") })

	assert.Error(t, p.ValidateNonce(okIdTok, "good"), "claims decode error")
}
