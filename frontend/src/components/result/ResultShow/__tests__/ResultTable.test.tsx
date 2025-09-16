// src/components/result/ResultShow/__tests__/ResultTable.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ResultTable from '../../ResultShow/ResultTable'

/* ========== 依存の薄いモック ========== */
vi.mock('@/components/common/RegisterToggle', () => ({
  // import { RegisterToggle } なので名前付きエクスポート
  RegisterToggle: ({ isRegistered, onToggle, widthClass }: any) => (
    <button
      data-testid="RegisterToggle"
      data-reg={String(isRegistered)}
      data-width={widthClass}
      onClick={onToggle}
    >
      toggle
    </button>
  ),
}))

vi.mock('@/components/common/Pagination', () => ({
  default: (props: any) => {
    const {
      page,
      totalPages,
      pageSize,
      onPageChange,
      onPageSizeChange,
      className,
      compact,
    } = props
    return (
      <div
        data-testid="Pagination"
        data-page={page}
        data-total={totalPages}
        data-size={pageSize}
        data-compact={String(!!compact)}
        data-class={className || ''}
      >
        <button onClick={() => onPageChange(2)}>go-2</button>
        <button onClick={() => onPageSizeChange(30)}>size-30</button>
      </div>
    )
  },
}))

vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
}))

/* ========== ヘルパ / ダミーデータ ========== */
type ResultQuestion = import('@/types/quiz').ResultQuestion

const row = (over: Partial<ResultQuestion> = {}): ResultQuestion => ({
  quizID: over.quizID ?? 1,
  questionNumber: over.questionNumber ?? 1,
  wordID: over.wordID ?? 101,
  wordName: over.wordName ?? 'apple',
  posID: over.posID ?? 1,
  correctJpmID: over.correctJpmID ?? 1,
  answerJpmID: over.answerJpmID ?? 1,
  isCorrect: over.isCorrect ?? true,
  choicesJpms: over.choicesJpms ?? [
    { japaneseMeanID: 1, name: '正解の意味' },
    { japaneseMeanID: 2, name: '別の意味' },
  ],
  registeredWord:
    over.registeredWord ??
    ({ isRegistered: false, quizCount: 3, correctCount: 2 } as any),
  timeMs: over.timeMs ?? 1000,
})

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter initialEntries={['/']}>{ui}</MemoryRouter>)

beforeEach(() => {
  vi.clearAllMocks()
})

// tbody を取得するヘルパー関数
const getTBody = () =>
  screen
    .getAllByRole('rowgroup')
    .find(
      (el) => (el as HTMLElement).tagName.toLowerCase() === 'tbody',
    ) as HTMLElement

