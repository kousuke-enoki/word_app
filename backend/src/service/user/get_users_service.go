package user

import (
	"context"
	"errors"

	"word_app/backend/ent"
	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/models"

	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
)

// user_list
func (s *EntUserClient) GetUsers(ctx context.Context, UserListRequest *models.UserListRequest) (*models.UserListResponse, error) {
	base := s.client.User().Query().
		Where(user.DeletedAtIsNil())
	userID := UserListRequest.UserID
	search := UserListRequest.Search
	sortBy := UserListRequest.SortBy
	order := UserListRequest.Order
	page := UserListRequest.Page
	limit := UserListRequest.Limit

	// 管理者チェック
	userEntity, err := s.client.User().Get(ctx, userID)
	if err != nil {
		logrus.Error(err)
		return nil, ErrDatabaseFailure
	}
	if !userEntity.IsAdmin {
		return nil, ErrUnauthorized
	}

	// 検索条件の追加
	base = addSearchFilter(base, search)

	// 総レコード数
	// 総件数は「未削除 + 検索条件」込みで数える
	totalCount, err := base.Clone().Count(ctx) // Clone() を使うと安全
	if err != nil {
		return nil, errors.New("failed to count users")
	}

	// Userに紐づくデータを取得 (ExternalAuthを含める)
	// query = query.WithExternalAuths()
	// ページネーション機能も考慮
	// 大文字小文字ゆらぎ対策に EqualFold を使用
	query := base.Clone().
		Offset((page - 1) * limit).Limit(limit).
		WithExternalAuths(func(q *ent.ExternalAuthQuery) {
			q.Where(externalauth.ProviderEqualFold("line"))
		})

	// ソート機能
	switch sortBy {
	case "name":
		if order == "asc" {
			query = query.Order(ent.Asc(user.FieldName))
		} else {
			query = query.Order(ent.Desc(user.FieldName))
		}
	case "email":
		if order == "asc" {
			query = query.Order(ent.Asc(user.FieldEmail))
		} else {
			query = query.Order(ent.Desc(user.FieldEmail))
		}
	case "role":
		if order == "asc" {
			query = query.Order(func(s *sql.Selector) {
				s.OrderBy(
					sql.Desc(s.C(user.FieldIsRoot)),  // root を先頭へ
					sql.Desc(s.C(user.FieldIsAdmin)), // admin を次へ
					sql.Asc(s.C(user.FieldIsTest)),   // test を最後へ
				)
			})
		} else {
			query = query.Order(func(s *sql.Selector) {
				s.OrderBy(
					sql.Asc(s.C(user.FieldIsRoot)),  // root を最後へ
					sql.Asc(s.C(user.FieldIsAdmin)), // admin を後ろへ
					sql.Desc(s.C(user.FieldIsTest)), // test を先頭へ
				)
			})
		}
	}

	// クエリ実行
	entUsers, err := query.All(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	// エンティティからレスポンス形式への変換
	users := convertEntUsersToResponse(entUsers)

	// 総ページ数を計算
	totalPages := (totalCount + limit - 1) / limit

	response := &models.UserListResponse{
		Users:      users,
		TotalPages: totalPages,
	}
	return response, nil
}

// 検索条件の追加
func addSearchFilter(query *ent.UserQuery, search string) *ent.UserQuery {
	if search != "" {
		query = query.Where(
			user.Or(
				user.NameContains(search),
				user.EmailContains(search),
			),
		)
	}
	return query
}

// エンティティからレスポンス形式に変換
func convertEntUsersToResponse(entUsers []*ent.User) []models.User {
	users := make([]models.User, 0, len(entUsers))
	for _, u := range entUsers {
		var emailPtr *string
		if u.Email != nil { // Ent も Nillable にした前提
			Email := *u.Email // string 取り出し
			emailPtr = &Email // ポインタ化（そのまま u.Email でも良い）
		}
		// password の設定有無
		isSet := false
		if u.Password != nil && *u.Password != "" { // ← Nillable の場合
			isSet = true
		}

		// LINE 連携有無：WithExternalAuths を "line" に絞っているので存在チェックのみでOK
		isLine := len(u.Edges.ExternalAuths) > 0
		createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
		updatedAt := u.UpdatedAt.Format("2006-01-02 15:04:05")

		users = append(users, models.User{
			ID:               u.ID,
			Name:             u.Name,
			IsAdmin:          u.IsAdmin,
			IsRoot:           u.IsRoot,
			IsTest:           u.IsTest,
			Email:            emailPtr,
			IsSettedPassword: isSet,
			IsLine:           isLine,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
		})
	}
	return users
}
