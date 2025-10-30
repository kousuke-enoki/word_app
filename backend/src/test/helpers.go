package test

import (
	"context"

	"word_app/backend/ent"
)

type RealEntClient struct{ *ent.Client }

func (r RealEntClient) Tx(ctx context.Context) (*ent.Tx, error)   { return r.Client.Tx(ctx) }
func (r RealEntClient) Word() *ent.WordClient                     { return r.Client.Word }
func (r RealEntClient) User() *ent.UserClient                     { return r.Client.User }
func (r RealEntClient) UserConfig() *ent.UserConfigClient         { return r.Client.UserConfig }
func (r RealEntClient) RootConfig() *ent.RootConfigClient         { return r.Client.RootConfig }
func (r RealEntClient) RegisteredWord() *ent.RegisteredWordClient { return r.Client.RegisteredWord }
func (r RealEntClient) WordInfo() *ent.WordInfoClient             { return r.Client.WordInfo }
func (r RealEntClient) JapaneseMean() *ent.JapaneseMeanClient     { return r.Client.JapaneseMean }
func (r RealEntClient) Quiz() *ent.QuizClient                     { return r.Client.Quiz }
func (r RealEntClient) QuizQuestion() *ent.QuizQuestionClient     { return r.Client.QuizQuestion }
func (r RealEntClient) ExternalAuth() *ent.ExternalAuthClient     { return r.Client.ExternalAuth }
func (r RealEntClient) UserDailyUsage() *ent.UserDailyUsageClient { return r.Client.UserDailyUsage }
func (r RealEntClient) EntClient() *ent.Client                    { return r.Client }
