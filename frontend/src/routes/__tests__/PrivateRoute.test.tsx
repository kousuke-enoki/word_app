/**
 * PrivateRoute.test.tsx
 *
 * - Vitest + @testing-library/react
 * - useAuth をケースごとにモックして挙動を検証する
 */
import { render, screen } from '@testing-library/react'
import React from 'react'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import PrivateRoute from '../PrivateRoute'

/* ------------------------------------------------------------------ */
/* ❶ useAuth を都度好きな値で返すようモック                            */
/* ------------------------------------------------------------------ */
vi.mock('@/hooks/useAuth', () => ({
  useAuth: vi.fn(),
}))
import { useAuth } from '@/hooks/useAuth' // ← 上でモック化済み

const mockedUseAuth = vi.mocked(useAuth) // 型付きのモック関数になる

const setAuthState = (state: Partial<ReturnType<typeof useAuth>>) => {
  mockedUseAuth.mockReturnValue({
    isLoggedIn: false,
    userRole: 'test',
    isLoading: false,
    ...state,
  })
}

/* ------------------------------------------------------------------ */
/* ❷ 現在パスを確認できるダミーコンポーネント                           */
/* ------------------------------------------------------------------ */
const WhereAmI = () => {
  const loc = useLocation()
  return <p data-testid="path">{loc.pathname}</p>
}

/* ------------------------------------------------------------------ */
/* ❸ レンダーヘルパ：<MemoryRouter> 内にルーティングを組み立てる        */
/* ------------------------------------------------------------------ */
const renderWithRouter = (ui: React.ReactElement, startPath = '/private') =>
  render(
    <MemoryRouter initialEntries={[startPath]}>
      <Routes>
        <Route path="/" element={<WhereAmI />} />
        <Route path="/private" element={ui} />
      </Routes>
    </MemoryRouter>,
  )

/* ------------------------------------------------------------------ */
/*                               TESTS                                */
/* ------------------------------------------------------------------ */
describe('PrivateRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('ロード中は Loading... を表示', () => {
    setAuthState({ isLoading: true })

    renderWithRouter(
      <PrivateRoute>
        <p>dummy</p>
      </PrivateRoute>,
    )

    expect(screen.getByText('Loading...'))
  })

  it('未ログインは / へリダイレクト', () => {
    setAuthState({ isLoggedIn: false })

    renderWithRouter(
      <PrivateRoute>
        <p>secret</p>
      </PrivateRoute>,
    )

    expect(screen.getByTestId('path').textContent).toBe('/') // ルートに飛ばされた
  })

  it('ログイン済み・権限制限なしなら子要素を表示', () => {
    setAuthState({ isLoggedIn: true, userRole: 'general' })

    renderWithRouter(
      <PrivateRoute>
        <p>secret</p>
      </PrivateRoute>,
    )

    expect(screen.getByText('secret'))
  })

  it('一般ユーザーが admin 専用ページにアクセスすると / へリダイレクト', () => {
    setAuthState({ isLoggedIn: true, userRole: 'general' })

    renderWithRouter(
      <PrivateRoute requiredRole="admin">
        <p>admin page</p>
      </PrivateRoute>,
    )

    expect(screen.getByTestId('path').textContent).toBe('/')
  })

  it('admin ユーザーは requiredRole="admin" を通過できる', () => {
    setAuthState({ isLoggedIn: true, userRole: 'admin' })

    renderWithRouter(
      <PrivateRoute requiredRole="admin">
        <p>admin page</p>
      </PrivateRoute>,
    )

    expect(screen.getByText('admin page'))
  })

  it('root ユーザーは requiredRole="admin" も通過できる', () => {
    setAuthState({ isLoggedIn: true, userRole: 'root' })

    renderWithRouter(
      <PrivateRoute requiredRole="admin">
        <p>root is fine</p>
      </PrivateRoute>,
    )

    expect(screen.getByText('root is fine'))
  })

  it('admin ユーザーは requiredRole="root" ではリダイレクトされる', () => {
    setAuthState({ isLoggedIn: true, userRole: 'admin' })

    renderWithRouter(
      <PrivateRoute requiredRole="root">
        <p>root only</p>
      </PrivateRoute>,
    )

    expect(screen.getByTestId('path').textContent).toBe('/')
  })
})
