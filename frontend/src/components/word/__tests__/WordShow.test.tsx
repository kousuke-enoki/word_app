// src/components/word/__tests__/WordShow.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, within } from '@testing-library/react'
import { fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import WordShow from '../WordShow'

/* ---- UI を薄くモック（安定化） ---- */
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({ children, ...rest }: any) => <button {...rest}>{children}</button>,
}))

/* ---- RegisterToggle を最小化（トグルボタン1つ） ---- */
vi.mock('@/components/common/RegisterToggle', () => ({
  RegisterToggle: ({ isRegistered, onToggle }: any) => (
    <button onClick={onToggle}>{isRegistered ? '解除' : '登録'}</button>
  ),
}))

/* ---- 品詞データ固定 ---- */
vi.mock('@/service/word/GetPartOfSpeech', () => ({
  getPartOfSpeech: [
    { id: 1, name: '名詞' },
    { id: 2, name: '動詞' },
    { id: 3, name: '形容詞' },
  ],
}))

/* ---- API をモック ---- */
vi.mock('@/axiosConfig', () => ({ default: { get: vi.fn() } }))
vi.mock('@/service/word/RegisterWord', () => ({ registerWord: vi.fn() }))
vi.mock('@/service/word/SaveMemo', () => ({ saveMemo: vi.fn() }))
vi.mock('@/service/word/DeleteWord', () => ({ deleteWord: vi.fn() }))

import axiosInstance from '@/axiosConfig'
import { deleteWord } from '@/service/word/DeleteWord'
import { registerWord } from '@/service/word/RegisterWord'
import { saveMemo } from '@/service/word/SaveMemo'

/* ---- 便利ルート（遷移と state を観測） ---- */
const ShowState: React.FC = () => {
  const loc = useLocation()
  return (
    <div>
      <div>一覧ページ</div>
      <pre data-testid="state">{JSON.stringify(loc.state)}</pre>
    </div>
  )
}
const EditSentinel: React.FC = () => <div>編集ページ</div>

/* ---- render ヘルパ（/words/:id で WordShow を表示） ---- */
const renderShow = (initialPath = '/words/1', state?: any) =>
  render(
    <MemoryRouter initialEntries={[{ pathname: initialPath, state }]}>
      <Routes>
        <Route path="/words/:id" element={<WordShow />} />
        <Route path="/words/edit/:id" element={<EditSentinel />} />
        <Route path="/words" element={<ShowState />} />
      </Routes>
    </MemoryRouter>,
  )

/* ---- テスト用データ ---- */
const baseWord = {
  id: 1,
  name: 'apple',
  wordInfos: [
    {
      id: 10,
      partOfSpeechId: 1,
      japaneseMeans: [{ name: 'りんご' }, { name: '林檎' }],
    },
  ],
  registrationCount: 2,
  isRegistered: false,
  attentionLevel: 1,
  quizCount: 3,
  correctCount: 2,
  memo: 'おいしい果物',
}

