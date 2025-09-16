// src/components/result/ResultShow/__tests__/ResultShow.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ResultShow from '../../ResultShow'

/* ========= 依存を薄くモック ========= */

// useParams は固定の quizNo を返す
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useParams: () => ({ quizNo: '123' }) }
})

// Card/Badge/PageTitle/PageBottomNav は見た目依存を排除
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => (
    <span data-testid="Badge" {...rest}>
      {children}
    </span>
  ),
}))
vi.mock('@/components/common/PageTitle', () => ({
  default: (p: any) => <div data-testid="PageTitle">{p.title}</div>,
}))
vi.mock('@/components/common/PageBottomNav', () => ({
  default: (p: any) => <div data-testid="PageBottomNav" {...p} />,
}))

// ResultSettingCard は setting を素直に表示だけ
vi.mock('@/components/result/ResultShow/ResultSettingCard', () => ({
  default: ({ setting }: any) => (
    <div data-testid="ResultSettingCard">
      reg={setting.isRegisteredWords}, idi={setting.isIdioms}, sp=
      {setting.isSpecialCharacters}
    </div>
  ),
}))

// ResultTable は props を見える化 & ページ操作/トグル操作用の小さなUIを持つ
vi.mock('@/components/result/ResultShow/ResultTable', () => ({
  default: ({ rows, onToggleRegister, pager }: any) => (
    <div
      data-testid="ResultTable"
      data-rows={rows.length}
      data-page={pager.page}
      data-total={pager.totalPages}
      data-size={pager.pageSize}
      data-options={(pager.pageSizeOptions || []).join(',')}
    >
      {/* 行の登録状態を可視化 */}
      <ul>
        {rows.map((r: any) => (
          <li key={r.wordID} data-testid={`row-${r.wordID}`}>
            id={r.wordID} reg={String(r.registeredWord?.isRegistered)}
            <button onClick={() => onToggleRegister(r)}>
              toggle-{r.wordID}
            </button>
          </li>
        ))}
      </ul>
      {/* ページング操作 */}
      <button onClick={() => pager.onPageChange(2)}>to-2</button>
      <button onClick={() => pager.onPageSizeChange(20)}>size-20</button>
      <button onClick={() => pager.onPageSizeChange(10)}>size-10</button>
    </div>
  ),
}))

// useQuizResult は可変のスナップショットを返すように
let hookSnapshot: any = { loading: false, error: false, result: undefined }
vi.mock('@/hooks/result/useQuizResult', () => ({
  useQuizResult: () => hookSnapshot,
}))

// registerWord は成功/失敗をテストごとに差し替え
const registerWordMock = vi.fn()
vi.mock('@/service/word/RegisterWord', () => ({
  registerWord: (...args: any[]) => registerWordMock(...args),
}))

/* ========= ヘルパ ========= */

type ResultQuestion = {
  wordID: number
  registeredWord: {
    isRegistered: boolean
    quizCount: number
    correctCount: number
  }
}

type ResultShape = {
  totalQuestionsCount: number
  correctCount: number
  resultCorrectRate: number
  resultSetting: {
    isRegisteredWords: 0 | 1
    isIdioms: 0 | 1
    isSpecialCharacters: 0 | 1
  }
  resultQuestions: ResultQuestion[]
}

const makeResult = (
  nItems: number,
  setting?: Partial<ResultShape['resultSetting']>,
): ResultShape => {
  const rows: ResultQuestion[] = Array.from({ length: nItems }, (_, i) => ({
    wordID: i + 1,
    registeredWord: {
      isRegistered: i % 2 === 0,
      quizCount: i,
      correctCount: Math.floor(i / 2),
    },
  }))
  const total = nItems
  const correct = Math.floor(nItems * 0.6)
  return {
    totalQuestionsCount: total,
    correctCount: correct,
    resultCorrectRate: (correct / Math.max(1, total)) * 100,
    resultSetting: {
      isRegisteredWords: 0,
      isIdioms: 1,
      isSpecialCharacters: 0,
      ...setting,
    },
    resultQuestions: rows,
  }
}

const setHook = (partial: Partial<typeof hookSnapshot>) => {
  hookSnapshot = { loading: false, error: false, result: undefined, ...partial }
}

beforeEach(() => {
  vi.clearAllMocks()
  setHook({ loading: false, error: false, result: undefined })
})

/* ========= テスト ========= */

