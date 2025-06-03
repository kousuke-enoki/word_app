package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ExternalAuth struct {
	ent.Schema
}

func (ExternalAuth) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider").
			NotEmpty(), // "line" など
		field.String("provider_user_id").
			NotEmpty(), // LINE側 sub
	}
}

func (ExternalAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("external_auths").
			Unique().
			Required(),
	}
}
