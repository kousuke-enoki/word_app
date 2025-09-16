// src/components/setting/__tests__/RootSetting.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import RootSetting from '../RootSetting'

/* ---------- 依存の薄いモック ---------- */
// Router 依存を避けるなら PageBottomNav をダミー化（見た目だけ）
vi.mock('../../common/PageBottomNav', () => ({
  default: (props: any) => <div data-testid="PageBottomNav" {...props} />,
}))
// UIカードはそのままでもOK。安定のため薄くしても良い。
// vi.mock('../../ui/card', () => ({ Card: ({ children, ...rest }: any) => <div {...rest}>{children}</div> }))

/* ---------- axios モック ---------- */
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn(), post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

/* ---------- ヘルパ ---------- */
const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>)

const getControls = () => {
  const combobox = screen.getByRole('combobox') as HTMLSelectElement
  const cbs = screen.getAllByRole('checkbox') as HTMLInputElement[]
  const [cbTest, cbMail, cbLine] = cbs
  const saveBtn = screen.getByRole('button', { name: /保存|保存中/ })
  return { combobox, cbTest, cbMail, cbLine, saveBtn }
}

beforeEach(() => {
  vi.resetAllMocks()
  vi.spyOn(window, 'alert').mockImplementation(() => {})
  vi.spyOn(console, 'error').mockImplementation(() => {}) // 失敗系のログを黙らせる
})

describe('RootSetting', () => {
  it('初期 GET 成功：UI に反映され、未変更のため保存は disabled', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'root',
          is_test_user_mode: true,
          is_email_authentication_check: true,
          is_line_authentication: false,
        },
      },
    })

    renderWithRouter(<RootSetting />)

    // 見出し
    expect(
      screen.getByRole('heading', { name: '管理設定画面' }),
    ).toBeInTheDocument()

    const { combobox, cbTest, cbMail, cbLine, saveBtn } = getControls()

    // 反映（待つ）
    await waitFor(() => {
      expect(combobox).toHaveValue('root')
      expect(cbTest).toBeChecked()
      expect(cbMail).toBeChecked()
      expect(cbLine).not.toBeChecked()
    })

    // 変更なし → 保存 disabled、メッセージなし
    expect(saveBtn).toBeDisabled()
    expect(
      screen.queryByText(/設定が保存されました。|保存に失敗|エラーが発生/),
    ).toBeNull()
  })

  it('初期 GET 失敗：alert が表示され、デフォルト値のまま', async () => {
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('GET fail'))

    renderWithRouter(<RootSetting />)

    await waitFor(() =>
      expect(window.alert).toHaveBeenCalledWith(
        'ルート設定の取得中にエラーが発生しました。',
      ),
    )

    const { combobox, cbTest, cbMail, cbLine, saveBtn } = getControls()
    // デフォルト state（admin/false/false/false）
    expect(combobox).toHaveValue('admin')
    expect(cbTest).not.toBeChecked()
    expect(cbMail).not.toBeChecked()
    expect(cbLine).not.toBeChecked()
    expect(saveBtn).toBeDisabled()
  })

  it('dirty 判定：変更 → enabled、元に戻す → disabled', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })

    renderWithRouter(<RootSetting />)
    const { combobox, saveBtn } = getControls()

    await screen.findByRole('combobox') // 初期反映

    // 変更 → enabled
    await userEvent.selectOptions(combobox, 'user')
    expect(saveBtn).toBeEnabled()

    // 元に戻す → disabled
    await userEvent.selectOptions(combobox, 'admin')
    expect(saveBtn).toBeDisabled()
  })

  it('保存成功 (status 200)：POST 内容が正しく、メッセージ表示＆dirty 解除', async () => {
    // 初期値
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })
    // 保存成功
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      status: 200,
      data: {},
    })

    renderWithRouter(<RootSetting />)
    const { combobox, cbTest, cbMail } = getControls()
    await screen.findByRole('combobox')

    // いくつか変更
    await userEvent.selectOptions(combobox, 'user')
    await userEvent.click(cbTest)
    await userEvent.click(cbMail)
    // line はそのまま false

    // 保存
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    // POST された payload
    await waitFor(() => {
      expect(axiosInstance.post).toHaveBeenCalledWith('/setting/root_config', {
        editing_permission: 'user',
        is_test_user_mode: true,
        is_email_authentication_check: true,
        is_line_authentication: false,
      })
    })

    // メッセージ＆dirty 解除（保存ボタン disabled）
    expect(
      await screen.findByText('設定が保存されました。'),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })

  it('保存 非200：メッセージ「保存に失敗しました。」、dirty 継続（enabled のまま）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      status: 500,
      data: {},
    })

    renderWithRouter(<RootSetting />)
    const { combobox } = getControls()
    await screen.findByRole('combobox')

    await userEvent.selectOptions(combobox, 'root')
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    expect(await screen.findByText('保存に失敗しました。')).toBeInTheDocument()
    // まだ dirty（初期値は admin のまま） → 保存は enabled
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled()
  })

  it('保存 例外：メッセージ「エラーが発生しました。」、dirty 継続', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('POST fail'))

    renderWithRouter(<RootSetting />)
    const { cbLine } = getControls()
    await screen.findByRole('combobox')

    await userEvent.click(cbLine)
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    expect(
      await screen.findByText('エラーが発生しました。'),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled()
  })

  it('ローディング：保存中… の間は disabled、終了後に元の文言へ', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })

    // 保存を少し遅延させてローディングを観測
    let resolveSave: (v: any) => void
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolveSave = res
        }),
    )

    renderWithRouter(<RootSetting />)
    const { combobox } = getControls()
    await screen.findByRole('combobox')

    await userEvent.selectOptions(combobox, 'user')
    await userEvent.click(screen.getByRole('button', { name: '保存' }))

    // 保存中の表示＆ disabled
    const saving = screen.getByRole('button', { name: '保存中...' })
    expect(saving).toBeDisabled()

    // 終了
    resolveSave!({ status: 200, data: {} })
    expect(
      await screen.findByText('設定が保存されました。'),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })

  it('複数フィールドを変更 → 片方だけ戻すとまだ dirty、全て戻すと非 dirty', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        Config: {
          editing_permission: 'admin',
          is_test_user_mode: false,
          is_email_authentication_check: false,
          is_line_authentication: false,
        },
      },
    })

    renderWithRouter(<RootSetting />)
    const { combobox, cbTest } = getControls()
    await screen.findByRole('combobox')

    // 2箇所変更
    await userEvent.selectOptions(combobox, 'root')
    await userEvent.click(cbTest)
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled()

    // 1つだけ戻す → まだ dirty
    await userEvent.selectOptions(combobox, 'admin')
    expect(screen.getByRole('button', { name: '保存' })).toBeEnabled()

    // もう1つも戻す → 非 dirty
    await userEvent.click(cbTest)
    expect(screen.getByRole('button', { name: '保存' })).toBeDisabled()
  })
})
