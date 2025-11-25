// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"time"

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
		field.Bool("is_email_authentication_check").
			Default(false),
		field.Bool("is_line_authentication").
			Default(false),
		field.Time("updated_at").
			Default(time.Now()).
			Immutable(),
	}
}

// Edges of the RootConfig.
func (RootConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}
