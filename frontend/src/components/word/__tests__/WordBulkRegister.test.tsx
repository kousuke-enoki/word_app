// src/components/word/__tests__/WordBulkRegister.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import WordBulkRegister from '../WordBulkRegister'

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>)

/* --------- UI を薄くモックして安定化 --------- */
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div {...rest} data-testid="Card">
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
  Input: ({ ...rest }: any) => <input {...rest} />,
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({ children, ...rest }: any) => <button {...rest}>{children}</button>,
}))

/* --------- axios をモック --------- */
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

const typeInTextarea = async (val: string) => {
  const ta = screen.getByPlaceholderText(
    'Paste your English paragraph here...',
  ) as HTMLTextAreaElement
  await userEvent.clear(ta)
  await userEvent.type(ta, val)
  return ta
}

const clickExtract = async () => {
  const btn = screen.getByRole('button', { name: '抽出' })
  await userEvent.click(btn)
}

const clickRegister = async () => {
  const btn = screen.getByRole('button', { name: 'まとめて登録' })
  await userEvent.click(btn)
}

const clickReset = async () => {
  const btn = screen.getByRole('button', { name: '初期化' })
  await userEvent.click(btn)
}

beforeEach(() => {
  vi.resetAllMocks() // 実装も含めて完全リセット
})