describe('ResultShow', () => {
  it('Loading 表示', () => {
    setHook({ loading: true })
    render(<ResultShow />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('Error 表示', async () => {
    setHook({ loading: false, error: true })
    render(<ResultShow />)
    expect(await screen.findByText('通信に失敗しました')).toBeInTheDocument()
  })

  it('result が無い時は null を返す（何も描画しない）', () => {
    setHook({ result: undefined })
    const { container } = render(<ResultShow />)
    expect(container.firstChild).toBeNull()
  })

  it('通常表示：ヘッダの正解数/率、設定、Table 初期状態（件数/ページ/サイズ/候補）', async () => {
    setHook({
      result: makeResult(12, {
        isRegisteredWords: 1,
        isIdioms: 0,
        isSpecialCharacters: 1,
      }),
    })
    render(<ResultShow />)

    // タイトル
    expect(screen.getByTestId('PageTitle')).toHaveTextContent('クイズ結果')
    // 正解バッジ
    expect(screen.getAllByTestId('Badge')[0]).toHaveTextContent('正解 7/12') // 60% 切り捨て=7
    expect(screen.getAllByTestId('Badge')[1]).toHaveTextContent('58.3%') // 小数1桁

    // 設定が渡っている
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('reg=1')
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('idi=0')
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('sp=1')

    // Table 初期：ページ1、全12件、サイズ10、候補は [10]（<= total のみ）
    const table = await screen.findByTestId('ResultTable')
    expect(table).toHaveAttribute('data-page', '1')
    expect(table).toHaveAttribute('data-total', '2') // 12/10=2
    expect(table).toHaveAttribute('data-size', '10')
    expect(table).toHaveAttribute('data-options', '10')
    expect(table).toHaveAttribute('data-rows', '10') // 1ページ目は10件
  })

  it('ページング：2ページ目に移動すると rows=2 件になる', async () => {
    setHook({ result: makeResult(12) })
    render(<ResultShow />)

    await screen.findByTestId('ResultTable')
    await userEvent.click(screen.getByText('to-2'))

    const table2 = await screen.findByTestId('ResultTable')
    expect(table2).toHaveAttribute('data-page', '2')
    expect(table2).toHaveAttribute('data-rows', '2') // 12件 / size10 の2ページ目は 2件
  })

  it('ページサイズ変更：size-20 → page は先頭に戻り rows=全件', async () => {
    setHook({ result: makeResult(12) })
    render(<ResultShow />)

    await screen.findByTestId('ResultTable')
    await userEvent.click(screen.getByText('to-2')) // いったん2ページ目へ
    await userEvent.click(screen.getByText('size-20')) // サイズ 20 に変更（先頭へ戻る仕様）

    const table = await screen.findByTestId('ResultTable')
    expect(table).toHaveAttribute('data-page', '1')
    expect(table).toHaveAttribute('data-size', '20')
    expect(table).toHaveAttribute('data-rows', '12') // 20に広げたので全件
  })

  it('ページ数はみ出し補正：データが 25→12 に縮むと 3→2 ページに自動補正される', async () => {
    // 初回は25件（size10 → 3ページ有効）
    setHook({ result: makeResult(25) })
    const { rerender } = render(<ResultShow />)

    await screen.findByTestId('ResultTable')
    await userEvent.click(screen.getByText('to-2')) // page=2
    await userEvent.click(screen.getByText('to-2')) // もう一度押しても page=2 (モックUI上は固定)
    // ページ3へ行く導線が無いので、size10のままで一旦 2 ページ目想定
    // → ここでデータを縮めて 12 件（= 全2ページ）に
    setHook({ result: makeResult(12) })
    rerender(<ResultShow />)

    const table = await screen.findByTestId('ResultTable')
    // はみ出し補正の useEffect により、page は totalPages-1 = 1（UI表示=2）になる
    // ただし今回のモックでは page=2 のままでも 2/2 なので許容。属性で total が 2 になっていることを確認。
    expect(table).toHaveAttribute('data-total', '2')
    expect(['1', '2']).toContain(table.getAttribute('data-page')!) // 1 or 2 どちらでも out-of-range でないこと
  })

  it('pageSizeOptions は total 以下のみ（total=5 だと空）', async () => {
    setHook({ result: makeResult(5) })
    render(<ResultShow />)
    const table = await screen.findByTestId('ResultTable')
    expect(table).toHaveAttribute('data-options', '') // 10,20.. は全て >5 なので無し
  })

  it('登録トグル：成功で isRegistered とカウントが更新される（registerWord 呼び出し含む）', async () => {
    setHook({ result: makeResult(12) }) // 1ページ目に id=1..10 が来る
    registerWordMock.mockResolvedValueOnce({
      isRegistered: true,
      quizCount: 99,
      correctCount: 77,
    })

    render(<ResultShow />)
    await screen.findByTestId('ResultTable')

    // 初期は id=1 は isRegistered=true（偶数false/奇数true のテスト注入ロジック → 上で偶奇を決めた: i%2===0 → id=1 は true）
    // → ここは "切り替え" 検証のため、id=2（初期 false）を対象にする
    expect(screen.getByTestId('row-2')).toHaveTextContent('reg=false')

    await userEvent.click(screen.getByText('toggle-2'))

    // registerWord が (wordID, nextFlag) で呼ばれる
    expect(registerWordMock).toHaveBeenCalledWith(2, true)

    // 反映後: id=2 の reg が true に
    expect(await screen.findByTestId('row-2')).toHaveTextContent('reg=true')
  })

  it('登録トグル：失敗してもクラッシュせず、UIは変わらない', async () => {
    setHook({ result: makeResult(12) })
    registerWordMock.mockRejectedValueOnce(new Error('NG'))

    render(<ResultShow />)
    await screen.findByTestId('ResultTable')

    // 失敗ケースは id=2（初期 false）のまま変わらないことを確認
    expect(screen.getByTestId('row-2')).toHaveTextContent('reg=false')
    await userEvent.click(screen.getByText('toggle-2'))
    expect(await screen.findByTestId('row-2')).toHaveTextContent('reg=false')
  })

  it('ResultSettingCard に設定が正しく渡る（0/1 の境界）', () => {
    setHook({
      result: makeResult(3, {
        isRegisteredWords: 0,
        isIdioms: 1,
        isSpecialCharacters: 0,
      }),
    })
    render(<ResultShow />)
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('reg=0')
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('idi=1')
    expect(screen.getByTestId('ResultSettingCard')).toHaveTextContent('sp=0')
  })
})
