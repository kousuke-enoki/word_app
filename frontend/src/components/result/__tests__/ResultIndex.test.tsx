// src/components/result/__tests__/ResultIndex.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { fireEvent, render, screen } from '@testing-library/react'
import { within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ResultIndex from '../ResultIndex'

/* ------------ ライブラリ & 依存モック ------------ */
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

// useNavigate をスパイ
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

// Router 依存のあるコンポーネントは薄くモック（描画安定化）
vi.mock('../../common/PageBottomNav', () => ({
  default: (props: any) => <div data-testid="PageBottomNav" {...props} />,
}))
// Pagination は props を受け取りやすく、操作できるように簡易 UI を自前で
vi.mock('../../common/Pagination', () => ({
  default: (props: any) => {
    const {
      page,
      totalPages,
      pageSize,
      onPageChange,
      pageSizeOptions,
      onPageSizeChange,
    } = props
    return (
      <div
        data-testid="Pagination"
        data-page={page}
        data-total={totalPages}
        data-size={pageSize}
      >
        <button onClick={() => onPageChange(1)}>to-1</button>
        <button onClick={() => onPageChange(2)}>to-2</button>
        {/* 2つ目のオプション（例: 20）に切り替えるスイッチ */}
        <button onClick={() => onPageSizeChange(pageSizeOptions[1])}>
          size-{pageSizeOptions[1]}
        </button>
      </div>
    )
  },
}))
// Card は見た目だけ
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
}))

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>)

/* ------------ テストデータ作成ヘルパ ------------ */
type ResultSummary = import('@/types/result').ResultSummary

const makeResult = (n: number, iso: string, extra?: Partial<ResultSummary>) =>
  ({
    quizNumber: n,
    createdAt: iso, // 並び替えに使う
    // 0=全/1=登録のみ/2=未登録のみ
    isRegisteredWords: 0,
    // 0=全て/1=含む/2=含まない
    isIdioms: 0,
    isSpecialCharacters: 0,
    choicesPosIds: [1, 3, 99], // 99 は未知 -> そのまま "99"
    totalQuestionsCount: 10,
    correctCount: 8,
    resultCorrectRate: 0.8 * 100, // 80
    ...extra,
  }) as ResultSummary

const makeMany = (count: number): ResultSummary[] => {
  // createdAt は古→新へ増加、コンポーネント側で「新しい順」にソートされる
  // quizNumber も 1..count としておく
  return Array.from({ length: count }, (_, i) => {
    const idx = i + 1
    const d = new Date(2025, 0, idx, 12, 0, 0) // 2025-01-(idx) 12:00
    return makeResult(idx, d.toISOString())
  })
}

beforeEach(() => {
  vi.resetAllMocks()
  navigateMock.mockReset()
})

