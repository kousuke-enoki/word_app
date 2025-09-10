// src/components/word/__tests__/WordEdit.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { renderWithClient } from '@/__tests__/testUtils'

import WordEdit from '../WordEdit'

/* ---- UI を薄くモック（安定化） ---- */
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
  Input: ({ ...rest }: any) => <input {...rest} />,
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({ children, ...rest }: any) => <button {...rest}>{children}</button>,
}))

/* ---- 品詞データ固定 ---- */
vi.mock('@/service/word/GetPartOfSpeech', () => {
  return {
    getPartOfSpeech: [
      { id: 1, name: '名詞' },
      { id: 2, name: '動詞' },
      { id: 3, name: '形容詞' },
      { id: 4, name: '副詞' },
    ],
  }
})

/* ---- axios モック ---- */
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn(), put: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

/* ---- ルーター：useParams/useNavigate をモック ---- */
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return {
    ...actual,
    useParams: () => ({ id: '10' }),
    useNavigate: () => navigateMock,
  }
})

beforeEach(() => {
  vi.resetAllMocks() // 実装も含めて完全リセット
})

/* ヘルパー：描画 */
const renderEdit = async () => {
  return renderWithClient(<WordEdit />)
}

/* 初期フェッチ用のダミー */
const fetchedWord = {
  id: 10,
  name: 'apple',
  wordInfos: [
    {
      id: 1,
      partOfSpeechId: 1, // 名詞
      japaneseMeans: [{ id: 11, name: 'りんご' }],
    },
    {
      id: 2,
      partOfSpeechId: 2, // 動詞
      japaneseMeans: [{ id: 21, name: 'はしる' }],
    },
  ],
}

