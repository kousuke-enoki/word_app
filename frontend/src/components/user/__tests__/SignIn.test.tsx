/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import {
  afterAll,
  afterEach,
  beforeAll,
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest'

import { server } from '@/__tests__/mswServer'
import { RuntimeConfigProvider } from '@/contexts/runtimeConfig/Provider'

/** ▼ useNavigate を先にモック（他の import より前） */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

import SignIn from '../SignIn'

/** ▼ ヘルパ：ENV と location を整える */
const ORIG_LOCATION = window.location
beforeAll(() => {
  // user-event と衝突しないように location を再定義可能に
  delete (window as any).location
  ;(window as any).location = { ...ORIG_LOCATION, href: 'about:blank' }
})
afterAll(() => {
  // 元に戻す
  delete (window as any).location
  ;(window as any).location = ORIG_LOCATION
})

beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
  // Vite 互換: Vitest の環境変数スタブ（無ければ下の fallback でもOK）
  if ((vi as any).stubEnv) {
    ;(vi as any).stubEnv('VITE_API_URL', 'https://api.example.com')
  } else {
    // fallback（古い Vitest 用）
    ;(import.meta as any).env = {
      ...(import.meta as any).env,
      VITE_API_URL: 'https://api.example.com',
    }
  }
  // beforeEachではデフォルトハンドラーを設定しない（各テストで明示的に設定）
})

afterEach(() => {
  if ((vi as any).unstubAllEnvs) (vi as any).unstubAllEnvs()
})

/** 入力ヘルパ */
const typeCredentials = async (email: string, password: string) => {
  const user = userEvent.setup()
  await user.type(screen.getByLabelText('Email'), email)
  await user.type(screen.getByLabelText('Password'), password)
  return user
}

/** コンポーネントをRuntimeConfigProviderでラップしてレンダリング */
const renderWithProvider = (ui: React.ReactElement) => {
  return render(
    <MemoryRouter>
      <RuntimeConfigProvider>{ui}</RuntimeConfigProvider>
    </MemoryRouter>,
  )
}

