import { render, screen, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'

import { server } from '@/__tests__/mswServer'
import MyPage from '@/components/user/MyPage'
import PrivateRoute from '@/routes/PrivateRoute'

const LocationDisplay = () => {
  const location = useLocation()
  return <div data-testid="location">{location.pathname}</div>
}

const renderWithAuthRoutes = (initialPath = '/mypage') =>
  render(
    <MemoryRouter initialEntries={[initialPath]}>
      <LocationDisplay />
      <Routes>
        <Route path="/" element={<div>Login Screen</div>} />
        <Route
          path="/mypage"
          element={
            <PrivateRoute>
              <MyPage />
            </PrivateRoute>
          }
        />
      </Routes>
    </MemoryRouter>,
  )

describe('認証・認可統合テスト', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('未ログインで保護ルートにアクセスするとログイン画面へリダイレクトされる', async () => {
    renderWithAuthRoutes('/mypage')

    await waitFor(() => {
      expect(screen.getByTestId('location').textContent).toBe('/')
    })
    expect(screen.getByText('Login Screen')).toBeInTheDocument()
  })

  it('ロール不足のユーザーは管理系リンクが表示されない', async () => {
    localStorage.setItem('token', 'user-token')
    server.use(
      rest.get('http://localhost:8080/auth/check', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            user: { id: 10, isAdmin: false, isRoot: false, isTest: false },
          }),
        ),
      ),
      rest.get('http://localhost:8080/users/my_page', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            user: {
              id: 10,
              name: '一般ユーザー',
              isAdmin: false,
              isRoot: false,
            },
          }),
        ),
      ),
    )

    renderWithAuthRoutes('/mypage')

    await screen.findByText(/一般ユーザー\s*さん/)
    expect(screen.getByTestId('location').textContent).toBe('/mypage')
    expect(screen.queryByRole('link', { name: /単語登録/ })).toBeNull()
    expect(screen.queryByRole('link', { name: /管理設定/ })).toBeNull()
    expect(screen.getByText(/👤\s*User/)).toBeInTheDocument()
  })

  it('十分なロールのユーザーには管理系リンクが表示される', async () => {
    localStorage.setItem('token', 'admin-token')
    server.use(
      rest.get('http://localhost:8080/auth/check', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            user: { id: 20, isAdmin: true, isRoot: false, isTest: false },
          }),
        ),
      ),
      rest.get('http://localhost:8080/users/my_page', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            user: { id: 20, name: '管理者', isAdmin: true, isRoot: false },
          }),
        ),
      ),
    )

    renderWithAuthRoutes('/mypage')

    await screen.findByText(/管理者\s*さん/)
    expect(screen.getByTestId('location').textContent).toBe('/mypage')
    const adminLink = await screen.findByRole('link', { name: /単語登録/ })
    expect(adminLink).toHaveAttribute('href', '/words/new')
    expect(screen.getByText(/🔧\s*Admin/)).toBeInTheDocument()
  })
})