describe('WordBulkRegister', () => {
  it('初期表示：見出し、文字カウンタ、ボタン状態', () => {
    renderWithRouter(<WordBulkRegister />)

    expect(
      screen.getByRole('heading', { name: '単語一括登録' }),
    ).toBeInTheDocument()
    expect(screen.getByText(/0\/5000/)).toBeInTheDocument()

    // 抽出は空入力なので disabled
    expect(screen.getByRole('button', { name: '抽出' })).toBeDisabled()

    // 初期化は enabled（仕様：loading && !tokens.length のみ disable）
    expect(screen.getByRole('button', { name: '初期化' })).toBeEnabled()

    // 候補カード/補助情報は出ていない
    expect(screen.queryByText(/候補:/)).not.toBeInTheDocument()
    expect(screen.queryByText('すでに登録済みの単語')).not.toBeInTheDocument()
    expect(
      screen.queryByText('データが存在しないため登録できない単語'),
    ).not.toBeInTheDocument()
  })

  it('抽出成功（候補あり）：候補・登録済み・未存在の表示、フィルタ、全選択/全解除、登録ボタン活性', async () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: {
        candidates: ['Apple', 'banana', 'Cherry'],
        registered: ['hello'],
        not_exists: ['xyz'],
      },
    })

    renderWithRouter(<WordBulkRegister />)

    await typeInTextarea('Some english paragraph.')
    expect(screen.getByRole('button', { name: '抽出' })).toBeEnabled()

    await clickExtract()

    // ローディング表示 → 成功メッセージ
    expect(
      await screen.findByText('抽出に成功しました（3 語）'),
    ).toBeInTheDocument()

    // 候補カードのツールバー類
    expect(screen.getByText('候補: 3')).toBeInTheDocument()
    expect(screen.getByText('選択: 0 / 200')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('検索...')).toBeInTheDocument()

    // 候補3件（ラベルテキストが見えている）
    expect(screen.getByText('Apple')).toBeInTheDocument()
    expect(screen.getByText('banana')).toBeInTheDocument()
    expect(screen.getByText('Cherry')).toBeInTheDocument()

    // 補助情報：登録済み/未存在
    expect(screen.getByText('すでに登録済みの単語')).toBeInTheDocument()
    expect(screen.getByText('hello')).toBeInTheDocument()
    expect(
      screen.getByText('データが存在しないため登録できない単語'),
    ).toBeInTheDocument()
    expect(screen.getByText('xyz')).toBeInTheDocument()

    // フィルタ： "an" で banana だけ残る
    await userEvent.type(screen.getByPlaceholderText('検索...'), 'an')
    expect(screen.queryByText('Apple')).not.toBeInTheDocument()
    expect(screen.getByText('banana')).toBeInTheDocument()
    expect(screen.queryByText('Cherry')).not.toBeInTheDocument()

    // フィルタ解除
    await userEvent.clear(screen.getByPlaceholderText('検索...'))
    expect(screen.getByText('Apple')).toBeInTheDocument()
    expect(screen.getByText('banana')).toBeInTheDocument()
    expect(screen.getByText('Cherry')).toBeInTheDocument()

    // 全選択 → カウント 3
    await userEvent.click(screen.getByRole('button', { name: '全選択' }))
    expect(screen.getByText('選択: 3 / 200')).toBeInTheDocument()

    // 全解除 → カウント 0
    await userEvent.click(screen.getByRole('button', { name: '全解除' }))
    expect(screen.getByText('選択: 0 / 200')).toBeInTheDocument()

    // 1件手動選択（ラベルクリックでチェック切替）
    await userEvent.click(screen.getByText('Apple'))
    expect(screen.getByText('選択: 1 / 200')).toBeInTheDocument()

    // 選択あり → 「まとめて登録」活性
    expect(screen.getByRole('button', { name: 'まとめて登録' })).toBeEnabled()
  })

  it('抽出成功（0件）：適切なメッセージ・候補カード非表示', async () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: [], registered: [], not_exists: [] },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('text')
    await clickExtract()

    expect(
      await screen.findByText('登録できる単語がありませんでした。'),
    ).toBeInTheDocument()
    expect(screen.queryByText(/候補:/)).not.toBeInTheDocument()
  })

  it('抽出失敗：エラーメッセージ', async () => {
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('500'))

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('text')
    await clickExtract()

    expect(await screen.findByText('抽出に失敗しました')).toBeInTheDocument()
  })

  it('登録：success のみ', async () => {
    // 抽出
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: ['dog', 'cat'], registered: [], not_exists: [] },
    })
    // 登録
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { success: ['dog', 'cat'] },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('animals')
    await clickExtract()
    await screen.findByText('抽出に成功しました（2 語）')

    // 2件選択
    await userEvent.click(screen.getByRole('button', { name: '全選択' }))
    await clickRegister()

    // POST 内容
    await waitFor(() => {
      const calls = (axiosInstance.post as any).mock.calls
      const [path, payload] = calls.at(-1)

      expect(path).toBe('/words/bulk_register')
      expect(payload.words).toEqual(expect.arrayContaining(['cat', 'dog']))
      expect(payload.words).toHaveLength(2) // 数も担保
    })

    // 成功メッセージ（詳細表示）
    expect(await screen.findByText(/✅ 登録成功（2件）/)).toBeInTheDocument()
    expect(screen.getByText('dog')).toBeInTheDocument()
    expect(screen.getByText('cat')).toBeInTheDocument()
  })

  it('登録：failed のみ', async () => {
    // 抽出
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: ['one'], registered: [], not_exists: [] },
    })
    // 登録
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { failed: [{ word: 'one', reason: 'duplicate' }] },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('numbers')
    await clickExtract()
    await screen.findByText('抽出に成功しました（1 語）')

    await userEvent.click(screen.getByText('one'))
    await clickRegister()

    expect(await screen.findByText(/❌ 登録失敗（1件）/)).toBeInTheDocument()
    expect(
      screen.getByText(/one.*duplicate|duplicate.*one/),
    ).toBeInTheDocument()
  })

  it('登録：success と failed 混在', async () => {
    // 抽出
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: ['a', 'b', 'c'], registered: [], not_exists: [] },
    })
    // 登録
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: {
        success: ['a'],
        failed: [
          { word: 'b', reason: 'not_found' },
          { word: 'c', reason: 'already_exists' },
        ],
      },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('abc')
    await clickExtract()
    await screen.findByText('抽出に成功しました（3 語）')

    await userEvent.click(screen.getByRole('button', { name: '全選択' }))
    await clickRegister()

    // 成功メッセージ（詳細表示）
    expect(await screen.findByText(/✅ 登録成功（1件）/)).toBeInTheDocument()
    expect(screen.getByText('a')).toBeInTheDocument()

    // 失敗メッセージ（詳細表示）
    expect(screen.getByText(/❌ 登録失敗（2件）/)).toBeInTheDocument()
    expect(screen.getByText('b')).toBeInTheDocument()
    expect(screen.getByText('c')).toBeInTheDocument()
  })

  it('登録失敗：エラーメッセージ', async () => {
    // 抽出
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: ['x'], registered: [], not_exists: [] },
    })
    // 登録失敗
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('500'))

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('x')
    await clickExtract()
    await screen.findByText('抽出に成功しました（1 語）')

    await userEvent.click(screen.getByRole('checkbox', { name: 'x' }))
    await clickRegister()

    expect(await screen.findByText('登録に失敗しました')).toBeInTheDocument()
  })

  it('選択 0 件では登録ボタン disabled → クリックできない', async () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: ['x', 'y'], registered: [], not_exists: [] },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('xy')
    await clickExtract()
    await screen.findByText('抽出に成功しました（2 語）')

    const regBtn = screen.getByRole('button', { name: 'まとめて登録' })
    expect(regBtn).toBeDisabled()
  })

  it('選択 201 件（>200）だと early return：axios へ登録 POST されない', async () => {
    const words = Array.from({ length: 201 }, (_, i) => `w${i + 1}`)
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { candidates: words, registered: [], not_exists: [] },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('many')
    await clickExtract()
    await screen.findByText(`抽出に成功しました（${words.length} 語）`)

    // 全選択 → 201 件
    await userEvent.click(screen.getByRole('button', { name: '全選択' }))
    expect(screen.getByText('選択: 201 / 200')).toBeInTheDocument()

    await clickRegister()

    // 1回目は /bulk_tokenize、2回目の /bulk_register が呼ばれていないことを確認
    const calls = (axiosInstance.post as any).mock.calls.map((c: any[]) => c[0])
    expect(
      calls.filter((p: string) => p === '/words/bulk_register').length,
    ).toBe(0)
  })

  it('初期化：全てクリアされる', async () => {
    // 抽出済み状態を作る
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: {
        candidates: ['apple'],
        registered: ['done'],
        not_exists: ['none'],
      },
    })

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('text')
    await clickExtract()
    await screen.findByText('抽出に成功しました（1 語）')
    expect(screen.getByText('候補: 1')).toBeInTheDocument()
    expect(screen.getByText('すでに登録済みの単語')).toBeInTheDocument()
    expect(
      screen.getByText('データが存在しないため登録できない単語'),
    ).toBeInTheDocument()

    // 初期化
    await clickReset()

    // クリア確認
    expect(screen.getByText(/0\/5000/)).toBeInTheDocument()
    expect(screen.queryByText(/候補:/)).not.toBeInTheDocument()
    expect(screen.queryByText('すでに登録済みの単語')).not.toBeInTheDocument()
    expect(
      screen.queryByText('データが存在しないため登録できない単語'),
    ).not.toBeInTheDocument()
  })

  it('ローディング表示：抽出中… / 登録中…', async () => {
    // 抽出は少し遅らせる
    // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
    let resolveExtract: Function
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolveExtract = res
        }),
    )

    renderWithRouter(<WordBulkRegister />)
    await typeInTextarea('delayed')
    await userEvent.click(screen.getByRole('button', { name: '抽出' }))

    // 抽出中…
    expect(screen.getByRole('button', { name: '抽出中…' })).toBeDisabled()

    // 抽出 resolve
    resolveExtract!({
      data: { candidates: ['a'], registered: [], not_exists: [] },
    })
    expect(
      await screen.findByText('抽出に成功しました（1 語）'),
    ).toBeInTheDocument()

    // 登録も遅延させる
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) =>
          setTimeout(() => res({ data: { success: ['a'] } }), 50),
        ),
    )

    await userEvent.click(screen.getByText('a'))
    await userEvent.click(screen.getByRole('button', { name: 'まとめて登録' }))
    expect(screen.getByRole('button', { name: '登録中…' })).toBeDisabled()

    expect(await screen.findByText(/✅ 登録成功（1件）/)).toBeInTheDocument()
    expect(screen.getByText('a')).toBeInTheDocument()
  })
})