/* ========== テスト本体 ========== */
describe('ResultTable', () => {
  it('基本：行が描画され、リンク/タイトル/数値列が正しい', () => {
    const rows = [
      row({
        questionNumber: 1,
        wordID: 10,
        wordName: 'alpha',
        isCorrect: true,
      }),
      row({
        questionNumber: 2,
        wordID: 20,
        wordName: 'beta',
        isCorrect: false,
        answerJpmID: 2,
      }),
    ]
    const onToggle = vi.fn()

    render(
      <MemoryRouter>
        <ResultTable rows={rows} onToggleRegister={onToggle} />
      </MemoryRouter>,
    )

    // tbody だけを見る
    const body = getTBody()
    const trs = within(body).getAllByRole('row')
    expect(trs).toHaveLength(2)

    // 1行目
    {
      const r = trs[0]
      within(r).getByText('1')
      const link = within(r).getByRole('link', { name: 'alpha' })
      expect(link).toHaveAttribute('href', '/words/10')
      expect(link).toHaveAttribute('title', 'alpha')

      // title が重複しうるので All を使って列を特定
      const [correctTd, selectTd] = within(r).getAllByTitle('正解の意味')
      expect(correctTd).toBeInTheDocument()
      expect(selectTd).toBeInTheDocument()

      within(r).getByText('3')
      within(r).getByText('2')

      const toggle = within(r).getByTestId('RegisterToggle')
      expect(toggle).toHaveAttribute('data-reg', 'false')
      expect(toggle).toHaveAttribute('data-width', 'w-24 sm:w-28')
    }

    // 2行目
    {
      const r = trs[1]
      const tds = within(r).getAllByRole('cell')
      // #列
      expect(tds[0]).toHaveTextContent('2')
      // 単語リンク（2列目のセル内）
      within(tds[1]).getByRole('link', { name: 'beta' })
      // 正解（3列目）: td 自身の title を検証
      expect(tds[2]).toHaveAttribute('title', '正解の意味')
      // 選択（4列目）
      expect(tds[3]).toHaveAttribute('title', '別の意味')
    }
  })

  it('正誤ハイライト：isCorrect=true は emerald 系、false は rose 系クラス', () => {
    const rows = [
      row({ questionNumber: 1, isCorrect: true, answerJpmID: 1 }),
      row({ questionNumber: 2, isCorrect: false, answerJpmID: 2 }),
    ]
    render(
      <MemoryRouter>
        <ResultTable rows={rows} onToggleRegister={() => {}} />
      </MemoryRouter>,
    )

    const body = getTBody()
    const trs = within(body).getAllByRole('row')

    // 1行目：選択セル（titleが重複するため2番目=選択セルを使う）
    {
      const [, selectTd] = within(trs[0]).getAllByTitle('正解の意味')
      const cls = selectTd.getAttribute('class') || ''
      expect(cls).toMatch(/emerald/i)
      expect(cls).not.toMatch(/rose/i)
    }

    // 2行目：選択セルは '別の意味'
    {
      const selectTd = within(trs[1]).getByTitle('別の意味')
      const cls = selectTd.getAttribute('class') || ''
      expect(cls).toMatch(/rose/i)
      expect(cls).not.toMatch(/emerald/i)
    }
  })

  it("正解/選択のIDが choicesJpms に存在しない場合、'-' を表示", () => {
    const rows = [
      row({
        correctJpmID: 999,
        answerJpmID: 888,
        choicesJpms: [{ japaneseMeanID: 1, name: 'X' }], // 999/888 は存在しない
      }),
    ]
    renderWithRouter(<ResultTable rows={rows} onToggleRegister={() => {}} />)

    // '-' が2箇所（正解/選択）に表示される
    const dashes = screen.getAllByText('-')
    expect(dashes.length).toBeGreaterThanOrEqual(2)
  })

  it('登録トグル：クリックで onToggleRegister(row) が呼ばれる', async () => {
    const rows = [
      row({ questionNumber: 1, wordID: 10 }),
      row({ questionNumber: 2, wordID: 20 }),
    ]
    const onToggle = vi.fn()
    render(
      <MemoryRouter>
        <ResultTable rows={rows} onToggleRegister={onToggle} />
      </MemoryRouter>,
    )

    const body = getTBody()
    const trs = within(body).getAllByRole('row')

    const toggle2 = within(trs[1]).getByTestId('RegisterToggle')
    await userEvent.click(toggle2)

    expect(onToggle).toHaveBeenCalledTimes(1)
    expect(onToggle.mock.calls[0][0]).toMatchObject({
      wordID: 20,
      questionNumber: 2,
    })
  })

  it('pager が無い時はフッター非表示', () => {
    renderWithRouter(<ResultTable rows={[row()]} onToggleRegister={() => {}} />)
    expect(screen.queryByTestId('Pagination')).toBeNull()
  })

  it('pager がある時：Pagination に props が渡り、イベントもフォワードされる', async () => {
    const onPageChange = vi.fn()
    const onPageSizeChange = vi.fn()
    const pager = {
      page: 1,
      totalPages: 5,
      pageSize: 10,
      onPageChange,
      onPageSizeChange,
      pageSizeOptions: [10, 30, 50],
      compact: true,
      className: 'custom-class',
    }

    renderWithRouter(
      <ResultTable
        rows={[row(), row({ questionNumber: 2 })]}
        onToggleRegister={() => {}}
        pager={pager}
      />,
    )

    const pg = await screen.findByTestId('Pagination')
    expect(pg).toHaveAttribute('data-page', '1')
    expect(pg).toHaveAttribute('data-total', '5')
    expect(pg).toHaveAttribute('data-size', '10')
    expect(pg).toHaveAttribute('data-compact', 'true')
    // コンポーネント側で `className="!mt-0 !mb-0"` を渡しているので、モックでは data-class に反映される
    expect(pg.getAttribute('data-class') || '').toContain('!mt-0')
    expect(pg.getAttribute('data-class') || '').toContain('!mb-0')

    // クリックでコールバックが呼ばれる
    await userEvent.click(screen.getByText('go-2'))
    expect(onPageChange).toHaveBeenCalledWith(2)
    await userEvent.click(screen.getByText('size-30'))
    expect(onPageSizeChange).toHaveBeenCalledWith(30)
  })

  it('レスポンシブ幅指定などの固定 props も流れている（RegisterToggle の widthClass）', () => {
    renderWithRouter(<ResultTable rows={[row()]} onToggleRegister={() => {}} />)
    const toggle = screen.getByTestId('RegisterToggle')
    expect(toggle).toHaveAttribute('data-width', 'w-24 sm:w-28')
  })
})
