// src/components/user/__tests__/SignUp.test.tsx

import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { server } from '@/__tests__/mswServer'

import SignUp from '../SignUp'

/** ------------ ルーターの navigate をモック（必ず他の import より前） ------------ */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

/** ------------ ThemeContext をモック ------------ */
const setThemeMock = vi.fn()
vi.mock('@/contexts/themeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

/** ------------ 共通ヘルパ ------------ */
const renderSignUp = () =>
  render(
    <MemoryRouter>
      <SignUp />
    </MemoryRouter>,
  )

const typeAll = async (name: string, email: string, password: string) => {
  const user = userEvent.setup()
  const by = (label: string) => screen.getByLabelText(label)

  if (name) await user.type(by('Name'), name)
  if (email) await user.type(by('Email'), email)
  if (password) await user.type(by('Password'), password)
  return user
}

beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})

afterEach(() => {
  // 念のためタイマーを実環境に戻す
  try {
    vi.useRealTimers()
  } catch {
    /* empty */
  }
})

describe('SignUp Component', () => {
  it('初期表示：見出し/入力/ボタン/リンクがある', () => {
    renderSignUp()

    // 見出し
    expect(
      screen.getByRole('heading', { name: 'サインアップ' }),
    ).toBeInTheDocument()

    // 入力（必須が付いているかも見る）
    const name = screen.getByLabelText('Name')
    const email = screen.getByLabelText('Email')
    const password = screen.getByLabelText('Password')
    expect(name).toBeRequired()
    expect(email).toBeRequired()
    expect(password).toBeRequired()

    // 送信ボタン（初期は有効で文言は「サインアップ」）
    const submit = screen.getByRole('button', { name: 'サインアップ' })
    expect(submit).toBeEnabled()

    // サインインへのリンク
    expect(screen.getByRole('link', { name: 'サインイン' })).toHaveAttribute(
      'href',
      '/sign_in',
    )
  })

  it('成功：token保存・/my_page遷移・テーマ適用・ボタン戻る（中間表示も）', async () => {
    // 設定APIは即成功
    server.use(
      rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
        res(ctx.status(200), ctx.json({ is_dark_mode: false })),
      ),
      // サインアップAPIは遅延を入れてpending状態を確認できるようにする
      rest.post(
        'http://localhost:8080/users/sign_up',
        async (req, res, ctx) => {
          // サインアップ時はトークンなしなので、Authorizationヘッダーがnullであることを確認
          const authHeader = req.headers.get('Authorization')
          expect(authHeader).toBeNull()
          await new Promise((resolve) => setTimeout(resolve, 50))
          return res(ctx.status(200), ctx.json({ token: 'tok-1' }))
        },
      ),
    )

    renderSignUp()
    const user = await typeAll('Alice', 'alice@example.com', 'passw0rd')

    const submit = screen.getByRole('button', { name: 'サインアップ' })
    await user.click(submit)

    // 送信中の確認
    const loadingButton = await screen.findByRole('button', {
      name: /サインアップ中/,
    })
    await waitFor(
      () => {
        expect(loadingButton).toBeDisabled()
      },
      { timeout: 1000 },
    )

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/my_page'))
    expect(localStorage.getItem('token')).toBe('tok-1')
    expect(localStorage.getItem('logoutMessage')).toBe('サインアップしました。')
    expect(setThemeMock).toHaveBeenCalledWith('light')

    // 最終的にボタンは戻る
    await waitFor(() =>
      expect(
        screen.getByRole('button', { name: 'サインアップ' }),
      ).toBeEnabled(),
    )
  })

  it('成功：dark テーマが適用される（is_dark_mode:true）', async () => {
    // 設定APIは成功（ダーク）
    server.use(
      rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
        res(ctx.status(200), ctx.json({ is_dark_mode: true })),
      ),
      rest.post('http://localhost:8080/users/sign_up', (_, res, ctx) =>
        res(ctx.status(200), ctx.json({ token: 'tok-dark' })),
      ),
    )

    renderSignUp()
    const user = await typeAll('Bob', 'bob@example.com', 'secretxxx')

    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/my_page'))
    expect(localStorage.getItem('token')).toBe('tok-dark')
    expect(localStorage.getItem('logoutMessage')).toBe('サインアップしました。')
    expect(setThemeMock).toHaveBeenCalledWith('dark')
  })

  it('成功：ユーザー設定取得が失敗しても遷移は行われ、setTheme は呼ばれない', async () => {
    server.use(
      rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
        res(ctx.status(500)),
      ),
      rest.post('http://localhost:8080/users/sign_up', (_, res, ctx) =>
        res(ctx.status(200), ctx.json({ token: 'tok-ok' })),
      ),
    )

    renderSignUp()
    const user = await typeAll('Carol', 'carol@example.com', 'xxxyyyzzz')
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/my_page'))
    expect(localStorage.getItem('token')).toBe('tok-ok')
    expect(localStorage.getItem('logoutMessage')).toBe('サインアップしました。')
    expect(setThemeMock).not.toHaveBeenCalled()
  })

  it('失敗：field errors + message を表示し、ボタン状態も戻る（navigate なし / token 未保存）', async () => {
    server.use(
      rest.post('http://localhost:8080/users/sign_up', (_, res, ctx) =>
        res(
          ctx.status(400),
          ctx.json({
            message: '入力に不備があります',
            errors: [
              { field: 'name', message: '名前を入力してください' },
              { field: 'email', message: 'メール形式が不正です' },
              { field: 'password', message: '8文字以上で入力してください' },
            ],
          }),
        ),
      ),
    )

    renderSignUp()
    const user = await typeAll('Taro', 'taro@example.com', 'password123') // ← 妥当な値
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    expect(await screen.findByText('入力に不備があります')).toBeInTheDocument()
    expect(screen.getByText('名前を入力してください')).toBeInTheDocument()
    expect(screen.getByText('メール形式が不正です')).toBeInTheDocument()
    expect(screen.getByText('8文字以上で入力してください')).toBeInTheDocument()

    // ボタンは戻る / navigate されない / token なし
    await waitFor(() =>
      expect(
        screen.getByRole('button', { name: 'サインアップ' }),
      ).toBeEnabled(),
    )
    expect(navigateMock).not.toHaveBeenCalled()
    expect(localStorage.getItem('token')).toBeNull()
  })

  it('失敗：message無し→デフォ文言', async () => {
    server.use(
      rest.post('http://localhost:8080/users/sign_up', (_, res, ctx) =>
        res(ctx.status(400), ctx.json({ errors: [] })),
      ),
    )

    renderSignUp()
    const user = await typeAll('Dave', 'dave@example.com', 'passpass')
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    expect(
      await screen.findByText('サインアップに失敗しました'),
    ).toBeInTheDocument()
    expect(navigateMock).not.toHaveBeenCalled()
    expect(localStorage.getItem('token')).toBeNull()
  })

  it('二重送信防止：送信中はdisabled・postは1回', async () => {
    let requestCount = 0
    server.use(
      rest.get('http://localhost:8080/setting/user_config', (_, res, ctx) =>
        res(ctx.status(200), ctx.json({ is_dark_mode: false })),
      ),
      rest.post(
        'http://localhost:8080/users/sign_up',
        async (req, res, ctx) => {
          requestCount++
          // 遅延を入れて、2回目のクリックが無視されることを確認できるようにする
          await new Promise((resolve) => setTimeout(resolve, 100))
          return res(ctx.status(200), ctx.json({ token: 'tok-dup' }))
        },
      ),
    )

    renderSignUp()
    const user = await typeAll('Eve', 'eve@example.com', 'abcd1234')
    const submit = screen.getByRole('button', { name: 'サインアップ' })

    await user.click(submit)
    await user.click(submit) // 2回目は無視されるべき

    const loadingButton = await screen.findByRole('button', {
      name: /サインアップ中/,
    })
    await waitFor(
      () => {
        expect(loadingButton).toBeDisabled()
      },
      { timeout: 1000 },
    )
    expect(requestCount).toBe(1)

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/my_page'))
  })
})
