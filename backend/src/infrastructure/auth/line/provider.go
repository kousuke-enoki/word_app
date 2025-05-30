package line

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"word_app/backend/config"
	auth_port "word_app/backend/src/usecase/port/auth"
)

var endpoint = oauth2.Endpoint{
	AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
	TokenURL: "https://api.line.me/oauth2/v2.1/token",
}

type Provider struct {
	cfg      *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewProvider(c config.LineOAuth) (auth_port.AuthProvider, error) {
	oidcProvider, err := oidc.NewProvider(context.Background(), "https://access.line.me")
	if err != nil {
		return nil, err
	}

	return &Provider{
		cfg: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.RedirectURI,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     endpoint,
		},
		verifier: oidcProvider.Verifier(&oidc.Config{
			ClientID: c.ClientID,
		}),
	}, nil
}

// ------------------- AuthProvider 実装 -------------------

func (p *Provider) AuthURL(state, nonce string) string {
	return p.cfg.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
	)
}

func (p *Provider) Exchange(ctx context.Context, code string) (*auth_port.Identity, error) {
	tok, err := p.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	rawID, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, err
	}

	idTok, err := p.verifier.Verify(ctx, rawID)
	if err != nil {
		return nil, err
	}

	var cl struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idTok.Claims(&cl); err != nil {
		return nil, err
	}

	return &auth_port.Identity{
		Provider: "line",
		Subject:  cl.Sub,
		Email:    cl.Email,
		Name:     cl.Name,
	}, nil
}

func (p *Provider) ValidateNonce(idTok *oidc.IDToken, expected string) error {
	var cl struct {
		Nonce string `json:"nonce"`
	}
	if err := idTok.Claims(&cl); err != nil {
		return err
	}
	if cl.Nonce != expected {
		return errors.New("oidc: nonce mismatch")
	}
	return nil
}
