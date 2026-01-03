// これがないと toBeInTheDocument / toHaveAttribute が未定義扱いになる
import '@testing-library/jest-dom/vitest'

import { TextEncoder } from 'node:util'

import { vi } from 'vitest'

import { server } from './mswServer'

globalThis.TextEncoder = TextEncoder as never

vi.stubGlobal('alert', vi.fn())

// vi.stubGlobal(
//   'localStorage',
//   (function () {
//     let store: Record<string, string> = {}
//     return {
//       getItem: (k: string) => (k in store ? store[k] : null),
//       setItem: (k: string, v: string) => (store[k] = v),
//       removeItem: (k: string) => delete store[k],
//       clear: () => (store = {}),
//     }
//   })(),
// )

// Object.defineProperty(global.navigator, 'clipboard', {
//   value: { writeText: vi.fn(), readText: vi.fn() },
// })

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
// 各ケースで server.use(...) を呼べば一時的にハンドラ差し替えが可能。
// 例:
//   server.use(
//     rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
//       res(ctx.status(500)),
//     ),
//   )
// afterEach で resetHandlers() が走るため、上記の差し替えはテスト単位でリセットされる。
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
