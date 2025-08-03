// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"errors"
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// JapaneseMean holds the schema definition for the JapaneseMean entity.
type JapaneseMean struct {
	ent.Schema
}

// Fields of the JapaneseMean.
func (JapaneseMean) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Validate(func(s string) error {
				// アルファベット以外の文字列のみ許可（ひらがな、カタカナ、漢字など）
				match, _ := regexp.MatchString(`^[^\x00-\x7F]+$`, s)
				if !match {
					return errors.New("name must contain only non-alphabetic characters (e.g., Hiragana, Katakana, Kanji)")
				}
				return nil
			}),
		field.Int("word_info_id").
			Positive(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the JapaneseMean.
func (JapaneseMean) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("word_info", WordInfo.Type).
			Ref("japanese_means").
			Field("word_info_id").
			Unique().
			Required(),
		edge.To("quiz_questions", QuizQuestion.Type),
	}
}

func (JapaneseMean) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("word_info_id", "name").
			Unique(),
	}
}
