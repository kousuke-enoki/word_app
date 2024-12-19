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
		field.Int("editing_permission").
			Default(0),
	}
}

// Edges of the RootConfig.
func (RootConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}
