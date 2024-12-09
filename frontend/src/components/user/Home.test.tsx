import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
// import { rest } from 'msw'
import { setupServer } from 'msw/node'
import Home from './Home'
import { handlers } from '../../mocks/handlers'

const server = setupServer(...handlers)

// モックサーバーのセットアップ
beforeAll(() => {
  server.listen({
    onUnhandledRequest: 'error', // 未処理のリクエストをエラーとして扱う
  })
})
beforeEach(() => {
  // テスト開始前に localStorage をクリア
  localStorage.clear()
})
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('Home Component', () => {
  test('ホームコンポーネントのテスト', async () => {
    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )
    expect(screen.getByText('トップページです。')).toBeInTheDocument()
  })
  test('ログイン済みのユーザーにはMyPageを表示する', async () => {
    // ローカルストレージに有効なトークンをセット
    localStorage.setItem('token', 'valid-token')

    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )

    // "Test User" というユーザー名が表示されることを確認
    await waitFor(() => {
      expect(
        screen.getByText((content) => content.includes('Test User')),
      ).toBeInTheDocument()
    })

    // "ログアウト" ボタンが表示されていることを確認
    // eslint-disable-next-line testing-library/prefer-presence-queries
    expect(screen.getByText('マイページ')).toBeInTheDocument()
    expect(screen.getByText('全単語リスト:')).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeInTheDocument()
  })

  test('未ログイン状態の場合はログインを促すメッセージを表示', async () => {
    // ローカルストレージにトークンをセットしない（未ログイン状態）
    localStorage.removeItem('token')

    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )

    // "トップページです。" が表示されていることを確認
    expect(screen.getByText('トップページです。')).toBeInTheDocument()
    // サインアップとサインインのリンクが表示されていることを確認
    expect(screen.getByText('サインアップページ')).toBeInTheDocument()
    expect(screen.getByText('サインインページ')).toBeInTheDocument()
  })

  test('認証エラー時にトークンが削除され、ログインを促すメッセージを表示', async () => {
    // ローカルストレージに無効なトークンをセット
    localStorage.clear()
    localStorage.setItem('token', 'invalid-token')

    render(
      <MemoryRouter>
        <Home />
      </MemoryRouter>,
    )

    // "ログインしてください" のメッセージが表示されることを確認
    await waitFor(() => {
      expect(screen.getByText('ログインしてください')).toBeInTheDocument()
    })

    // ローカルストレージからトークンが削除されていることを確認
    expect(localStorage.getItem('token')).toBeNull()
  })
})
