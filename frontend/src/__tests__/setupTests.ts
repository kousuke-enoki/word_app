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
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
