// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserConfig struct {
	ent.Schema
}

func (UserConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Positive(),
		field.Bool("is_dark_mode").
			Default(false),
		field.Time("deleted_at").
			Nillable().
			Optional(),
	}
}

func (UserConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("user_config").
			Field("user_id").
			Unique().
			Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}
