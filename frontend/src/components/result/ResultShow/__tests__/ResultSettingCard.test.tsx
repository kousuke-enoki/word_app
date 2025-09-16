// src/components/result/__tests__/ResultSettingCard.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, within } from '@testing-library/react'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ResultSettingCard from '../ResultSettingCard'

// Card を薄くモック（見た目依存を排除）
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
}))

// QuizStatus を固定文言でモック（実装差分に影響されないように）
vi.mock('@/lib/QuizStatus', () => ({
  QuizStatus: {
    registered: ['全', '登録のみ'], // index: 0/1
    idioms: ['含む', '含まない'], // index: 0/1
    special: ['含む', '含まない'], // index: 0/1
  },
}))

type Setting = {
  isRegisteredWords: 0 | 1
  isIdioms: 0 | 1
  isSpecialCharacters: 0 | 1
}

const renderWith = (setting: Setting) =>
  render(<ResultSettingCard setting={setting as any} />)

beforeEach(() => {
  vi.clearAllMocks()
})

describe('ResultSettingCard', () => {
  it('見出しと定義リストの見出しが表示される', () => {
    renderWith({
      isRegisteredWords: 0,
      isIdioms: 0,
      isSpecialCharacters: 0,
    })

    // 見出し
    expect(
      screen.getByRole('heading', { name: '出題条件', level: 2 }),
    ).toBeDefined()

    // 定義項目のラベル
    expect(screen.getByText('登録単語')).toBeInTheDocument()
    expect(screen.getByText('慣用句')).toBeInTheDocument()
    expect(screen.getByText('特殊単語')).toBeInTheDocument()

    // Card ラッパーがある（UIラッパの有無だけ確認）
    expect(screen.getByTestId('Card')).toBeInTheDocument()
  })

  // 2 x 2 x 2 の全パターンを網羅
  it.each([
    // isRegisteredWords, isIdioms, isSpecialCharacters, expected texts
    [
      { r: 0, i: 0, s: 0 },
      { reg: '全', idi: '含む', sp: '含む' },
    ],
    [
      { r: 0, i: 0, s: 1 },
      { reg: '全', idi: '含む', sp: '含まない' },
    ],
    [
      { r: 0, i: 1, s: 0 },
      { reg: '全', idi: '含まない', sp: '含む' },
    ],
    [
      { r: 0, i: 1, s: 1 },
      { reg: '全', idi: '含まない', sp: '含まない' },
    ],
    [
      { r: 1, i: 0, s: 0 },
      { reg: '登録のみ', idi: '含む', sp: '含む' },
    ],
    [
      { r: 1, i: 0, s: 1 },
      { reg: '登録のみ', idi: '含む', sp: '含まない' },
    ],
    [
      { r: 1, i: 1, s: 0 },
      { reg: '登録のみ', idi: '含まない', sp: '含む' },
    ],
    [
      { r: 1, i: 1, s: 1 },
      { reg: '登録のみ', idi: '含まない', sp: '含まない' },
    ],
  ])(
    '全パターン表示 (registered=%j)',
    (
      idx: { r: number; i: number; s: number },
      expected: { reg: string; idi: string; sp: string },
    ) => {
      renderWith({
        isRegisteredWords: idx.r as 0 | 1,
        isIdioms: idx.i as 0 | 1,
        isSpecialCharacters: idx.s as 0 | 1,
      })

      // 「登録単語」行の dd
      {
        const row = screen.getByText('登録単語').closest('div')!
        const dd = within(row).getByText(expected.reg)
        expect(dd).toBeInTheDocument()
      }

      // 「慣用句」行の dd
      {
        const row = screen.getByText('慣用句').closest('div')!
        const dd = within(row).getByText(expected.idi)
        expect(dd).toBeInTheDocument()
      }

      // 「特殊単語」行の dd
      {
        const row = screen.getByText('特殊単語').closest('div')!
        const dd = within(row).getByText(expected.sp)
        expect(dd).toBeInTheDocument()
      }
    },
  )

  it('要素の構造: <dl> 内に3つの項目がある', () => {
    renderWith({ isRegisteredWords: 0, isIdioms: 1, isSpecialCharacters: 1 })
    const card = screen.getByTestId('Card')
    const dl = card.querySelector('dl') as HTMLDListElement
    expect(dl).toBeTruthy()
    const terms = within(dl).getAllByRole('term') // <dt>
    const defs = within(dl).getAllByRole('definition') // <dd>
    expect(terms).toHaveLength(3)
    expect(defs).toHaveLength(3)
    // それぞれの並びを軽く確認（任意）
    expect(terms.map((n) => n.textContent)).toEqual([
      '登録単語',
      '慣用句',
      '特殊単語',
    ])
  })
})
