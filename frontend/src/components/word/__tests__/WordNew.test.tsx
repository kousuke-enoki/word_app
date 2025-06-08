/* eslint-disable @typescript-eslint/no-explicit-any */
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { fireEvent, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import WordNew from '../WordNew'
import { renderWithClient, queryClient } from '@/__tests__/testUtils'

/* axios モック */
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

/* useNavigate モック */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>(
    'react-router-dom',
  )
  return { ...actual, useNavigate: () => navigateMock }
})

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
})
afterEach(() => queryClient.clear())

/* ---------------- テスト本体 ----------------------- */
describe('WordNew Component', () => {
  /* 1) 初期描画 */
  it('フォームは空で表示される', () => {
    renderWithClient(
      <MemoryRouter initialEntries={['/words/new']}>
        <Routes>
          <Route path="/words/new" element={<WordNew />} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByRole('heading', { name: '単語登録フォーム' }))
    // 単語名 input 空
    expect(screen.getByLabelText('単語名:'),(''))
    // 品詞は “選択してください” (value 0)
    expect(screen.getByRole('combobox'),('0'))
  })

  /* 2) バリデーション */
  it('無効な入力は送信されずエラーを表示', async () => {
    renderWithClient(
      <MemoryRouter initialEntries={['/words/new']}>
        <Routes>
          <Route path="/words/new" element={<WordNew />} />
        </Routes>
      </MemoryRouter>,
    )
  
    // 数字を入力 → input の中身は '' のまま
    await userEvent.type(screen.getByLabelText('単語名:'), '1234')
  
    // form 要素を直接 submit させてブラウザ検証を回避
    const form = screen.getByRole('form')   // name フィルタを外す
    fireEvent.submit(form)
  
    expect(
      await screen.findByText('単語名は半角アルファベットのみ入力できます。'))
    expect(axiosInstance.post).not.toHaveBeenCalled()
  })

  /* 3) 登録成功フロー */
  it('正しい値で送信すると POST → /words/:id へ遷移', async () => {
    // axios.post 成功レスポンス
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { id: 99, name: 'apple' },
    })

    renderWithClient(
      <MemoryRouter initialEntries={['/words/new']}>
        <Routes>
          <Route path="/words/new" element={<WordNew />} />
        </Routes>
      </MemoryRouter>,
    )

    // --- フォーム入力 ---
    await userEvent.type(screen.getByLabelText('単語名:'), 'apple')
    await userEvent.selectOptions(screen.getByRole('combobox'), '1') // 名詞など
    await userEvent.type(screen.getByLabelText('日本語訳:'), 'りんご')
    await userEvent.click(screen.getByRole('button', { name: '単語を登録' }))

    // POST されたこと
    await waitFor(() =>
      expect(axiosInstance.post).toHaveBeenCalledWith('/words/new', {
        name: 'apple',
        wordInfos: [
          { partOfSpeechId: 1, japaneseMeans: [{ name: 'りんご' }] },
        ],
      }),
    )

    // navigate が呼ばれる
    expect(navigateMock).toHaveBeenCalledWith('/words/99', {
      state: { successMessage: 'appleが正常に登録されました！' },
    })
  })

  /* 4) 登録失敗フロー */
  it('API 失敗時はエラーメッセージを表示', async () => {
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('500'))

    renderWithClient(
      <MemoryRouter initialEntries={['/words/new']}>
        <Routes>
          <Route path="/words/new" element={<WordNew />} />
        </Routes>
      </MemoryRouter>,
    )

    await userEvent.type(screen.getByLabelText('単語名:'), 'apple')
    await userEvent.selectOptions(screen.getByRole('combobox'), '1')
    await userEvent.type(screen.getByLabelText('日本語訳:'), 'りんご')
    await userEvent.click(screen.getByRole('button', { name: '単語を登録' }))

    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'))
    expect(navigateMock).not.toHaveBeenCalled()
  })
})
