// src/components/user/Home.test.tsx
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import Home from '../Home'

/* ThemeContext をモック */
const setThemeMock = vi.fn()
vi.mock('@/contexts/themeContext', () => ({
  useTheme: () => ({ setTheme: setThemeMock }),
}))

beforeEach(() => {
  localStorage.clear()
  setThemeMock.mockClear()
})

describe('Home', () => {
  it('通常表示', () => {
    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )
    expect(screen.getByText('トップページです。'))
    expect(
      screen.getByRole('link', { name: /サインインはここから！/ }))
    expect(setThemeMock).toHaveBeenCalledWith('light')
  })

  it('logoutMessage を表示し localStorage から削除', () => {
    localStorage.setItem('logoutMessage', 'サインアウトしました')

    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )

    expect(screen.getByText('サインアウトしました'))
    expect(localStorage.getItem('logoutMessage')).toBeNull()
  })
})
