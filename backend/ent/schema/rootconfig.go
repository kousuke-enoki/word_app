package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// RootConfig holds the schema definition for the RootConfig entity.
type RootConfig struct {
	ent.Schema
}

// Fields of the RootConfig.
func (RootConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("editing_permission").
			Default("admin").
			NotEmpty(),
		field.Bool("is_test_user_mode").
			Default(false),
		field.Bool("is_email_authentication").
			Default(false),
	}
}

// Edges of the RootConfig.
func (RootConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}
