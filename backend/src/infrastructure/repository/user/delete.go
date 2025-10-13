// src/infrastructure/repository/user/ent_user_repo.go
package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
)

// type EntUserRepo struct{ c *ent.Client }

// func NewEntUserRepo(c *ent.Client) *EntUserRepo { return &EntUserRepo{c: c} }

type UserLite struct {
	ID     int
	IsTest bool
}

// DeleteIfTest は (id, is_test=true) の行を原子的に削除。
// 戻り値: deleted=trueなら削除済み。deleted=falseなら存在しないか is_test=false。
func (r *EntUserRepo) DeleteIfTest(ctx context.Context, id int) (deleted bool, err error) {
	n, err := r.client.User().
		Delete().
		Where(user.ID(id), user.IsTest(true)).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// 存在チェック（エラー詳細を返したい場合のみ使用・省略可）
func (r *EntUserRepo) Exists(ctx context.Context, id int) (bool, error) {
	exist, err := r.client.User().Query().Where(user.ID(id)).Exist(ctx)
	return exist, err
}

// is_test判定（Existsと合わせて Forbidden を出し分けしたい場合にのみ使用）
func (r *EntUserRepo) IsTest(ctx context.Context, id int) (bool, error) {
	u, err := r.client.User().Query().Where(user.ID(id)).Only(ctx)
	if ent.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return u.IsTest, nil
}
