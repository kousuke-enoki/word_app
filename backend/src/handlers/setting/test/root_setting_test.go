// setting_handler_test.go
package setting_test

import (
	"testing"
)

func TestGetRootSettingHandler(t *testing.T) {
	// 	t.Run("TestGetRootSettingHandler", func(t *testing.T) {
	// 		gin.SetMode(gin.TestMode)

	// 		mockSvc := new(mocks.SettingClient)
	// 		h := setting.NewSettingHandler(mockSvc)
	// 		rootCfg := &ent.RootConfig{
	// 			ID:                         1,
	// 			EditingPermission:          "admin",
	// 			IsTestUserMode:             false,
	// 			IsEmailAuthenticationCheck: false,
	// 			IsLineAuthentication:       false,
	// 		}
	// 		mockSvc.
	// 			On("GetRootConfig", mock.Anything, 99).
	// 			Return(rootCfg, nil)

	// 		c, w := test.NewTestCtx("GET", "/setting/root_config", nil)
	// 		test.InjectUser(c, 99, true) // root ユーザー
	// 		h.GetRootSettingHandler()(c)

	// 		assert.Equal(t, http.StatusOK, w.Code)
	// 		var got ent.RootConfig
	// 		_ = json.Unmarshal(w.Body.Bytes(), &got)
	// 		assert.Equal(t, rootCfg.EditingPermission, got.EditingPermission)
	// 	})

	// 	t.Run("TestGetRootSettingHandler_role_error", func(t *testing.T) {
	// 		gin.SetMode(gin.TestMode)

	// 		mockSvc := new(mocks.SettingClient)
	// 		h := setting.NewSettingHandler(mockSvc)
	// 		c2, w2 := test.NewTestCtx("GET", "/setting/root_config", nil)
	// 		test.InjectUser(c2, 50, false) // root ではない
	// 		h.GetRootSettingHandler()(c2)
	// 		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	// 	})
	// 	gin.SetMode(gin.TestMode)

	// 	mockSvc := new(mocks.SettingClient)
	// 	h := setting.NewSettingHandler(mockSvc)

	// 	t.Run("TestSaveRootSettingHandler", func(t *testing.T) {

	// 		payload := models.RootConfig{
	// 			EditingPermission: "admin",
	// 			IsTestUserMode:    true,
	// 			IsEmailAuthCheck:  false,
	// 			IsLineAuth:        false,
	// 		}
	// 		updated := &ent.RootConfig{
	// 			ID:                         1,
	// 			EditingPermission:          "admin",
	// 			IsTestUserMode:             true,
	// 			IsEmailAuthenticationCheck: false,
	// 			IsLineAuthentication:       false,
	// 		}
	// 		mockSvc.
	// 			On("UpdateRootConfig", mock.Anything, 99, "admin", true, false, false).
	// 			Return(updated, nil)
	// 		c, w := test.NewTestCtx("POST", "/setting/root_config", payload)
	// 		test.InjectUser(c, 99, true)
	// 		h.SaveRootSettingHandler()(c)

	// 		assert.Equal(t, http.StatusOK, w.Code)

	// 	})
	// 	t.Run("TestSaveRootSettingHandler_validation_error", func(t *testing.T) {
	// 		payload := models.RootConfig{
	// 			EditingPermission: "admin",
	// 			IsTestUserMode:    true,
	// 			IsEmailAuthCheck:  false,
	// 			IsLineAuth:        false,
	// 		}
	// 		bad := payload
	// 		bad.EditingPermission = "hacker"

	// 		c2, w2 := test.NewTestCtx("POST", "/setting/root_config", bad)
	// 		test.InjectUser(c2, 99, true)
	// 		h.SaveRootSettingHandler()(c2)
	// 		assert.Equal(t, http.StatusBadRequest, w2.Code)

	// })
	//
	//	t.Run("TestSaveRootSettingHandler_role_error", func(t *testing.T) {
	//		payload := models.RootConfig{
	//			EditingPermission: "admin",
	//			IsTestUserMode:    true,
	//			IsEmailAuthCheck:  false,
	//			IsLineAuth:        false,
	//		}
	//		c3, w3 := test.NewTestCtx("POST", "/setting/root_config", payload)
	//		test.InjectUser(c3, 10, false) // root でない
	//		h.SaveRootSettingHandler()(c3)
	//		assert.Equal(t, http.StatusUnauthorized, w3.Code)
	//	})
}
