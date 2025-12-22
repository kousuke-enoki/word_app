import { rest } from 'msw'
import { setupServer } from 'msw/node'

export const server = setupServer(
  // 共通で使う GET /setting/user_config
  rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
    res(ctx.status(200), ctx.json({ is_dark_mode: false })),
  ),
  // GET /public/runtime-config - ランタイム設定取得（デフォルト）
  rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
    res(
      ctx.status(200),
      ctx.json({
        is_test_user_mode: false,
        is_line_authentication: false,
        version: '1.0.0',
      }),
    ),
  ),
  // POST /users/sign_in - サインイン（デフォルトは成功レスポンス）
  rest.post('http://localhost:8080/users/sign_in', async (req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ token: 'default-token' }))
  }),
  // POST /users/sign_up - サインアップ（デフォルトは成功レスポンス）
  rest.post('http://localhost:8080/users/sign_up', async (req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ token: 'default-token' }))
  }),
)
