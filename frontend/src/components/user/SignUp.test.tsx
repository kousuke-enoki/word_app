import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import SignUp from './SignUp'

// MSWのセットアップ
const server = setupServer()

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
beforeEach(() => {
  localStorage.clear()
})
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('SignUp Component', () => {
  test('フォーム要素が正しく表示される', () => {
    render(
      <MemoryRouter>
        <SignUp />
      </MemoryRouter>,
    )

    // 見出し
    expect(
      screen.getByRole('heading', { name: 'サインアップ' }),
    ).toBeInTheDocument()
    // 入力フォーム
    expect(screen.getByLabelText('Name:')).toBeInTheDocument()
    expect(screen.getByLabelText('Email:')).toBeInTheDocument()
    expect(screen.getByLabelText('Password:')).toBeInTheDocument()
    // ボタン
    expect(
      screen.getByRole('button', { name: 'サインアップ' }),
    ).toBeInTheDocument()
  })

  test('サインアップ成功時、トークンをlocalStorageに保存し、成功メッセージが表示される', async () => {
    // 200レスポンスで成功をモック
    server.use(
      rest.post('http://localhost:8080/users/sign_up', (req, res, ctx) => {
        return res(
          ctx.status(200),
          ctx.json({
            token: 'mocked-signup-token',
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <SignUp />
      </MemoryRouter>,
    )

    // フォームに入力
    await userEvent.type(screen.getByLabelText('Name:'), 'Test User')
    await userEvent.type(screen.getByLabelText('Email:'), 'test@example.com')
    await userEvent.type(screen.getByLabelText('Password:'), 'password123')

    // サインアップボタンをクリック
    await userEvent.click(screen.getByRole('button', { name: 'サインアップ' }))

    // 成功メッセージが表示されるまで待機
    await waitFor(() => {
      expect(screen.getByText('Sign up successful!')).toBeInTheDocument()
    })

    // localStorageにトークンが保存されている
    expect(localStorage.getItem('token')).toBe('mocked-signup-token')
  })

  test('サインアップ失敗時、エラーメッセージ一覧が表示され、成功メッセージは表示されない', async () => {
    // 400等で失敗レスポンスをモックし、FieldErrorを返す
    server.use(
      rest.post('http://localhost:8080/users/sign_up', (req, res, ctx) => {
        return res(
          ctx.status(400),
          ctx.json({
            errors: [
              { field: 'email', message: 'Email is already taken' },
              {
                field: 'password',
                message: 'Password must be at least 8 characters',
              },
            ],
          }),
        )
      }),
    )

    render(
      <MemoryRouter>
        <SignUp />
      </MemoryRouter>,
    )

    // フォームに誤った入力(重複メール、短いパスワードなど)
    await userEvent.type(screen.getByLabelText('Name:'), 'User')
    await userEvent.type(
      screen.getByLabelText('Email:'),
      'existing@example.com',
    )
    await userEvent.type(screen.getByLabelText('Password:'), 'short')

    // 送信
    await userEvent.click(screen.getByRole('button', { name: 'サインアップ' }))

    // エラーメッセージが表示されること
    await waitFor(() => {
      expect(screen.getByText('Email is already taken')).toBeInTheDocument()
      // eslint-disable-next-line testing-library/no-wait-for-multiple-assertions
      expect(
        screen.getByText('Password must be at least 8 characters'),
      ).toBeInTheDocument()
    })

    // 成功メッセージは表示されない
    expect(screen.queryByText('Sign up successful!')).not.toBeInTheDocument()

    // localStorageにはトークンが保存されていない
    expect(localStorage.getItem('token')).toBeNull()
  })
})
