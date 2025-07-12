/* eslint-disable @typescript-eslint/no-explicit-any */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'

const navigateMock = vi.fn()

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>(
    'react-router-dom',
  )
  return { ...actual, useNavigate: () => navigateMock }
})

vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

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

    render(
      <MemoryRouter>
        <LineCallback />
      </MemoryRouter>,
    )

    expect(
      screen.getByText('LINE 認証処理中...'),
    )
  })

  /** 2. コールバック成功で token 保存＆ /mypage へ */
  it('LINE 認証成功時、token 保存して /mypage へ遷移', async () => {
    setSearch('?code=sucCode&state=sucState')
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { token: 'jwt-token' },
    })

    render(
      <MemoryRouter>
        <LineCallback />
      </MemoryRouter>,
    )

    await waitFor(() => {
      expect(axiosInstance.get).toHaveBeenCalledWith(
        '/users/auth/line/callback',
        { params: { code: 'sucCode', state: 'sucState' } },
      )
      expect(localStorage.getItem('token')).toBe('jwt-token')
      expect(navigateMock).toHaveBeenCalledWith('/mypage')
    })
  })

  /** 3. 失敗したら /signin?err=line へ */
  it('認証エラー時は /signin?err=line へリダイレクト', async () => {
    setSearch('?code=badCode&state=badState')
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('401'))

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
