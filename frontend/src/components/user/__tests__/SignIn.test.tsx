/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
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

/** ▼ useNavigate を先にモック（他の import より前） */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

/** ▼ axios モック */
vi.mock('@/axiosConfig', () => ({
  default: {
    get: vi.fn(), // /public/runtime-config
    post: vi.fn(), // /users/sign_in
  },
}))
import axiosInstance from '@/axiosConfig'

/** ▼ useRuntimeConfig をモック */
const useRuntimeConfigMock = vi.fn()
vi.mock('@/contexts/runtimeConfig/useRuntimeConfig', () => ({
  useRuntimeConfig: () => useRuntimeConfigMock(),
}))

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
  // useRuntimeConfig のデフォルト値（ローディング中）
  useRuntimeConfigMock.mockReturnValue({
    config: { is_test_user_mode: false, is_line_authentication: false },
    isLoading: true,
  })
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

describe('SignIn Component', () => {
  it('初期表示：見出し/入力/ボタンがある（設定ロード中は LINE ボタン非表示）', async () => {
    // ローディング中の状態
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: true,
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

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
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: false,
    })
    // 再レンダリングをシミュレートするために、別のテストで確認
  })

  it('設定API成功：LINE 有効なら LINE ボタンが表示される', async () => {
    // LINE認証が有効な状態
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: true },
      isLoading: false,
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    // 表示されるまで待つ
    const lineBtn = await screen.findByRole('button', {
      name: 'LINEでログイン',
    })
    expect(lineBtn).toBeInTheDocument()
  })

  it('設定API失敗：LINE ボタンは表示されない（サイレント）', async () => {
    // LINE認証が無効な状態（デフォルト）
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: false,
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    // LINE認証が無効なのでボタンは表示されない
    expect(screen.queryByRole('button', { name: 'LINEでログイン' })).toBeNull()
  })

  it('メールサインイン成功：token を保存して /mypage へ navigate、ローディング表示が戻る', async () => {
    // 設定API
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: false,
    })
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((resolve) =>
          setTimeout(() => resolve({ data: { token: 't-123' } }), 50),
        ),
    )
    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const user = userEvent.setup()
    await user.type(screen.getByLabelText('Email'), 'me@example.com')
    await user.type(screen.getByLabelText('Password'), 'secret123')

    await user.click(screen.getByRole('button', { name: 'メールでサインイン' }))

    // ← ここで loading が描画される猶予ができる
    expect(
      await screen.findByRole('button', { name: 'サインイン中…' }),
    ).toBeInTheDocument()

    await waitFor(() =>
      expect(axiosInstance.post).toHaveBeenCalledWith('/users/sign_in', {
        email: 'me@example.com',
        password: 'secret123',
      }),
    )

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/mypage'))
    expect(localStorage.getItem('token')).toBe('t-123')

    // 最終的に文言が戻る
    await screen.findByRole('button', { name: 'メールでサインイン' })
  })

  it('メールサインイン失敗（サーバーから message あり）：その文言を表示し、ローディング解除', async () => {
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: false,
    })
    ;(axiosInstance.post as any).mockRejectedValueOnce({
      response: { data: { message: 'ユーザーが存在しません' } },
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const user = await typeCredentials('ng@example.com', 'wrong')
    await user.click(screen.getByRole('button', { name: 'メールでサインイン' }))

    // エラーメッセージ
    expect(
      await screen.findByText('ユーザーが存在しません'),
    ).toBeInTheDocument()

    // ローディング解除
    await screen.findByRole('button', { name: 'メールでサインイン' })
    expect(navigateMock).not.toHaveBeenCalled()
    expect(localStorage.getItem('token')).toBeNull()
  })

  it('メールサインイン失敗（message なし）：デフォルト文言を表示', async () => {
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: true },
      isLoading: false,
    })
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('boom'))

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const user = await typeCredentials('x@example.com', 'xxx')
    await user.click(screen.getByRole('button', { name: 'メールでサインイン' }))

    expect(
      await screen.findByText('サインインに失敗しました'),
    ).toBeInTheDocument()
    await screen.findByRole('button', { name: 'メールでサインイン' })
    expect(navigateMock).not.toHaveBeenCalled()
  })

  it('LINEでログイン：クリックで window.location.href が LINE ログイン URL に変わる', async () => {
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: true },
      isLoading: false,
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const btn = await screen.findByRole('button', { name: 'LINEでログイン' })
    const user = userEvent.setup()
    await user.click(btn)

    expect(window.location.href).toBe(
      'https://api.example.com/users/auth/line/login',
    )
  })

  it('サインアップ導線：/sign_up リンクがある', () => {
    useRuntimeConfigMock.mockReturnValue({
      config: { is_test_user_mode: false, is_line_authentication: false },
      isLoading: false,
    })

    render(
      <MemoryRouter>
        <SignIn />
      </MemoryRouter>,
    )

    const link = screen.getByRole('link', { name: 'サインアップ' })
    expect(link).toHaveAttribute('href', '/sign_up')
  })
})
