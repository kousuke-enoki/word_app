import { render, screen, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { server } from '@/__tests__/mswServer'

const navigateMock = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

import LineCallback from '../LineCallback'

beforeEach(() => {
  navigateMock.mockClear()
  vi.clearAllMocks()
  localStorage.clear()
})

const setSearch = (qs: string) =>
  window.history.replaceState({}, '', `/line/callback${qs}`)

describe('LineCallback Component', () => {
  /** 1. 描画時はメッセージが出る */
  it('初期表示は「LINE 認証処理中…」だけ', () => {
    setSearch('?code=aaa&state=bbb')

    // MSWハンドラーを設定（リクエストを捕捉して警告を防ぐ）
    // このテストでは初期表示だけを確認するため、ハンドラーは設定するが結果は確認しない
    server.use(
      rest.get(
        'http://localhost:8080/users/auth/line/callback',
        (req, res, ctx) => {
          // 即座に成功レスポンスを返す（テストでは遷移を確認しない）
          return res(ctx.status(200), ctx.json({ token: 'test-token' }))
        },
      ),
    )

    render(
      <MemoryRouter>
        <LineCallback />
      </MemoryRouter>,
    )

    // 初期表示を確認（useEffect が実行される前の状態）
    expect(screen.getByText('LINE 認証処理中...'))
  })

  /** 2. コールバック成功で token 保存＆ /mypage へ */
  it('LINE 認証成功時、token 保存して /mypage へ遷移', async () => {
    setSearch('?code=sucCode&state=sucState')

    server.use(
      rest.get(
        'http://localhost:8080/users/auth/line/callback',
        (req, res, ctx) => {
          expect(req.url.searchParams.get('code')).toBe('sucCode')
          expect(req.url.searchParams.get('state')).toBe('sucState')
          return res(ctx.status(200), ctx.json({ token: 'jwt-token' }))
        },
      ),
    )

    render(
      <MemoryRouter>
        <LineCallback />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(localStorage.getItem('token')).toBe('jwt-token')
      expect(navigateMock).toHaveBeenCalledWith('/mypage')
    })
  })

  /** 3. 失敗したら /signin?err=line へ */
  it('認証エラー時は /signin?err=line へリダイレクト', async () => {
    setSearch('?code=badCode&state=badState')

    server.use(
      rest.get(
        'http://localhost:8080/users/auth/line/callback',
        (req, res, ctx) => {
          return res(ctx.status(401), ctx.json({ message: 'Unauthorized' }))
        },
      ),
    )

    render(
      <MemoryRouter>
        <LineCallback />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith('/signin?err=line')
      expect(localStorage.getItem('token')).toBeNull()
    })
  })
})