describe('ResultIndex', () => {
  it('初期: 読み込み中 → 解決後に一覧', async () => {
    let resolveGet: (v: any) => void
    ;(axiosInstance.get as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolveGet = res
        }),
    )

    renderWithRouter(<ResultIndex />)
    // ローディング表示
    expect(screen.getByText('読み込み中…')).toBeInTheDocument()

    // 解決: 空配列
    resolveGet!({ data: [] })
    expect(await screen.findByText('全 0 件')).toBeInTheDocument()
    expect(screen.getByText('成績がありません')).toBeInTheDocument()
  })

  it('取得失敗: エラーメッセージ', async () => {
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('NG'))
    renderWithRouter(<ResultIndex />)
    expect(
      await screen.findByText('成績の取得に失敗しました'),
    ).toBeInTheDocument()
  })

  it('空リスト: 件数表示とプレースホルダ行', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: [] })
    renderWithRouter(<ResultIndex />)
    expect(await screen.findByText('全 0 件')).toBeInTheDocument()
    expect(screen.getByText('成績がありません')).toBeInTheDocument()
    // Pagination は最低1ページ
    const pg = screen.getByTestId('Pagination')
    expect(pg).toHaveAttribute('data-page', '1')
    expect(pg).toHaveAttribute('data-total', '1')
    expect(pg).toHaveAttribute('data-size', '10')
  })
  it('通常リスト: 新しい順に並ぶ & 列表示（マッピング/数値/品詞名/未知ID）', async () => {
    // このテスト専用のレスポンスをセット
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: [
        // 202（最新）: 登録のみ / 慣用句=含まない / 特殊=全て / 品詞=代名, 動詞
        makeResult(202, '2025-02-01T09:00:00.000Z', {
          isRegisteredWords: 1,
          isIdioms: 2,
          isSpecialCharacters: 0,
          choicesPosIds: [2, 3],
          totalQuestionsCount: 20,
          correctCount: 18,
          resultCorrectRate: 90.0,
        }),
        // 101: 全 / 含む / 含まない / 名詞, 形容, 99
        makeResult(101, '2025-01-01T09:00:00.000Z', {
          isRegisteredWords: 0,
          isIdioms: 1,
          isSpecialCharacters: 2,
          choicesPosIds: [1, 4, 99],
          totalQuestionsCount: 15,
          correctCount: 9,
          resultCorrectRate: 60.0,
        }),
        // 303（最も古い）: 未登録のみ / 全て / 含む / 副詞
        makeResult(303, '2024-12-31T09:00:00.000Z', {
          isRegisteredWords: 2,
          isIdioms: 0,
          isSpecialCharacters: 1,
          choicesPosIds: [5],
          totalQuestionsCount: 10,
          correctCount: 5,
          resultCorrectRate: 50.0,
        }),
      ],
    })

    renderWithRouter(<ResultIndex />)

    // 件数表示
    expect(await screen.findByText('全 3 件')).toBeInTheDocument()

    // ▼ 並び順（新しい順）を厳密に確認したい場合はこうすると堅い
    const dataRows = screen.getAllByRole('row').slice(1) // thead除外
    const firstCells = dataRows.map((r) =>
      r.querySelector('td')!.textContent?.trim(),
    )
    expect(firstCells).toEqual(['202', '101', '303'])

    // ▼ 各行をスコープして検証（重複テキスト対策）
    {
      const row = screen.getByText('202').closest('tr')!
      const r = within(row)
      r.getByText('登録のみ')
      r.getByText('含まない')
      r.getByText('全て')
      r.getByText('代名, 動詞')
      r.getByText('20')
      r.getByText('18')
      r.getByText('90.0%')
    }
    {
      const row = screen.getByText('101').closest('tr')!
      const r = within(row)
      r.getByText('全')
      r.getByText('含む')
      r.getByText('含まない')
      r.getByText('名詞, 形容, 99')
      r.getByText('15')
      r.getByText('9')
      r.getByText('60.0%')
    }
    {
      const row = screen.getByText('303').closest('tr')!
      const r = within(row)
      r.getByText('未登録のみ')
      r.getByText('全て')
      r.getByText('含む')
      r.getByText('副詞')
      r.getByText('10')
      r.getByText('5')
      r.getByText('50.0%')
    }
  })

  it('行クリックで /results/:quizNumber へ遷移', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: [makeResult(777, '2025-03-01T00:00:00.000Z')],
    })
    renderWithRouter(<ResultIndex />)

    // TD のテキストクリックで TR の onClick が反応する想定（イベントバブリング）
    await screen.findByText('全 1 件')
    await userEvent.click(screen.getByText('777'))

    expect(navigateMock).toHaveBeenCalledWith('/results/777')
  })

  it('Enter キーでも遷移（tabIndex=0 の行をフォーカスして Enter）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: [makeResult(888, '2025-04-01T00:00:00.000Z')],
    })
    renderWithRouter(<ResultIndex />)

    await screen.findByText('全 1 件')
    // tbody の行（ヘッダ行を除外）
    const allRows = screen.getAllByRole('row')
    const dataRow = allRows.find((r) => r.textContent?.includes('888'))!
    dataRow.focus()
    fireEvent.keyDown(dataRow, { key: 'Enter', code: 'Enter' })
    expect(navigateMock).toHaveBeenCalledWith('/results/888')
  })

  it('ページング: 初期はページ1/サイズ10 → 2ページ目へ → サイズ20で先頭へ戻る', async () => {
    // 25件（新しい順で 25..1 が表示される想定）
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: makeMany(25),
    })
    renderWithRouter(<ResultIndex />)

    // 初期: 件数と Pagination の属性
    expect(await screen.findByText('全 25 件')).toBeInTheDocument()
    const pg = screen.getByTestId('Pagination')
    expect(pg).toHaveAttribute('data-page', '1')
    expect(pg).toHaveAttribute('data-total', '3') // 25/10 = 3ページ
    expect(pg).toHaveAttribute('data-size', '10')

    // 表示行数は 10 件（ヘッダ除く）
    const rowsPage1 = screen.getAllByRole('row').slice(1) // 先頭は thead
    expect(rowsPage1).toHaveLength(10)

    // 2ページへ（モック Pagination のボタン）
    await userEvent.click(screen.getByText('to-2'))
    // 再描画後の Pagination 属性
    const pg2 = screen.getByTestId('Pagination')
    expect(pg2).toHaveAttribute('data-page', '2')

    // 2ページ目も 10 件
    const rowsPage2 = screen.getAllByRole('row').slice(1)
    expect(rowsPage2).toHaveLength(10)

    // ページサイズ 20 へ変更 → 先頭へ戻る（page が 1 に戻る）
    await userEvent.click(screen.getByText(/size-20/))
    const pg3 = screen.getByTestId('Pagination')
    expect(pg3).toHaveAttribute('data-size', '20')
    expect(pg3).toHaveAttribute('data-page', '1')

    // 先頭ページの表示件数は 20
    const rowsSize20 = screen.getAllByRole('row').slice(1)
    expect(rowsSize20).toHaveLength(20)
  })
})
