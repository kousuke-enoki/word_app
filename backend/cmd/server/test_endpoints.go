// Package main provides test endpoints for development/testing purposes.
// These endpoints should be disabled in production.
package main

import (
	"context"
	"net/http"
	"time"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/interfaces/sqlexec"
	"word_app/backend/src/usecase/apperror"

	"github.com/gin-gonic/gin"
)

// setupTestEndpoints adds test endpoints for error classification testing.
// These endpoints are for development/testing purposes only.
// TODO: Disable in production via environment variable or build tag.

// 下記のテスト用エンドポイントをたたくコマンド
// curl http://localhost:8080/test/db-error
// curl http://localhost:8080/test/timeout

func setupTestEndpoints(router *gin.Engine, runner sqlexec.Runner) {
	// DBエラーテスト: 存在しないテーブルにアクセス
	router.GET("/test/db-error", func(c *gin.Context) {
		// 存在しないテーブルにクエリを実行
		_, err := runner.ExecContext(c.Request.Context(), "SELECT * FROM non_existent_table_12345")
		if err != nil {
			httperr.Write(c, apperror.Internalf("database error", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// タイムアウトエラーテスト: コンテキストにタイムアウトを設定
	router.GET("/test/timeout", func(c *gin.Context) {
		// 100ミリ秒のタイムアウトを設定
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
		defer cancel()

		// 意図的に長時間処理をシミュレート
		time.Sleep(200 * time.Millisecond)

		// タイムアウトが発生したかチェック
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				httperr.Write(c, apperror.Internalf("timeout", ctx.Err()))
				return
			}
		default:
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}
