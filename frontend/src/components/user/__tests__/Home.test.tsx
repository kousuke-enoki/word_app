import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import Home from '../Home'

/* ThemeContext をモック */
const setThemeMock = vi.fn()
vi.mock('@/contexts/themeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

/* 外部 UI を薄くモック（レンダリング安定化） */
vi.mock('@/components/ui/PageShell', () => ({
  PageShell: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="PageShell">{children}</div>
  ),
}))
vi.mock('@/components/ui/card', () => ({
  Card: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="Card">{children}</div>
  ),
  PageContainer: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="PageContainer">{children}</div>
  ),
}))
vi.mock('@/components/ui/ui', () => ({
  // Link 内に入れても role=link の名前計算に使われるよう、span でOK
  Button: ({
    children,
    ...rest
  }: React.PropsWithChildren<React.HTMLAttributes<HTMLSpanElement>>) => (
    <span {...rest}>{children}</span>
  ),
}))

beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
})

const renderHome = () =>
  render(
    <MemoryRouter>
      <Home />
    </MemoryRouter>,
  )

describe('Home', () => {
  it('通常表示：見出しと主要リンクがあり、テーマをlightに設定', async () => {
    renderHome()

    // 見出し
    expect(
      screen.getByRole('heading', { name: '英単語を、もっと覚えやすく。' }),
    ).toBeInTheDocument()

    // 主要導線（role=link でチェック）
    expect(screen.getByRole('link', { name: 'サインイン' })).toHaveAttribute(
      'href',
      '/sign_in',
    )

    // logoutMessage はない
    expect(screen.queryByText(/サインアウトしました/)).not.toBeInTheDocument()

    // useEffect の呼び出しを待つ
    await waitFor(() => expect(setThemeMock).toHaveBeenCalledWith('light'))
  })

  it('logoutMessage を表示し、useEffect 後に localStorage から削除される', async () => {
    localStorage.setItem('logoutMessage', 'サインアウトしました')
    renderHome()

    // 初回描画では表示されている
    expect(screen.getByText('サインアウトしました')).toBeInTheDocument()

    // 削除は useEffect 後に行われるので待つ
    await waitFor(() =>
      expect(localStorage.getItem('logoutMessage')).toBeNull(),
    )
  })

  it('機能カードが3枚あり、各リンク先が正しい', () => {
    renderHome()
    const links = [
      { name: /単語リスト/, href: '/words' },
      { name: /登録/, href: '/words/new' },
      { name: /クイズ/, href: '/quizs' },
    ]
    for (const { name, href } of links) {
      const link = screen.getByRole('link', { name })
      expect(link).toHaveAttribute('href', href)
    }
  })

  it('logoutMessage は一度表示されたら、再マウント時には表示されない（ワンタイム表示）', async () => {
    localStorage.setItem('logoutMessage', 'サインアウトしました')
    const { unmount } = renderHome()

    expect(screen.getByText('サインアウトしました')).toBeInTheDocument()
    await waitFor(() =>
      expect(localStorage.getItem('logoutMessage')).toBeNull(),
    )

    // 再マウント（次回訪問想定）
    unmount()
    renderHome()
    expect(screen.queryByText('サインアウトしました')).not.toBeInTheDocument()
  })
})
