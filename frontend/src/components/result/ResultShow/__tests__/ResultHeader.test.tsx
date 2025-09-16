// src/components/result/__tests__/ResultHeader.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen } from '@testing-library/react'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ResultHeader from '../ResultHeader'

/** UI層を薄くモック（表示テキストの検証に集中） */
vi.mock('@/components/ui/card', () => ({
  Badge: ({ children, ...rest }: any) => (
    <span data-testid="Badge" {...rest}>
      {children}
    </span>
  ),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

describe('ResultHeader', () => {
  it('見出しとバッジ2つが表示される（基本ケース）', () => {
    render(<ResultHeader correct={8} total={10} rate={80} />)

    // 見出し
    expect(
      screen.getByRole('heading', { name: 'クイズ結果' }),
    ).toBeInTheDocument()

    // バッジ2個
    const badges = screen.getAllByTestId('Badge')
    expect(badges).toHaveLength(2)

    // 正解/総数の表示
    expect(screen.getByText('正解 8/10')).toBeInTheDocument()

    // 率は小数1桁
    expect(screen.getByText('80.0%')).toBeInTheDocument()
  })

  it('率は toFixed(1) で四捨五入される', () => {
    render(<ResultHeader correct={10} total={15} rate={66.666} />)
    expect(screen.getByText('66.7%')).toBeInTheDocument()
  })

  it('0% の表示（正解0 / 総数>0）', () => {
    render(<ResultHeader correct={0} total={10} rate={0} />)
    expect(screen.getByText('正解 0/10')).toBeInTheDocument()
    expect(screen.getByText('0.0%')).toBeInTheDocument()
  })

  it('100% の表示（正解=総数）', () => {
    render(<ResultHeader correct={10} total={10} rate={100} />)
    expect(screen.getByText('正解 10/10')).toBeInTheDocument()
    expect(screen.getByText('100.0%')).toBeInTheDocument()
  })

  it.each([
    { correct: 1, total: 3, rate: 33.333, expected: '33.3%' },
    { correct: 7, total: 11, rate: 63.636, expected: '63.6%' },
    { correct: 12, total: 37, rate: 32.432, expected: '32.4%' },
  ])('いくつかの端数ケース %#', ({ correct, total, rate, expected }) => {
    render(<ResultHeader correct={correct} total={total} rate={rate} />)
    expect(screen.getByText(`正解 ${correct}/${total}`)).toBeInTheDocument()
    expect(screen.getByText(expected)).toBeInTheDocument()
  })
})
