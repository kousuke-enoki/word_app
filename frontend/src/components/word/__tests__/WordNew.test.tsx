/* eslint-disable @typescript-eslint/no-explicit-any */
import { fireEvent, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { renderWithClient } from '@/__tests__/testUtils'

import WordNew from '../WordNew'

/* axios をモック */
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

/* 品詞データを固定（テストが壊れないように） */
vi.mock('@/service/word/GetPartOfSpeech', () => {
  return {
    getPartOfSpeech: [
      { id: 1, name: '名詞' },
      { id: 2, name: '動詞' },
      { id: 3, name: '形容詞' },
    ],
  }
})

/* useNavigate をモック */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  // バリデーションの alert を無害化
  vi.spyOn(window, 'alert').mockImplementation(() => {})
})
afterEach(() => {
  vi.restoreAllMocks()
})

const renderWordNew = () =>
  renderWithClient(<WordNew />, {
    router: {
      initialEntries: ['/words/new'],
      routes: [{ path: '/words/new', element: <WordNew /> }],
    },
  })

/* ---------------- テスト本体 ----------------------- */
describe('WordNew Component', () => {
  /* 1) 初期描画 */
  it('フォームは空で表示される', () => {
    renderWordNew()

    // 見出し
    expect(
      screen.getByRole('heading', { name: '単語登録' }),
    ).toBeInTheDocument()

    // 単語名 input は空（placeholder=example）
    const nameInput = screen.getByPlaceholderText('example') as HTMLInputElement
    expect(nameInput).toHaveValue('')

    // 品詞は “選択してください”(value 0)
    const posSelect = screen.getByRole('combobox') as HTMLSelectElement
    expect(posSelect).toHaveValue('0')

    // 日本語訳 input も空（placeholder=意味）
    const jpInput = screen.getByPlaceholderText('意味') as HTMLInputElement
    expect(jpInput).toHaveValue('')
  })

  /* 2) バリデーション */
  it('無効な入力は送信されずエラーを表示', async () => {
    renderWordNew()

    // 数字を入力 → ハンドラで拒否され、値は '' のまま
    const nameInput = screen.getByPlaceholderText('example') as HTMLInputElement
    await userEvent.type(nameInput, '1234')
    expect(nameInput).toHaveValue('')

    // form を送信（ブラウザネイティブ検証は避ける）
    const form = screen.getByRole('form', { name: 'word-create-form' })
    fireEvent.submit(form)

    // エラーが表示され、POST は呼ばれない
    expect(
      await screen.findByText('単語名は半角アルファベットのみ入力できます。'),
    ).toBeInTheDocument()

    // 追加のエラー（品詞/日本語訳）も出ているはずだが、最低限の主張のみ検証
    expect(axiosInstance.post).not.toHaveBeenCalled()
  })

  /* 3) 登録成功フロー */
  it('正しい値で送信すると POST → /words/:id へ遷移', async () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { id: 99, name: 'apple' },
    })

    renderWordNew()

    // --- フォーム入力 ---
    await userEvent.type(screen.getByPlaceholderText('example'), 'apple')
    await userEvent.selectOptions(screen.getByRole('combobox'), '1') // 名詞
    await userEvent.type(screen.getByPlaceholderText('意味'), 'りんご')

    await userEvent.click(
      screen.getByRole('button', { name: '単語を登録する' }),
    )

    // POST されたこと
    await waitFor(() =>
      expect(axiosInstance.post).toHaveBeenCalledWith('/words/new', {
        name: 'apple',
        wordInfos: [{ partOfSpeechId: 1, japaneseMeans: [{ name: 'りんご' }] }],
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

    renderWordNew()

    await userEvent.type(screen.getByPlaceholderText('example'), 'apple')
    await userEvent.selectOptions(screen.getByRole('combobox'), '1')
    await userEvent.type(screen.getByPlaceholderText('意味'), 'りんご')
    await userEvent.click(
      screen.getByRole('button', { name: '単語を登録する' }),
    )

    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'),
    ).toBeInTheDocument()
    expect(navigateMock).not.toHaveBeenCalled()
  })
})
