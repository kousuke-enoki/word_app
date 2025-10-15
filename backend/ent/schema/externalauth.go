// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ExternalAuth struct {
	ent.Schema
}

func (ExternalAuth) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Positive(),
		field.String("provider").
			NotEmpty(), // "line" など
		field.String("provider_user_id").
			NotEmpty(), // LINE側 sub
		field.Time("deleted_at").
			Nillable().
			Optional(),
	}
}

func (ExternalAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("external_auths").
			Field("user_id").
			Unique().
			Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}

func (ExternalAuth) Indexes() []ent.Index {
	return []ent.Index{
		// ① provider + provider_user_id は世界で一意（外部IDの衝突防止）
		index.Fields("provider", "provider_user_id").Unique(),

		// ② 同一ユーザー内で provider は一意（同じproviderの二重連携防止）
		//    Edges("user") を使うと user_id に対するインデックス/ユニークを張れる
		index.Edges("user").Fields("provider").Unique(),
	}
}
