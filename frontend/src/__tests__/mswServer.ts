import { rest } from 'msw'
import { setupServer } from 'msw/node'

export const server = setupServer(
  // 共通で使う GET /setting/user_config
  rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
    res(ctx.status(200), ctx.json({ /* 適当な設定 */ })),
  ),
  // 必要なら他のエンドポイントもここに
)