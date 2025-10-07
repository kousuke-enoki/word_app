package line_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" // key ID 用
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"word_app/backend/src/infrastructure/auth/line" // Provider のパッケージパスに合わせてください
)

const (
	clientID     = "client-id"
	clientSecret = "client-secret"
	redirectURI  = "https://example.com/callback"
)

// --- helper: RSA 鍵と kid 付き JWKS を生成 -----------------------------

type jwkKey struct {
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
}

type jwkSet struct {
	Keys []jwkKey `json:"keys"`
}

func newRSAKey(t *testing.T) (*rsa.PrivateKey, jwkSet, string) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// kid は公開鍵の SHA-1 フィンガープリント程度で十分
	fprint := sha1.Sum(priv.N.Bytes())
	kid := fmt.Sprintf("%x", fprint[:4])

	pub := priv.Public().(*rsa.PublicKey)
	jwk := jwkKey{
		Kty: "RSA",
		N:   b64(pub.N.Bytes()),
		E:   b64(big.NewInt(int64(pub.E)).Bytes()),
		Alg: "RS256",
		Use: "sig",
		Kid: kid,
	}
	return priv, jwkSet{Keys: []jwkKey{jwk}}, kid
}

func b64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// --- helper: ID Token を RS256 署名 ------------------------------

func signedIDToken(t *testing.T, priv *rsa.PrivateKey, kid, issuer, nonce string) string {
	t.Helper()
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   issuer, // 後で設定
		"sub":   "line-user-sub",
		"aud":   clientID,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		"nonce": nonce,
		"email": "foo@example.com",
		"name":  "Foo Bar",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	s, err := token.SignedString(priv)
	require.NoError(t, err)
	return s
}

// --------------------------- AuthURL ------------------------------

func TestProvider_AuthURL(t *testing.T) {
	p := line.NewTestProvider(
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://auth.example/authorize",
				TokenURL: "https://auth.example/token",
			},
		},
		nil, // verifier not needed for AuthURL test
	)

	state := "xyz-state"
	nonce := "abc-nonce"
	u := p.AuthURL(state, nonce)

	parsed, err := url.Parse(u)
	require.NoError(t, err)
	q := parsed.Query()

	require.Equal(t, state, q.Get("state"))
	require.Equal(t, nonce, q.Get("nonce"))
	require.Equal(t, clientID, q.Get("client_id"))
	require.Equal(t, redirectURI, q.Get("redirect_uri"))
	require.Contains(t, q.Get("scope"), "openid")
}

// ------------------------ Exchange (table) ------------------------

