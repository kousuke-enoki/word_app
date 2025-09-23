package line

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"

	"word_app/backend/config"
	"word_app/backend/src/utils/tempjwt"
)

var endpoint = oauth2.Endpoint{
	AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
	TokenURL: "https://api.line.me/oauth2/v2.1/token",
}

type Provider struct {
	cfg *oauth2.Config
	// verifier *oidc.IDTokenVerifier
	secret string
}

func NewProvider(c config.LineOAuthCfg) (*Provider, error) {
	return &Provider{
		cfg: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURI,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     endpoint,
		},
	}, nil
}

func NewTestProvider(cfg *oauth2.Config, _ string) *Provider { // テスト用
	return &Provider{cfg: cfg}
}

// ------------------- AuthProvider 実装 -------------------

func (p *Provider) AuthURL(state, nonce string) string {
	return p.cfg.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
	)
}

type idClaims struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Nonce string `json:"nonce"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

func (p *Provider) Exchange(ctx context.Context, code string) (*tempjwt.Identity, error) {
	tok, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token_exchange_failed: %w", err)
	}

	raw, ok := tok.Extra("id_token").(string)
	if !ok || raw == "" {
		return nil, errors.New("id_token missing")
	}

	var cl idClaims
	parsed, err := jwt.ParseWithClaims(raw, &cl, func(t *jwt.Token) (interface{}, error) {
		// LINEはHS256（HMAC）で署名 → ClientSecret を使って検証
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected alg: %v", t.Header["alg"])
		}
		return []byte(p.cfg.ClientSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("id_token_verify_failed: %w", err)
	}

	// 主要クレームの妥当性
	if cl.Iss != "https://access.line.me" {
		return nil, fmt.Errorf("invalid iss: %s", cl.Iss)
	}
	if cl.Aud != p.cfg.ClientID {
		return nil, fmt.Errorf("invalid aud: %s", cl.Aud)
	}
	now := time.Now().Unix()
	if cl.Exp < now {
		return nil, errors.New("id_token expired")
	}

	return &tempjwt.Identity{
		Provider: "line",
		Subject:  cl.Sub,
		Email:    cl.Email,
		Name:     cl.Name,
		Nonce:    cl.Nonce,
	}, nil
}

func (p *Provider) ValidateNonce(actual, expected string) error {
	if expected == "" {
		return nil
	} // 仕様次第
	if actual != expected {
		return errors.New("oidc: nonce mismatch")
	}
	return nil
}
