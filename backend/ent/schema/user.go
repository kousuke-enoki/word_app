// Package schema contains Ent entity definitions used to generate
// the application's database models.
package schema

import (
	"errors"
	"regexp"
	"time"

	"entgo.io/ent"

	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Unique().
			Nillable().
			Optional().
			Validate(func(email string) error {
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRegex.MatchString(email) {
					return errors.New("invalid email format")
				}
				return nil
			}),
		// Go の string は値型なので 必ず何かの値（空文字含む）を持ち、「未設定」を表現できません。
		// *string（ポインタ）にすると、nil という不在状態を表現できます。
		// Ent の field.String(...).Nillable().Optional() は、DB列を NULL 許可にし、**mutator に値がセットされない限り「未設定」**として扱います。
		// このとき NotEmpty() や Validate(func(string) error) は **「値がセットされた時だけ」**評価されます（nil ならスキップ）。
		// よって、LINE だけの初期登録では email = nil のまま保存でき、**「あとからメール追加」**の時にだけ正規表現と NotEmpty が効きます。
		field.String("password").
			Sensitive().
			Nillable().
			Optional(),
		field.String("name").
			Default("JohnDoe").
			Comment("Name of the user.\n If not specified, defaults to \"John Doe\".").
			NotEmpty().
			Validate(func(name string) error {
				if len(name) < 3 || len(name) > 40 {
					return errors.New("name must be between 3 and 40 characters")
				}
				return nil
			}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Nillable().
			Optional(),
		field.Bool("isAdmin").
			Default(false),
		field.Bool("isRoot").
			Default(false),
		field.Bool("isTest").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("registered_words", RegisteredWord.Type).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("quizzes", Quiz.Type).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("user_config", UserConfig.Type).
			Unique().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("external_auths", ExternalAuth.Type).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("user_daily_usage", UserDailyUsage.Type).
			Unique().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 有効ユーザーのemailだけユニーク
		index.Fields("email").
			Annotations(entsql.IndexWhere("deleted_at IS NULL")).
			Unique(),
	}
}