func TestProvider_Exchange(t *testing.T) {
	type want struct {
		ok          bool
		errContains string
	}

	tests := []struct {
		name string
		// setUp returns *line.AuthProvider, fakeServer.CloseFn
		setUp func(t *testing.T) (*line.AuthProvider, func())
		want
	}{
		{
			name: "success",
			setUp: func(t *testing.T) (*line.AuthProvider, func()) {
				priv, jwks, kid := newRSAKey(t)
				nonce := "nonce123"

				// fake OIDC discovery / jwks / token
				mux := http.NewServeMux()
				srv := httptest.NewServer(mux)
				issuerURL := srv.URL

				confResp, _ := json.Marshal(map[string]interface{}{
					"issuer":                 issuerURL,
					"jwks_uri":               issuerURL + "/keys",
					"authorization_endpoint": issuerURL + "/authorize",
					"token_endpoint":         issuerURL + "/token",
				})
				jwksResp, _ := json.Marshal(jwks)

				mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write(confResp)
				})
				mux.HandleFunc("/keys", func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write(jwksResp)
				})
				mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
					_ = r.ParseForm()
					idToken := signedIDToken(t, priv, kid, issuerURL, nonce)
					resp := map[string]interface{}{
						"access_token": "dummy",
						"id_token":     idToken,
						"token_type":   "Bearer",
						"expires_in":   3600,
					}
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(resp)
				})

				// Provider
				oidcProvider, err := oidc.NewProvider(context.Background(), issuerURL)
				require.NoError(t, err)
				p := line.NewTestProvider(
					&oauth2.Config{
						ClientID:     clientID,
						ClientSecret: clientSecret,
						RedirectURL:  redirectURI,
						Endpoint: oauth2.Endpoint{
							AuthURL:  issuerURL + "/authorize",
							TokenURL: issuerURL + "/token",
						},
					},
					oidcProvider.Verifier(&oidc.Config{ClientID: clientID}),
				)
				return p, srv.Close
			},
			want: want{ok: true},
		},
		{
			name: "token endpoint 500",
			setUp: func(_ *testing.T) (*line.AuthProvider, func()) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, "fail", http.StatusInternalServerError)
				}))
				p := line.NewTestProvider(
					&oauth2.Config{
						ClientID:     clientID,
						ClientSecret: clientSecret,
						RedirectURL:  redirectURI,
						Endpoint: oauth2.Endpoint{
							TokenURL: srv.URL,
						},
					},
					&oidc.IDTokenVerifier{}, // never used because Exchange fails first
				)
				return p, srv.Close
			},
			want: want{ok: false, errContains: "500"},
		},
		{
			name: "id_token missing",
			setUp: func(_ *testing.T) (*line.AuthProvider, func()) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					resp := map[string]interface{}{
						"access_token": "dummy",
						"token_type":   "Bearer",
					}
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(resp)
				}))
				p := line.NewTestProvider(
					&oauth2.Config{
						ClientID:     clientID,
						ClientSecret: clientSecret,
						RedirectURL:  redirectURI,
						Endpoint: oauth2.Endpoint{
							TokenURL: srv.URL,
						},
					},
					&oidc.IDTokenVerifier{},
				)
				return p, srv.Close
			},
			want: want{ok: false, errContains: "id_token"},
		},
		{
			name: "signature invalid",
			setUp: func(t *testing.T) (*line.AuthProvider, func()) {
				// 一度正常に構築 → その後 JWK に含まれない鍵で署名
				_, jwks, _ := newRSAKey(t)

				mux := http.NewServeMux()
				srv := httptest.NewServer(mux)
				issuerURL := srv.URL

				confResp, _ := json.Marshal(map[string]interface{}{
					"issuer":         issuerURL,
					"jwks_uri":       issuerURL + "/keys",
					"token_endpoint": issuerURL + "/token",
				})
				jwksResp, _ := json.Marshal(jwks)

				mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write(confResp)
				})
				mux.HandleFunc("/keys", func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write(jwksResp)
				})
				// ==== token エンドポイント側だけ *別鍵* で署名 ===
				privBad, _, kidBad := newRSAKey(t)
				mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
					idTok := signedIDToken(t, privBad, kidBad, issuerURL, "nonce")
					resp := map[string]any{
						"access_token": "x",
						"token_type":   "Bearer",
						"expires_in":   3600,
						"id_token":     idTok,
					}
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(resp)
				})

				oidcProvider, err := oidc.NewProvider(context.Background(), issuerURL)
				require.NoError(t, err)
				p := line.NewTestProvider(
					&oauth2.Config{
						ClientID:     clientID,
						ClientSecret: clientSecret,
						Endpoint: oauth2.Endpoint{
							TokenURL: issuerURL + "/token",
						},
					},
					oidcProvider.Verifier(&oidc.Config{ClientID: clientID}),
				)
				return p, srv.Close
			},
			want: want{ok: false, errContains: "verify missing"}, // Verify が失敗
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, closeFn := tt.setUp(t)
			defer closeFn()

			id, err := p.Exchange(context.Background(), "dummy-code")
			if tt.ok {
				require.NoError(t, err)
				require.NotNil(t, id)
				require.Equal(t, "line-user-sub", id.Subject)
			} else {
				require.ErrorContains(t, err, tt.errContains)
				require.Nil(t, id)
			}
		})
	}
}

// ----------------------- ValidateNonce -----------------------------

func TestProvider_ValidateNonce(t *testing.T) {
	p := line.NewTestProvider(nil, nil) // verifier 不要

	// ===== helper (前回答の buildIDToken を import しても OK) =========
	buildIDToken := func(t *testing.T, nonce string) *oidc.IDToken {
		t.Helper()
		priv, jwks, kid := newRSAKey(t)

		mux := http.NewServeMux()
		srv := httptest.NewServer(mux)
		issuer := srv.URL

		discovery, _ := json.Marshal(map[string]any{
			"issuer":   issuer,
			"jwks_uri": issuer + "/keys",
		})
		jwksJSON, _ := json.Marshal(jwks)

		mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(discovery)
		})
		mux.HandleFunc("/keys", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(jwksJSON)
		})

		raw := signedIDToken(t, priv, kid, issuer, nonce)

		prov, err := oidc.NewProvider(context.Background(), issuer)
		require.NoError(t, err)

		idTok, err := prov.Verifier(&oidc.Config{ClientID: clientID}).Verify(context.Background(), raw)
		require.NoError(t, err)

		srv.Close()
		return idTok
	}

	// ====== （オプション）invalid JSON 用の雑なトークン =============
	// makeBroken := func(raw string) *oidc.IDToken {
	// 	tok := &oidc.IDToken{}
	// 	v := reflect.ValueOf(tok).Elem()
	// 	f := v.FieldByName("raw")
	// 	if !f.IsValid() { // バージョン差分対策
	// 		for i := 0; i < v.NumField(); i++ {
	// 			if strings.Contains(strings.ToLower(v.Type().Field(i).Name), "raw") {
	// 				f = v.Field(i)
	// 				break
	// 			}
	// 		}
	// 	}
	// 	reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().SetBytes([]byte(raw))
	// 	return tok
	// }

	// ======================= テーブル ================================
	tests := []struct {
		name      string
		idTok     *oidc.IDToken
		expected  string
		wantError bool
	}{
		{"match", buildIDToken(t, "foo"), "foo", false},
		{"mismatch", buildIDToken(t, "bar"), "baz", true},
	}

	for _, tt := range tests {
		err := p.ValidateNonce(tt.idTok, tt.expected)
		if tt.wantError {
			require.Error(t, err, tt.name)
		} else {
			require.NoError(t, err, tt.name)
		}
	}
}
