/* eslint-disable @typescript-eslint/no-explicit-any */
/* ---------------- はじめに：useNavigate を先モック ---------------- */
import { beforeEach,describe, expect, it, vi } from 'vitest'

const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

/* ---------------- axios を丸ごとモック ---------------- */
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn(), get: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

/* ---------------- ThemeContext も最低限だけ ---------------- */
export const setThemeMock = vi.fn()
vi.mock('@/contexts/ThemeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

/* ---------------- テストに必要な依存 ---------------- */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'

import SignUp from '../SignUp'

/* ---------------- 毎回きれいな状態で開始 ---------------- */
beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})

/* =====================================================
 *                       TESTS
 * ===================================================== */
describe('SignUp Component', () => {
  /* ---------- 1. サインアップ成功フロー ---------- */
  it('正常にサインアップすると token 保存 → /mypage へ遷移', async () => {
    /* 1) /users/sign_up → JWT を返す */
    ;(axiosInstance.post as any).mockResolvedValueOnce({ data: { token: 'jwt-signup' } })
    /* 2) /setting/user_config → ダークモード無効 (今回は呼ばれないが用意だけ) */
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: { is_dark_mode: false } })

    render(
      <MemoryRouter>
        <SignUp />
      </MemoryRouter>,
    )

    // -------- 入力はリアルタイマーで行う --------
    await userEvent.type(screen.getByLabelText('Name:'),     'Taro')
    await userEvent.type(screen.getByLabelText('Email:'),    'taro@example.com')
    await userEvent.type(screen.getByLabelText('Password:'), 'secret123')
    await userEvent.click(screen.getByRole('button', { name: 'サインアップ' }))

    // -------- submit 後だけフェイクタイマー --------
    vi.useFakeTimers()
    await vi.runAllTicks()        // axios → state 更新
    vi.runOnlyPendingTimers()     // setTimeout(0) で navigate('/')
    await vi.runAllTicks()
    vi.useRealTimers()

    // -------- 検証 --------
    expect(localStorage.getItem('token')).toBe('jwt-signup')
    expect(localStorage.getItem('logoutMessage')).toBe('サインアップしました。')
    expect(screen.getByText('Sign up successful!'))
    expect(screen.getByText('Sign up successful!'))           // component の文言
    expect(navigateMock).toHaveBeenCalledWith('/mypage')
  })

  /* ---------- 2. バリデーションエラー ---------- */
  it('サインアップ失敗時はフィールドエラーを表示 & token は保存しない', async () => {
    const apiError = {
      response: {
        data: {
          errors: [
            { field: 'email',    message: 'メールが既に使用されています' },
            { field: 'password', message: '8文字以上で入力してください' },
          ],
        },
      },
    }
    ;(axiosInstance.post as any).mockRejectedValueOnce(apiError)

    render(
      <MemoryRouter>
        <SignUp />
      </MemoryRouter>,
    )

    await userEvent.type(screen.getByLabelText('Name:'),     'Hanako')
    await userEvent.type(screen.getByLabelText('Email:'),    'hanako@example.com')
    await userEvent.type(screen.getByLabelText('Password:'), 'short')
    await userEvent.click(screen.getByRole('button', { name: 'サインアップ' }))

    /* waitFor で「エラーメッセージ2つ出た」ことを確認 */
    await waitFor(() => {
      expect(screen.getByText('メールが既に使用されています'))
      expect(screen.getByText('8文字以上で入力してください'))
    })

    expect(localStorage.getItem('token')).toBeNull()
    expect(navigateMock).not.toHaveBeenCalled()
  })
})