describe('SignIn Component', () => {
  it('初期表示：見出し/入力/ボタンがある（設定ロード中は LINE ボタン非表示）', async () => {
    // ローディング中の状態をシミュレート（遅延レスポンス）
    server.use(
      rest.get(
        'http://localhost:8080/public/runtime-config',
        async (_, res, ctx) => {
          await new Promise((resolve) => setTimeout(resolve, 100))
          return res(
            ctx.status(200),
            ctx.json({
              is_test_user_mode: false,
              is_line_authentication: false,
              version: '1.0.0',
            }),
          )
        },
      ),
    )

    renderWithProvider(<SignIn />)

    // 見出しなど
    expect(
      screen.getByRole('heading', { name: 'サインイン' }),
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
    expect(screen.getByLabelText('Password')).toBeInTheDocument()

    // ロード中は LINE ボタンなし（読み込み中表示）
    expect(screen.queryByRole('button', { name: 'LINEでログイン' })).toBeNull()
    expect(screen.getByText('読み込み中...')).toBeInTheDocument()

    // ローディング完了後もLINE認証が無効の場合は表示されない
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )
    expect(screen.queryByRole('button', { name: 'LINEでログイン' })).toBeNull()
  })

  it('設定API成功：LINE 有効なら LINE ボタンが表示される', async () => {
    // LINE認証が有効な状態
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: true,
            version: '1.0.0',
          }),
        ),
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    // 表示されるまで待つ
    const lineBtn = await screen.findByRole(
      'button',
      {
        name: 'LINEでログイン',
      },
      { timeout: 3000 },
    )
    expect(lineBtn).toBeInTheDocument()
  })

  it('設定API失敗：LINE ボタンは表示されない（サイレント）', async () => {
    // LINE認証が無効な状態（デフォルト）
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: false,
            version: '1.0.0',
          }),
        ),
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    // LINE認証が無効なのでボタンは表示されない
    expect(screen.queryByRole('button', { name: 'LINEでログイン' })).toBeNull()
  })

  it('メールサインイン成功：token を保存して /mypage へ navigate、ローディング表示が戻る', async () => {
    // 設定API
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: false,
            version: '1.0.0',
          }),
        ),
      ),
    )

    // サインイン成功ハンドラー
    let requestBody: any
    server.use(
      rest.post(
        'http://localhost:8080/users/sign_in',
        async (req, res, ctx) => {
          requestBody = await req.json()
          // リクエストインターセプターの検証：Authorizationヘッダーが付与されていないことを確認（サインイン時はトークンなし）
          const authHeader = req.headers.get('Authorization')
          expect(authHeader).toBeNull()
          await new Promise((resolve) => setTimeout(resolve, 50))
          return res(ctx.status(200), ctx.json({ token: 't-123' }))
        },
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    const user = userEvent.setup()
    await user.type(screen.getByLabelText('Email'), 'me@example.com')
    await user.type(screen.getByLabelText('Password'), 'secret123')

    const submitButton = screen.getByRole('button', {
      name: 'メールでサインイン',
    })
    await user.click(submitButton)

    // ← ここで loading が描画される猶予ができる
    expect(
      await screen.findByRole(
        'button',
        { name: 'サインイン中…' },
        { timeout: 3000 },
      ),
    ).toBeInTheDocument()

    // リクエストボディの検証
    await waitFor(
      () => {
        expect(requestBody).toEqual({
          email: 'me@example.com',
          password: 'secret123',
        })
      },
      { timeout: 3000 },
    )

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/mypage'), {
      timeout: 3000,
    })
    expect(localStorage.getItem('token')).toBe('t-123')

    // 最終的に文言が戻る
    await screen.findByRole(
      'button',
      { name: 'メールでサインイン' },
      { timeout: 3000 },
    )
  })

  it('メールサインイン失敗（サーバーから message あり）：その文言を表示し、ローディング解除', async () => {
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: false,
            version: '1.0.0',
          }),
        ),
      ),
    )

    // サインイン失敗ハンドラー（メッセージあり）
    server.use(
      rest.post(
        'http://localhost:8080/users/sign_in',
        async (req, res, ctx) => {
          // axiosのエラーハンドリングに合わせて、エラーレスポンスの構造を正しく設定
          return res(
            ctx.status(401),
            ctx.json({ message: 'ユーザーが存在しません' }),
          )
        },
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    const user = await typeCredentials('ng@example.com', 'wrong')
    await user.click(screen.getByRole('button', { name: 'メールでサインイン' }))

    // エラーメッセージ（タイムアウトを長めに設定）
    expect(
      await screen.findByText('ユーザーが存在しません', {}, { timeout: 3000 }),
    ).toBeInTheDocument()

    // ローディング解除
    await screen.findByRole(
      'button',
      { name: 'メールでサインイン' },
      { timeout: 3000 },
    )
    expect(navigateMock).not.toHaveBeenCalled()
    expect(localStorage.getItem('token')).toBeNull()
  })

  it('メールサインイン失敗（message なし）：デフォルト文言を表示', async () => {
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: true,
            version: '1.0.0',
          }),
        ),
      ),
    )

    // サインイン失敗ハンドラー（メッセージなし）
    server.use(
      rest.post(
        'http://localhost:8080/users/sign_in',
        async (req, res, ctx) => {
          return res(ctx.status(500), ctx.json({}))
        },
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    const user = await typeCredentials('x@example.com', 'xxx')
    await user.click(screen.getByRole('button', { name: 'メールでサインイン' }))

    expect(
      await screen.findByText(
        'サインインに失敗しました',
        {},
        { timeout: 3000 },
      ),
    ).toBeInTheDocument()
    await screen.findByRole(
      'button',
      { name: 'メールでサインイン' },
      { timeout: 3000 },
    )
    expect(navigateMock).not.toHaveBeenCalled()
  })

  it('LINEでログイン：クリックで window.location.href が LINE ログイン URL に変わる', async () => {
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: true,
            version: '1.0.0',
          }),
        ),
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    const btn = await screen.findByRole(
      'button',
      { name: 'LINEでログイン' },
      { timeout: 3000 },
    )
    const user = userEvent.setup()
    await user.click(btn)

    expect(window.location.href).toBe(
      'https://api.example.com/users/auth/line/login',
    )
  })

  it('サインアップ導線：/sign_up リンクがある', async () => {
    server.use(
      rest.get('http://localhost:8080/public/runtime-config', (_, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            is_test_user_mode: false,
            is_line_authentication: false,
            version: '1.0.0',
          }),
        ),
      ),
    )

    renderWithProvider(<SignIn />)

    // ローディング完了を待つ
    await waitFor(
      () => {
        expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    const link = screen.getByRole('link', { name: 'サインアップ' })
    expect(link).toHaveAttribute('href', '/sign_up')
  })
})
