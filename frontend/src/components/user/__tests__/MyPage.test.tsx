const navigateMock = vi.fn()

/** 必ず「他の import より前」で行う **/
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { server } from '@/__tests__/mswServer'

import MyPage from '../MyPage'

/* -------------------- 共通セットアップ -------------------- */
beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})

// 空白や改行を無視して 'Testさん' を探すマッチャ
const textEq = (expected: string) => (_: string, el?: Element | null) =>
  !!el && el.textContent?.replace(/\s+/g, '') === expected.replace(/\s+/g, '')

/* -------------------- テスト本体 -------------------- */
describe('MyPage Component', () => {
  it('通常ユーザーの場合、ユーザー名だけが表示される', async () => {
    localStorage.setItem('token', 'test-token')
    server.use(
      rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
        return res(
          ctx.status(200),
          ctx.json({
            user: { id: 1, name: 'Test User', isAdmin: false, isRoot: false },
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // ユーザー名の表示（空白/改行を無視）
    await screen.findByText(textEq('TestUserさん'))

    // Admin/Root 専用リンクがないことを確認（存在しない固定文言は使わない）
    expect(screen.queryByRole('link', { name: /単語登録/ })).toBeNull()
    expect(screen.queryByRole('link', { name: /管理設定/ })).toBeNull()

    // 任意：Userバッジの確認
    expect(screen.getByText(/👤\s*User/)).toBeInTheDocument()
  })

  it('管理ユーザーには Admin バッジと「単語登録」カードリンクが表示', async () => {
    localStorage.setItem('token', 'admin-token')
    server.use(
      rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
        return res(
          ctx.status(200),
          ctx.json({
            user: { id: 2, name: 'Admin', isAdmin: true, isRoot: false },
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    // ユーザー名（空白/改行を無視）
    await screen.findByText(textEq('Adminさん'))

    // バッジの存在確認
    expect(screen.getByText(/🔧\s*Admin/)).toBeInTheDocument()

    // カードリンクの確認（部分一致でOK）
    const adminLink = await screen.findByRole('link', { name: /単語登録/ })
    expect(adminLink).toHaveAttribute('href', '/words/new')
  })

  it('root ユーザーには Root バッジと「管理設定」カードリンクが表示', async () => {
    localStorage.setItem('token', 'root-token')
    server.use(
      rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
        return res(
          ctx.status(200),
          ctx.json({
            user: { id: 3, name: 'Root', isAdmin: false, isRoot: true },
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )
    await screen.findByText(textEq('Rootさん'))

    expect(screen.getByText(/⭐\s*Root/)).toBeInTheDocument()

    const rootLink = screen.getByRole('link', { name: /管理設定/ })
    expect(rootLink).toHaveAttribute('href', '/user/rootSetting')
  })

  it('認証エラー時に token を削除し 2 秒後にトップへリダイレクト', async () => {
    localStorage.setItem('token', 'expired-token')
    server.use(
      rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
        return res(
          ctx.status(401),
          ctx.json({
            message: 'ログインしてください',
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )
    await screen.findByText('ユーザー情報がありません。')
    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('logoutMessage')).toBe('ログインしてください')
  })
  it('サインアウトで token が消え、トップへ navigate', async () => {
    localStorage.setItem('token', 'dummy')
    server.use(
      rest.get('http://localhost:8080/users/my_page', (req, res, ctx) => {
        return res(
          ctx.status(200),
          ctx.json({
            user: { id: 1, name: 'Test', isAdmin: false, isRoot: false },
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <MyPage />
      </MemoryRouter>,
    )

    await screen.findByText(textEq('Testさん'))

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'サインアウト' }))

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/'))
    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('logoutMessage')).toBe('ログアウトしました')
  })
})
