/* eslint-disable @typescript-eslint/no-explicit-any */
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route,Routes } from 'react-router-dom'
import {
  afterEach,
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest'

import { queryClient,renderWithClient } from '@/__tests__/testUtils'

vi.mock('@/axiosConfig', () => ({
  default: {
    get : vi.fn(),
    put : vi.fn(),
  },
}))
import axiosInstance from '@/axiosConfig'

const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>(
    'react-router-dom',
  )
  return { ...actual, useNavigate: () => navigateMock }
})

import WordEdit, {
  WordForUpdate,
} from '../WordEdit'

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
})

afterEach(() => {
  queryClient?.clear()
})

/* ------------- テストデータ ---------------------- */
const stubWord: WordForUpdate = {
  id: 1,
  name: 'apple',
  wordInfos: [
    {
      id: 11,
      partOfSpeechId: 1,        // 名詞
      japaneseMeans: [{ id: 111, name: 'りんご' }],
    },
  ],
}

/* =================================================
 *                    TESTS
 * ================================================= */
describe('WordEdit Component', () => {
  /** 1. 取得中の表示 */
  it('ローディング中は「読み込み中…」が出る', async () => {
    // まだ解決しない Promise を返して “ロード中” を再現
    ;(axiosInstance.get as any).mockImplementation(
      () => new Promise(() => {}),
    )

    renderWithClient(
      <MemoryRouter initialEntries={['/words/1/edit']}>
        <Routes>
          <Route path="/words/:id/edit" element={<WordEdit />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText('読み込み中...'))
  })

  /** 2. 正常取得するとフォームが初期値で埋まる */
  it('単語データ取得後にフォームへ反映される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: stubWord })

    renderWithClient(
      <MemoryRouter initialEntries={['/words/1/edit']}>
        <Routes>
          <Route path="/words/:id/edit" element={<WordEdit />} />
        </Routes>
      </MemoryRouter>,
    );

    // 単語名が入っていること
    expect(
      await screen.findByDisplayValue('apple'))

    // 日本語訳が入っていること
    expect(screen.getByDisplayValue('りんご'))
  })

  /** 3. バリデーション : 不正値を入れるとエラーメッセージ */
  it('無効な入力は送信されずフィールドエラーが出る', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: stubWord })

    renderWithClient(
      <MemoryRouter initialEntries={['/words/1/edit']}>
        <Routes>
          <Route path="/words/:id/edit" element={<WordEdit />} />
        </Routes>
      </MemoryRouter>,
    );

    // 単語名を“数字”に変更 → 半角英字制約に違反
    await userEvent.clear(await screen.findByDisplayValue('apple'))
    // await userEvent.type(screen.getByRole('textbox'), '1234')
    await userEvent.type(
      screen.getByLabelText('単語名:'),   // ← ラベルで単独指定
      '1234',
    )

    // 送信
    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    expect(
      await screen.findByText('単語名は半角アルファベットのみ入力できます。'))

    // PUT は呼ばれていない
    expect(axiosInstance.put).not.toHaveBeenCalled()
  })

  /** 4. 更新成功フロー  */
  it('正しい値で送信すると PUT され /words/1 へ遷移', async () => {
    // ① 取得 → フォーム初期化
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: stubWord })

    // ② 更新 API 成功時のレスポンス
    ;(axiosInstance.put as any).mockResolvedValueOnce({ data: { name: 'apple' } })

    renderWithClient(
      <MemoryRouter initialEntries={['/words/1/edit']}>
        <Routes>
          <Route path="/words/:id/edit" element={<WordEdit />} />
        </Routes>
      </MemoryRouter>,
    );

    // 何も変更せずそのまま送信
    await userEvent.click(
      await screen.findByRole('button', { name: '単語を更新' }),
    )

    // PUT が呼ばれたこと
    await waitFor(() =>
      expect(axiosInstance.put).toHaveBeenCalledWith(
        '/words/1',
        stubWord,
      ),
    )

    // navigate('/') が呼ばれるまで待機
    expect(navigateMock).toHaveBeenCalledWith('/words/1', {
      state: { successMessage: 'appleが正常に更新されました！' },
    })
  })

  /** 5. 更新失敗フロー */
  it('更新 API が失敗するとエラーバナーが出る', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: stubWord })
    ;(axiosInstance.put as any).mockRejectedValueOnce(new Error('500'))

    renderWithClient(
      <MemoryRouter initialEntries={['/words/1/edit']}>
        <Routes>
          <Route path="/words/:id/edit" element={<WordEdit />} />
        </Routes>
      </MemoryRouter>,
    );

    await userEvent.click(
      await screen.findByRole('button', { name: '単語を更新' }),
    )

    expect(
      await screen.findByText('単語情報の更新中にエラーが発生しました。'))
    expect(navigateMock).not.toHaveBeenCalled()
  })
})
