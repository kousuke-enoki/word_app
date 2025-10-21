package contextutil

// type UserRoles struct {
// 	UserID  int
// 	IsAdmin bool
// 	IsRoot  bool
// 	IsTest  bool
// }

// errorを返すのは、userIDがcontextに存在しない場合のみ
// つまり未ログイン状態で呼び出された場合なので、基本的にはログイン必須のAPIでしか使わない
// その場合はtopに戻す
// func GetUserRoles(c *gin.Context) (*UserRoles, error) {
// 	principal, ok := jwt.GetPrincipal(c)
// 	if !ok {
// 		return nil, errors.New("userID not found in context")
// 	}

// 	return &UserRoles{
// 		UserID:  principal.UserID,
// 		IsAdmin: principal.IsAdmin,
// 		IsRoot:  principal.IsRoot,
// 		IsTest:  principal.IsTest,
// 	}, nil
// }
