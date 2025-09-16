// src/components/setting/__tests__/UserSetting.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import UserSetting from '../UserSetting'

/* -------- 依存を薄くモック -------- */
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn(), post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

const setThemeMock = vi.fn()
vi.mock('@/contexts/themeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

// 画面下部ナビ（useNavigate 依存のため見た目だけ）
vi.mock('../../common/PageBottomNav', () => ({
  default: (props: any) => <div data-testid="PageBottomNav" {...props} />,
}))

// Card はそのままでもOK。安定化したければ薄くモックしても良い。
// vi.mock('../../ui/card', () => ({ Card: ({ children, ...rest }: any) => <div {...rest}>{children}</div> }))

/* -------- ヘルパ -------- */
const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>)

const getControls = () => {
  const checkbox = screen.getByRole('checkbox') as HTMLInputElement
  const saveBtn = screen.getByRole('button', { name: /保存|保存中/ })
  const onOffText = () => screen.getByText(/ON|OFF/) // 表示テキスト
  return { checkbox, saveBtn, onOffText }
}

beforeEach(() => {
  vi.resetAllMocks()
  vi.spyOn(window, 'alert').mockImplementation(() => {})
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

describe('UserSetting', () => {
  it('初期 GET 成功：is_dark_mode=true → チェックON/ON表示、保存はdisabled', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: true } },
    })

    renderWithRouter(<UserSetting />)

    expect(
      screen.getByRole('heading', { name: 'ユーザー設定画面' }),
    ).toBeInTheDocument()

    const { checkbox, saveBtn, onOffText } = getControls()

    await waitFor(() => {
      expect(checkbox).toBeChecked()
      expect(onOffText()).toHaveTextContent('ON')
    })

    expect(saveBtn).toBeDisabled()
    expect(
      screen.queryByText(/設定を保存しました。|保存に失敗|保存中にエラー/),
    ).toBeNull()
  })

  it('初期 GET 失敗：alert が呼ばれ、デフォルト false/OFF、保存disabled', async () => {
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('get fail'))

    renderWithRouter(<UserSetting />)

    await waitFor(() =>
      expect(window.alert).toHaveBeenCalledWith(
        'ユーザー設定の取得中にエラーが発生しました。',
      ),
    )

    const { checkbox, onOffText, saveBtn } = getControls()
    expect(checkbox).not.toBeChecked()
    expect(onOffText()).toHaveTextContent('OFF')
    expect(saveBtn).toBeDisabled()
  })

  it('dirty 判定：トグルで enabled、元に戻すと disabled', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: false } },
    })

    renderWithRouter(<UserSetting />)
    const { checkbox, saveBtn } = getControls()

    await screen.findByRole('checkbox')

    // 変更 → enabled
    await userEvent.click(checkbox)
    expect(saveBtn).toBeEnabled()

    // 元に戻す → disabled
    await userEvent.click(checkbox)
    expect(saveBtn).toBeDisabled()
  })

  it('保存成功(200)：dark に保存 → setTheme("dark")、メッセージ、dirty解除', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: false } },
    })
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      status: 200,
      data: {},
    })

    renderWithRouter(<UserSetting />)
    const { checkbox } = getControls()
    await screen.findByRole('checkbox')

    // false -> true にして保存
    await userEvent.click(checkbox)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    // POST と setTheme
    await waitFor(() => {
      expect(axiosInstance.post).toHaveBeenCalledWith('/setting/user_config', {
        is_dark_mode: true,
      })
    })
    expect(setThemeMock).toHaveBeenCalledWith('dark')

    // メッセージ & dirty解除
    expect(await screen.findByText('設定を保存しました。')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })

  it('保存成功(200)：light に保存 → setTheme("light")', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: true } },
    })
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      status: 200,
      data: {},
    })

    renderWithRouter(<UserSetting />)
    const { checkbox } = getControls()
    await screen.findByRole('checkbox')

    // true -> false にして保存
    await userEvent.click(checkbox)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    await waitFor(() => {
      expect(axiosInstance.post).toHaveBeenCalledWith('/setting/user_config', {
        is_dark_mode: false,
      })
    })
    expect(setThemeMock).toHaveBeenCalledWith('light')
    expect(await screen.findByText('設定を保存しました。')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })

  it('保存 非200：メッセージ「保存に失敗しました。」、dirty 継続', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: false } },
    })
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      status: 500,
      data: {},
    })

    renderWithRouter(<UserSetting />)
    const { checkbox } = getControls()
    await screen.findByRole('checkbox')

    await userEvent.click(checkbox)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    expect(await screen.findByText('保存に失敗しました。')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled() // まだdirty
    expect(setThemeMock).not.toHaveBeenCalled()
  })

  it('保存 例外：メッセージ「保存中にエラーが発生しました。」、dirty 継続', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: false } },
    })
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('post fail'))

    renderWithRouter(<UserSetting />)
    const { checkbox } = getControls()
    await screen.findByRole('checkbox')

    await userEvent.click(checkbox)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    expect(
      await screen.findByText('保存中にエラーが発生しました。'),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled()
    expect(setThemeMock).not.toHaveBeenCalled()
  })

  it('ローディング：保存中… 表示中は disabled → 完了後にメッセージ＆disabled', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { Config: { is_dark_mode: false } },
    })

    let resolvePost: (v: any) => void
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolvePost = res
        }),
    )

    renderWithRouter(<UserSetting />)
    const { checkbox } = getControls()
    await screen.findByRole('checkbox')

    await userEvent.click(checkbox)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    const saving = screen.getByRole('button', { name: '保存中...' })
    expect(saving).toBeDisabled()

    resolvePost!({ status: 200, data: {} })
    expect(await screen.findByText('設定を保存しました。')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })
})
