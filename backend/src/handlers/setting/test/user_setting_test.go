// setting_handler_test.go
package setting_test

import (
	"testing"
)

func TestUserSettingHandler(t *testing.T) {
	// gin.SetMode(gin.TestMode)

	// mockSvc := new(mocks.SettingClient)
	// h := setting.NewSettingHandler(mockSvc)

	// t.Run("GetUserSetting_OK", func(t *testing.T) {
	// 	cfg := &ent.UserConfig{ID: 1, IsDarkMode: true}
	// 	mockSvc.
	// 		On("GetUserConfig", mock.Anything, 99).
	// 		Return(cfg, nil)

	// 	c, w := test.NewTestCtx("GET", "/setting/user_config", nil)
	// 	test.InjectUser(c, 99, false) // 任意ユーザー
	// 	h.GetUserSettingHandler()(c)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	var got ent.UserConfig
	// 	_ = json.Unmarshal(w.Body.Bytes(), &got)
	// 	assert.True(t, got.IsDarkMode)
	// })

	// t.Run("GetUserSetting_UserNotFound", func(t *testing.T) {
	// 	mockSvc.
	// 		On("GetUserConfig", mock.Anything, 404).
	// 		Return(nil, setting.ErrUserNotFound)

	// 	c, w := test.NewTestCtx("GET", "/setting/user_config", nil)
	// 	test.InjectUser(c, 404, false)
	// 	h.GetUserSettingHandler()(c)

	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// })

	// t.Run("SaveUserSetting_OK", func(t *testing.T) {
	// 	payload := models.UserConfig{IsDarkMode: false}
	// 	updated := &ent.UserConfig{ID: 1, IsDarkMode: false}

	// 	mockSvc.
	// 		On("UpdateUserConfig", mock.Anything, 99, false).
	// 		Return(updated, nil)

	// 	c, w := test.NewTestCtx("POST", "/setting/user_config", payload)
	// 	test.InjectUser(c, 99, false)
	// 	h.SaveUserSettingHandler()(c)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// })

	// t.Run("SaveUserSetting_UserNotFound", func(t *testing.T) {
	// 	payload := models.UserConfig{IsDarkMode: true}
	// 	mockSvc.
	// 		On("UpdateUserConfig", mock.Anything, 404, true).
	// 		Return(nil, setting.ErrUserNotFound)

	// 	c, w := test.NewTestCtx("POST", "/setting/user_config", payload)
	// 	test.InjectUser(c, 404, false)
	// 	h.SaveUserSettingHandler()(c)

	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// })
	// t.Run("TestSaveUserSettingHandler", func(t *testing.T) {

	// 	payload := models.UserConfig{
	// 		IsDarkMode: false,
	// 	}
	// 	updated := &ent.UserConfig{
	// 		ID:         1,
	// 		IsDarkMode: false,
	// 	}
	// 	mockSvc.
	// 		On("UpdateUserConfig", mock.Anything, 99, "admin", true, false).
	// 		Return(updated, nil)

	// 	c, w := test.NewTestCtx("POST", "/setting/user_config", payload)
	// 	test.InjectUser(c, 99, true)
	// 	h.SaveUserSettingHandler()(c)

	// 	assert.Equal(t, http.StatusOK, w.Code)

	// })
	// // t.Run("TestSaveUserSettingHandler_validation_error", func(t *testing.T) {
	// // 	payload := models.UserConfig{
	// // 		IsDarkMode: false,
	// // 	}
	// // 	bad := payload
	// // 	bad.IsDarkMode = "hacker"

	// // 	c2, w2 := test.NewTestCtx("POST", "/setting/user_config", bad)
	// // 	test.InjectUser(c2, 99, true)
	// // 	h.SaveUserSettingHandler()(c2)
	// // 	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// // })
}