beforeEach(() => {
  vi.useRealTimers()
  vi.resetAllMocks()
  vi.spyOn(window, 'alert').mockImplementation(() => {})
  vi.spyOn(window, 'confirm').mockImplementation(() => true)
})
// afterEach で念のため実タイマーへ戻す（漏れ防止）
afterEach(() => {
  vi.useRealTimers()
  vi.clearAllTimers()
})
describe('WordShow', () => {
  it('(1) 初期表示：取得成功（値/バッジ/メモ/ボタン）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    renderShow()

    // タイトル
    expect(
      await screen.findByRole('heading', { name: 'apple' }),
    ).toBeInTheDocument()

    // 日本語訳と品詞
    // expect(screen.getByText('日本語訳: りんご, 林檎')).toBeInTheDocument()
    const infoCard = screen.getAllByTestId('Card')[0]

    // 「日本語訳」は <span> を起点に親 <p> を検証
    const jpLabel = within(infoCard).getByText('日本語訳:', {
      selector: 'span',
    })
    expect(jpLabel.parentElement).toHaveTextContent(
      /日本語訳:\s*りんご,\s*林檎/,
    )

    // 「品詞」は <p> に限定して完全一致
    const posP = within(infoCard).getByText(/^品詞:\s*名詞$/, { selector: 'p' })
    expect(posP).toBeInTheDocument()

    // バッジ群
    expect(screen.getByText('全登録数: 2')).toBeInTheDocument()
    expect(screen.getByText('注意レベル: 1')).toBeInTheDocument()
    expect(screen.getByText('テスト回数: 3')).toBeInTheDocument()
    expect(screen.getByText('チェック回数: 2')).toBeInTheDocument()

    // メモ初期値
    const ta = screen.getByRole('textbox') as HTMLTextAreaElement
    expect(ta.value).toBe('おいしい果物')

    // 登録トグル表示（未登録 → 登録ボタン）
    expect(screen.getByRole('button', { name: '登録' })).toBeInTheDocument()
  })

  it('(2) 初期表示：取得失敗 → alert + No word', async () => {
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('boom'))
    renderShow()

    expect(
      await screen.findByText('No word details found.'),
    ).toBeInTheDocument()
    expect(window.alert).toHaveBeenCalledWith(
      '単語情報の取得中にエラーが発生しました。',
    )
  })

  it('(3) location.state の successMessage が初期表示される', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    renderShow('/words/1', { successMessage: '編集で保存しました！' })

    await screen.findByRole('heading', { name: 'apple' })
    expect(screen.getByText('編集で保存しました！')).toBeInTheDocument()
  })

  it('(4) 登録：成功 → カウント/表示更新', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(registerWord as any).mockResolvedValueOnce({
      id: 1,
      isRegistered: true,
      registrationCount: 3,
    })

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })
    fireEvent.click(screen.getByRole('button', { name: '登録' }))
    expect(
      await screen.findByRole('button', { name: '解除' }),
    ).toBeInTheDocument() // メッセージが DOM に現れるのを “実タイマー” で待つ
    expect(await screen.findByText('登録しました。')).toBeInTheDocument()
  })

  it('(5) 登録：失敗 → alert、表示は変わらない', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(registerWord as any).mockRejectedValueOnce(new Error('nope'))

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    await userEvent.click(screen.getByRole('button', { name: '登録' }))
    expect(window.alert).toHaveBeenCalledWith(
      '単語の登録中にエラーが発生しました。',
    )

    // isRegistered 変わらず / カウントも 2 のまま
    expect(screen.getByRole('button', { name: '登録' })).toBeInTheDocument()
    expect(screen.getByText('全登録数: 2')).toBeInTheDocument()
  })

  it('(6) メモ保存：成功 → API 引数/メッセージ表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(saveMemo as any).mockResolvedValueOnce({ ok: true })

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    fireEvent.change(screen.getByRole('textbox'), {
      target: { value: '新しいメモ' },
    })
    fireEvent.click(screen.getByRole('button', { name: '保存する' }))

    expect(saveMemo).toHaveBeenCalledWith(1, '新しいメモ')
    expect(await screen.findByText('メモを保存しました！')).toBeInTheDocument()
  })

  it('(7) メモ保存：失敗 → alert', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(saveMemo as any).mockRejectedValueOnce(new Error('x'))

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    await userEvent.click(screen.getByRole('button', { name: '保存する' }))
    expect(window.alert).toHaveBeenCalledWith(
      'メモの保存中にエラーが発生しました。',
    )
  })

  it('(8) 削除：confirm キャンセル → API 呼ばれない/遷移しない', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(deleteWord as any).mockResolvedValueOnce({ ok: true })
    ;(window.confirm as any).mockReturnValueOnce(false)

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    await userEvent.click(screen.getByRole('button', { name: '削除する' }))
    expect(deleteWord).not.toHaveBeenCalled()
    // 依然 WordShow の内容が見える
    expect(screen.getByRole('heading', { name: 'apple' })).toBeInTheDocument()
  })

  it('(9) 削除：成功 → メッセージ表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(deleteWord as any).mockResolvedValueOnce({ ok: true })

    renderShow('/words/1', { page: 3 })
    await screen.findByRole('heading', { name: 'apple' })
    fireEvent.click(screen.getByRole('button', { name: '削除する' }))
    expect(await screen.findByText('単語を削除しました。')).toBeInTheDocument()
  })

  it('(10) 削除：失敗 → メッセージ表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    ;(deleteWord as any).mockRejectedValueOnce(new Error('ng'))

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })
    fireEvent.click(screen.getByRole('button', { name: '削除する' }))
    // 実タイマーで “出現” を待つ
    expect(
      await screen.findByText('単語の削除に失敗しました。'),
    ).toBeInTheDocument()
  })

  it('(11) 編集ボタン → /words/edit/:id に遷移', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    await userEvent.click(screen.getByRole('button', { name: '編集する' }))
    expect(await screen.findByText('編集ページ')).toBeInTheDocument()
  })

  it('(12) 一覧に戻る → /words へ遷移（state 既定値 or 渡されたものを使用）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: baseWord })
    const init = {
      search: 'zz',
      sortBy: 'registrationCount',
      order: 'desc',
      page: 2,
      limit: 50,
    }
    renderShow('/words/1', init)
    await screen.findByRole('heading', { name: 'apple' })

    await userEvent.click(screen.getByRole('button', { name: '一覧に戻る' }))

    expect(await screen.findByText('一覧ページ')).toBeInTheDocument()
    const passed = JSON.parse(
      (screen.getByTestId('state') as HTMLElement).textContent || '{}',
    )
    expect(passed).toMatchObject(init)
  })

  it('(13) 品詞ID 未定義なら「未定義」表示', async () => {
    const w = {
      ...baseWord,
      wordInfos: [
        { id: 99, partOfSpeechId: 999, japaneseMeans: [{ name: '？？' }] },
      ],
    }
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: w })

    renderShow()
    await screen.findByRole('heading', { name: 'apple' })

    expect(screen.getByText('品詞: 未定義')).toBeInTheDocument()
  })

  it('(14) Loading 表示 → 解消', async () => {
    // resolve を少し遅らせて Loading が出ることを確認
    ;(axiosInstance.get as any).mockImplementationOnce(
      () =>
        new Promise((resolve) =>
          setTimeout(() => resolve({ data: baseWord }), 10),
        ),
    )

    renderShow()
    // 先に Loading...
    expect(screen.getByText('Loading...')).toBeInTheDocument()

    // その後解消
    expect(
      await screen.findByRole('heading', { name: 'apple' }),
    ).toBeInTheDocument()
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument()
  })
})
