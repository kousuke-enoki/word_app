// app/usecase/user/list_users.go
package user

import (
	"context"

	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/shared/ucerr"
)

// type ListUsersUsecase struct {
// 	UserRepo      repository.UserUsecase         // FindByID 用（既存）
// 	UserQueryRepo repository.UserQueryRepository // 上で定義した一覧Repo
// }

func (uc *UserUsecase) ListUsers(ctx context.Context, in ListUsersInput) (*UserListResponse, error) {
	// 1) 権限チェック（rootのみ）
	viewer, err := uc.userRepo.FindByID(ctx, in.ViewerID)
	if err != nil || viewer == nil {
		return nil, err
	}
	if !viewer.IsRoot {
		return nil, ucerr.Forbidden("forbidden")
	}

	// 2) ページング計算
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}
	offset := (in.Page - 1) * in.Limit

	// 3) Repoに委譲
	res, err := uc.userRepo.ListUsers(ctx, repository.UserListFilter{
		Search: in.Search,
		SortBy: in.SortBy,
		Order:  in.Order,
		Offset: offset,
		Limit:  in.Limit,
	})
	if err != nil {
		return nil, err
	}

	// 4) Domain -> DTO(models) 変換
	users := make([]models.User, 0, len(res.Users))
	for _, u := range res.Users {
		users = append(users, toUserDTO(u))
	}
	totalPages := (res.TotalCount + in.Limit - 1) / in.Limit
	return &UserListResponse{
		Users:      users,
		TotalPages: totalPages,
	}, nil
}

func toUserDTO(u *domain.User) models.User {
	var emailPtr *string
	if u.Email != nil {
		e := *u.Email
		emailPtr = &e
	}
	createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := u.UpdatedAt.Format("2006-01-02 15:04:05")

	return models.User{
		ID:               u.ID,
		Name:             u.Name,
		IsAdmin:          u.IsAdmin,
		IsRoot:           u.IsRoot,
		IsTest:           u.IsTest,
		Email:            emailPtr,
		IsSettedPassword: u.HasPassword,
		IsLine:           u.HasLine,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
