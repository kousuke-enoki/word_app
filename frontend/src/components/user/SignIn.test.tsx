import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import SignIn from './SignIn'

/**
 * MSW を使わないシンプルな axiosモックでもOKですが、
 * 他のコンポーネントと同様に MSW を使用する前提として例示します。
 *
 * 以下の handlers を独自に設定している想定です。
 * テスト内で server.use(...) で切り替え可能です。
 */

// テスト用サーバーを準備
const server = setupServer()

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => {
  server.resetHandlers()
  localStorage.clear() // 毎テスト後に localStorage をクリア
})
afterAll(() => server.close())

describe('SignIn Component', () => {
  test('フォーム要素が正しく表示される', () => {
    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    // 見出し
    expect(
      screen.getByRole('heading', { name: 'サインイン' }),
    ).toBeInTheDocument()
    // Email, Passwordのラベル
    expect(screen.getByLabelText('Email:')).toBeInTheDocument()
    expect(screen.getByLabelText('Password:')).toBeInTheDocument()
    // ボタン
    expect(
      screen.getByRole('button', { name: 'サインイン' }),
    ).toBeInTheDocument()
  })

  test('サインイン成功時、localStorage にトークンが保存されメッセージが表示される', async () => {
    // サーバーが成功レスポンス(200)を返すようにモック
    server.use(
      rest.post('http://localhost:8080/users/sign_in', (req, res, ctx) => {
        return res(ctx.status(200), ctx.json({ token: 'mocked-jwt-token' }))
      }),
    )

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    // ユーザー入力
    await userEvent.type(screen.getByLabelText('Email:'), 'test@example.com')
    await userEvent.type(screen.getByLabelText('Password:'), 'Password123!')

    // フォーム送信
    await userEvent.click(screen.getByRole('button', { name: 'サインイン' }))

    // 成功メッセージが表示されるまで待機
    await waitFor(() => {
      expect(screen.getByText('Sign in successful!')).toBeInTheDocument()
    })

    // localStorage にトークンが保存されている
    expect(localStorage.getItem('token')).toBe('mocked-jwt-token')
  })

  test('サインイン失敗時、エラーメッセージが表示される', async () => {
    // サーバーが失敗レスポンス(401など)を返すようにモック
    server.use(
      rest.post('http://localhost:8080/users/sign_in', (req, res, ctx) => {
        return res(ctx.status(401), ctx.json({ error: 'Invalid credentials' }))
      }),
    )

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    await userEvent.type(screen.getByLabelText('Email:'), 'bad@example.com')
    await userEvent.type(screen.getByLabelText('Password:'), 'wrongpassword')
    await userEvent.click(screen.getByRole('button', { name: 'サインイン' }))

    // 失敗メッセージが表示される
    await waitFor(() => {
      expect(
        screen.getByText('Sign in failed. Please try again.'),
      ).toBeInTheDocument()
    })

    // トークンは保存されていない
    expect(localStorage.getItem('token')).toBeNull()
  })
})