/* ========= テスト本体 ========= */
describe('WordEdit', () => {
  it('読み込み中 → 完了でフォームが出る', async () => {
    // get を遅延して "読み込み中…" を確認
    let resolveGet!: (v: unknown) => void
    ;(axiosInstance.get as any).mockImplementationOnce(
      () => new Promise((res) => (resolveGet = res)),
    )
    renderEdit()
    // ローディング表示
    expect(screen.getByText('読み込み中…')).toBeInTheDocument()

    // 解決 → フォーム表示
    resolveGet({ data: fetchedWord })
    expect(
      await screen.findByRole('heading', { name: '単語更新' }),
    ).toBeInTheDocument()

    // 単語名 input は "apple"（ラベル関連付けがないので value で取得）
    expect(screen.getByDisplayValue('apple')).toBeInTheDocument()
  })

  it('取得エラー時：エラーメッセージと「再取得」→ 正常表示に復帰', async () => {
    ;(axiosInstance.get as any)
      .mockRejectedValueOnce(new Error('500')) // 最初は失敗
      .mockResolvedValueOnce({ data: fetchedWord }) // 再取得で成功

    renderEdit()

    // エラー表示
    expect(
      await screen.findByText('単語情報の取得中にエラーが発生しました。'),
    ).toBeInTheDocument()

    // 再取得で回復
    await userEvent.click(screen.getByRole('button', { name: '再取得' }))
    expect(
      await screen.findByRole('heading', { name: '単語更新' }),
    ).toBeInTheDocument()
  })

  it('データなし：null が返ると「データが存在しません。」', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: null })
    renderEdit()
    expect(
      await screen.findByText('データが存在しません。'),
    ).toBeInTheDocument()
  })

  it('品詞オプション：他インデックスの選択肢は除外される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: fetchedWord })
    renderEdit()
    await screen.findByRole('heading', { name: '単語更新' })

    // 先頭の品詞セレクト
    const selects = screen.getAllByRole('combobox')
    const firstSelect = selects[0]
    const opts = within(firstSelect)
      .getAllByRole('option')
      .map((o) => o.textContent)

    // index=0 では他インデックス(index=1)の "動詞" は除外される想定
    expect(opts).not.toContain('動詞')
    // 未選択用 "選択してください" と自分以外の残りは含まれる
    expect(opts).toEqual(
      expect.arrayContaining(['選択してください', '名詞', '形容詞', '副詞']),
    )
  })

  it('バリデーション：名前/品詞/日本語訳のエラーを表示し、PUTは呼ばれない', async () => {
    // fetched をあえて NG っぽい値に変更してから検証してもOKだが、
    // ここでは入力で NG に寄せる
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: fetchedWord })
    renderEdit()
    await screen.findByRole('heading', { name: '単語更新' })

    // 単語名に数字（NG）
    const nameInput = screen.getByDisplayValue('apple') as HTMLInputElement
    await userEvent.clear(nameInput)
    await userEvent.type(nameInput, '123')

    // 品詞を 0 に変更（NG）
    const selects = screen.getAllByRole('combobox')
    await userEvent.selectOptions(selects[0], '0')

    // 日本語訳にアルファベット（NG）
    const jpInput = screen.getByDisplayValue('りんご') as HTMLInputElement
    await userEvent.clear(jpInput)
    await userEvent.type(jpInput, 'abc')

    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    // エラーメッセージが表示
    expect(
      await screen.findByText('単語名は半角アルファベットのみ入力できます。'),
    ).toBeInTheDocument()
    expect(screen.getByText('品詞を選択してください。')).toBeInTheDocument()
    expect(
      screen.getByText(
        '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      ),
    ).toBeInTheDocument()

    expect(axiosInstance.put).not.toHaveBeenCalled()
  })

  it('更新成功：PUT して 詳細へ navigate（成功メッセージ付き）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: fetchedWord })
    ;(axiosInstance.put as any).mockResolvedValueOnce({
      data: { name: 'banana' },
    })

    renderEdit()
    await screen.findByRole('heading', { name: '単語更新' })

    // 値を更新
    const nameInput = screen.getByDisplayValue('apple') as HTMLInputElement
    await userEvent.clear(nameInput)
    await userEvent.type(nameInput, 'banana')

    const selects = screen.getAllByRole('combobox')
    // 先頭の品詞を 3(形容詞) に
    await userEvent.selectOptions(selects[0], '3')

    // 日本語訳も変更
    const meansInputs = screen
      .getAllByRole('textbox')
      .filter((el) => (el as HTMLInputElement).value !== 'banana')
    await userEvent.clear(meansInputs[0] as HTMLInputElement)
    await userEvent.type(meansInputs[0] as HTMLInputElement, 'バナナ')

    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    await waitFor(() => {
      // PUT の payload を軽く検査
      const [path, payload] = (axiosInstance.put as any).mock.calls.at(0)
      expect(path).toBe('/words/10')
      expect(payload.name).toBe('banana')
      expect(payload.wordInfos[0].partOfSpeechId).toBe(3)
      expect(payload.wordInfos[0].japaneseMeans[0].name).toBe('バナナ')
    })

    expect(navigateMock).toHaveBeenCalledWith('/words/10', {
      state: { successMessage: 'bananaが正常に更新されました！' },
    })
  })

  it('更新失敗：PUT 失敗でエラーメッセージ、navigate されない', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: fetchedWord })
    ;(axiosInstance.put as any).mockRejectedValueOnce(new Error('500'))

    renderEdit()
    await screen.findByRole('heading', { name: '単語更新' })

    // なんでも良いので submit
    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    expect(
      await screen.findByText('単語情報の更新中にエラーが発生しました。'),
    ).toBeInTheDocument()
    expect(navigateMock).not.toHaveBeenCalled()
  })

  it('「単語詳細に戻る」：/words/:id へ遷移', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: fetchedWord })
    renderEdit()
    await screen.findByRole('heading', { name: '単語更新' })

    await userEvent.click(
      screen.getByRole('button', { name: '単語詳細に戻る' }),
    )
    expect(navigateMock).toHaveBeenCalledWith('/words/10')
  })
})
