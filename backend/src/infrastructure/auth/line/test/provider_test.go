package line_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"word_app/backend/src/infrastructure/auth/line"
)

const (
	clientID     = "client-id"
	clientSecret = "client-secret"
	redirectURI  = "https://example.com/callback"
)

// ============ HS256 の id_token を作るヘルパ ============

func hsIDToken(t *testing.T, secret, issuer, nonce string) string {
	t.Helper()
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   issuer,
		"sub":   "line-user-sub",
		"aud":   clientID,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		"nonce": nonce,
		"email": "foo@example.com",
		"name":  "Foo Bar",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

// --------------------------- AuthURL ------------------------------

func TestProvider_AuthURL(t *testing.T) {
	// テスト用 NewTestProvider（HS256版）を使う
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
		clientSecret, // ← HS256 の secret
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

// ------------------------ Exchange（HS256） ------------------------

func TestProvider_Exchange(t *testing.T) {
	// 	type want struct {
	// 		ok          bool
	// 		errContains string
	// 	}

	// 	tests := []struct {
	// 		name  string
	// 		setUp func(t *testing.T) (*line.Provider, func())
	// 		want
	// 	}{
	// 		{
	// 			name: "success",
	// 			setUp: func(t *testing.T) (*line.Provider, func()) {
	// 				nonce := "nonce123"

	// 				mux := http.NewServeMux()
	// 				srv := httptest.NewServer(mux)

	// 				mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
	// 					_ = r.ParseForm()
	// 					// Provider 側は iss を "https://access.line.me" とチェックしている実装なので合わせる
	// 					idToken := hsIDToken(t, clientSecret, "https://access.line.me", nonce)
	// 					resp := map[string]any{
	// 						"access_token": "dummy",
	// 						"id_token":     idToken,
	// 						"token_type":   "Bearer",
	// 						"expires_in":   3600,
	// 						"scope":        "openid profile",
	// 					}
	// 					w.Header().Set("Content-Type", "application/json")
	// 					_ = json.NewEncoder(w).Encode(resp)
	// 				})

	// 				p := line.NewTestProvider(
	// 					&oauth2.Config{
	// 						ClientID:     clientID,
	// 						ClientSecret: clientSecret,
	// 						RedirectURL:  redirectURI,
	// 						Endpoint: oauth2.Endpoint{
	// 							TokenURL: srv.URL + "/token",
	// 						},
	// 					},
	// 					clientSecret,
	// 				)
	// 				return p, srv.Close
	// 			},
	// 			want: want{ok: true},
	// 		},
	// 		{
	// 			name: "token endpoint 500",
	// 			setUp: func(t *testing.T) (*line.Provider, func()) {
	// 				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	// 					http.Error(w, "fail", http.StatusInternalServerError)
	// 				}))
	// 				p := line.NewTestProvider(
	// 					&oauth2.Config{
	// 						ClientID:     clientID,
	// 						ClientSecret: clientSecret,
	// 						RedirectURL:  redirectURI,
	// 						Endpoint: oauth2.Endpoint{
	// 							TokenURL: srv.URL,
	// 						},
	// 					},
	// 					clientSecret,
	// 				)
	// 				return p, srv.Close
	// 			},
	// 			// 実装は "token_exchange_failed: ..." を返すようにしているのでそちらを期待
	// 			want: want{ok: false, errContains: "token_exchange_failed"},
	// 		},
	// 		{
	// 			name: "id_token missing",
	// 			setUp: func(t *testing.T) (*line.Provider, func()) {
	// 				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	// 					resp := map[string]any{
	// 						"access_token": "dummy",
	// 						"token_type":   "Bearer",
	// 					}
	// 					w.Header().Set("Content-Type", "application/json")
	// 					_ = json.NewEncoder(w).Encode(resp)
	// 				}))
	// 				p := line.NewTestProvider(
	// 					&oauth2.Config{
	// 						ClientID:     clientID,
	// 						ClientSecret: clientSecret,
	// 						RedirectURL:  redirectURI,
	// 						Endpoint: oauth2.Endpoint{
	// 							TokenURL: srv.URL,
	// 						},
	// 					},
	// 					clientSecret,
	// 				)
	// 				return p, srv.Close
	// 			},
	// 			want: want{ok: false, errContains: "id_token"},
	// 		},
	// 		{
	// 			name: "signature invalid (wrong secret)",
	// 			setUp: func(t *testing.T) (*line.Provider, func()) {
	// 				mux := http.NewServeMux()
	// 				srv := httptest.NewServer(mux)

	// 				mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
	// 					// 誤った secret で署名 → Provider 側検証は失敗する
	// 					idTok := hsIDToken(t, "bad-secret", "https://access.line.me", "n")
	// 					resp := map[string]any{
	// 						"access_token": "x",
	// 						"token_type":   "Bearer",
	// 						"expires_in":   3600,
	// 						"id_token":     idTok,
	// 					}
	// 					w.Header().Set("Content-Type", "application/json")
	// 					_ = json.NewEncoder(w).Encode(resp)
	// 				})

	// 				p := line.NewTestProvider(
	// 					&oauth2.Config{
	// 						ClientID:     clientID,
	// 						ClientSecret: clientSecret,
	// 						Endpoint: oauth2.Endpoint{
	// 							TokenURL: srv.URL + "/token",
	// 						},
	// 					},
	// 					clientSecret, // 正しい secret（＝署名と不一致）
	// 				)
	// 				return p, srv.Close
	// 			},
	// 			// 実装は "id_token_verify_failed: ..." を返している
	// 			want: want{ok: false, errContains: "id_token_verify_failed"},
	// 		},
	// 	}

	// 	for _, tt := range tests {
	// 		t.Run(tt.name, func(t *testing.T) {
	// 			p, closeFn := tt.setUp(t)
	// 			defer closeFn()

	// 			id, err := p.Exchange(context.Background(), "dummy-code")
	// 			if tt.ok {
	// 				require.NoError(t, err)
	// 				require.NotNil(t, id)
	// 				require.Equal(t, "line-user-sub", id.Subject)
	// 			} else {
	// 				require.ErrorContains(t, err, tt.errContains)
	// 				require.Nil(t, id)
	// 			}
	// 		})
	// 	}
	// }

	// // ----------------------- ValidateNonce -----------------------------
	// // ※ Provider に ValidateNonce(*oidc.IDToken, expected string) を残している前提。
	// //   これは go-oidc で作った ID Token を使って “nonce claims の一致” だけを検証するユーティリティ扱いです。
	// //   Exchange 経路は HS256 に寄せましたが、このテストはそのままでも動きます。

	// // ========= 最低限の RS256/JWKS モック（ValidateNonce用） =========

	// func b64raw(b []byte) string {
	//   return base64.RawURLEncoding.EncodeToString(b)
	// }

	// type jwkKey struct {
	//   Kty string `json:"kty"`
	//   E   string `json:"e"`
	//   N   string `json:"n"`
	//   Alg string `json:"alg"`
	//   Use string `json:"use"`
	//   Kid string `json:"kid"`
	// }
	// type jwkSet struct{ Keys []jwkKey `json:"keys"` }

	// func TestProvider_ValidateNonce(t *testing.T) {
	// 	p := line.NewTestProvider(nil, clientSecret) // ValidateNonce は verifier 不要

	// 	// RSA鍵・JWKS・discovery を立て、oidc で IDToken を一度検証してから渡す
	// 	buildIDToken := func(t *testing.T, nonce string) *oidc.IDToken {
	// 		t.Helper()

	// 		// 1) 署名用 RSA 鍵
	// 		priv, err := rsa.GenerateKey(rand.Reader, 2048)
	// 		require.NoError(t, err)

	// 		// kid は公開鍵のフィンガープリント等でOK
	// 		fp := sha1.Sum(priv.N.Bytes())
	// 		kid := fmt.Sprintf("%x", fp[:8])

	// 		// 2) JWKS に公開鍵を載せる
	// 		pub := &priv.PublicKey
	// 		jwks := jwkSet{Keys: []jwkKey{{
	// 			Kty: "RSA",
	// 			Alg: "RS256",
	// 			Use: "sig",
	// 			Kid: kid,
	// 			N:   b64raw(pub.N.Bytes()),
	// 			E:   b64raw(big.NewInt(int64(pub.E)).Bytes()),
	// 		}}}

	// 		// 3) OIDC Discovery/JWKS のモック
	// 		mux := http.NewServeMux()
	// 		srv := httptest.NewServer(mux)
	// 		issuer := srv.URL

	// 		discovery, _ := json.Marshal(map[string]any{
	// 			"issuer":   issuer,
	// 			"jwks_uri": issuer + "/keys",
	// 		})
	// 		jwksJSON, _ := json.Marshal(jwks)

	// 		mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
	// 			_, _ = w.Write(discovery)
	// 		})
	// 		mux.HandleFunc("/keys", func(w http.ResponseWriter, _ *http.Request) {
	// 			_, _ = w.Write(jwksJSON)
	// 		})

	// 		// 4) RS256 で署名して kid を付ける
	// 		now := time.Now()
	// 		claims := jwt.MapClaims{
	// 			"iss":   issuer,
	// 			"sub":   "sub",
	// 			"aud":   clientID,
	// 			"exp":   now.Add(time.Hour).Unix(),
	// 			"iat":   now.Unix(),
	// 			"nonce": nonce,
	// 		}
	// 		tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 		tok.Header["kid"] = kid
	// 		raw, err := tok.SignedString(priv)
	// 		require.NoError(t, err)

	// 		// 5) go-oidc で検証 → *oidc.IDToken を返す
	// 		prov, err := oidc.NewProvider(context.Background(), issuer)
	// 		require.NoError(t, err)

	// 		idTok, err := prov.Verifier(&oidc.Config{ClientID: clientID}).Verify(context.Background(), raw)
	// 		require.NoError(t, err)

	// 		srv.Close()
	// 		return idTok
	// 	}

	// 	tests := []struct {
	// 		name      string
	// 		idTok     *oidc.IDToken
	// 		expected  string
	// 		wantError bool
	// 	}{
	// 		{"match", buildIDToken(t, "foo"), "foo", false},
	// 		{"mismatch", buildIDToken(t, "bar"), "baz", true},
	// 	}

	//	for _, tt := range tests {
	//		err := p.ValidateNonce(tt.idTok, tt.expected)
	//		if tt.wantError {
	//			require.Error(t, err, tt.name)
	//		} else {
	//			require.NoError(t, err, tt.name)
	//		}
	//	}
}
