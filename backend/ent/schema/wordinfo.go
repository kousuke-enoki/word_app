package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// WordInfo holds the schema definition for the WordInfo entity.
type WordInfo struct {
	ent.Schema
}

// Fields of the WordInfo.
func (WordInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int("word_id").
			Positive(),
		field.Int("part_of_speech_id").
			Positive(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the WordInfo.
func (WordInfo) Edges() []ent.Edge {
	return []ent.Edge{
			edge.From("word", Word.Type).
					Ref("word_infos").
					Field("word_id").
					Unique().
					Required(),
			edge.From("part_of_speech", PartOfSpeech.Type).
					Ref("word_infos").
					Field("part_of_speech_id").
					Unique().
					Required(),
			edge.To("japanese_means", JapaneseMean.Type),
		}
}
