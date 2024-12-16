import { rest } from 'msw'

export const handlers = [
  // GETリクエストのモック
  rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
    // Authorizationヘッダーを取得
    const authHeader = req.headers.get('Authorization')

    // トークンが有効な場合
    if (authHeader === 'Bearer valid-token') {
      return res(
        ctx.status(200),
        ctx.json({
          user: {
            name: 'Test User',
            admin: false,
          },
        }),
      )
    }

    // トークンが無効または存在しない場合
    return res(
      ctx.status(401),
      ctx.json({
        message: 'ログインしてください',
      }),
    )
  }),

  // OPTIONSリクエストのモック
  rest.options('http://localhost:8080/users/my_page', (req, res, ctx) => {
    return res(ctx.status(200))
  }),
]
