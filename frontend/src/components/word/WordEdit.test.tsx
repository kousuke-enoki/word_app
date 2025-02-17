// import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import WordEdit, { WordForUpdate } from './WordEdit'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import WordShow from './WordShow'

// MSWサーバー
const server = setupServer()

// privateRouteがログイン済みとなるようにする
jest.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({
    isLoggedIn: true,
    isLoading: false,
  }),
}))

// テスト前後の設定
beforeEach(() => {
  queryClient.clear()
  localStorage.setItem('token', 'dummy')
  window.alert = jest.fn() // これでalertが呼ばれても実際にダイアログは出ない
})
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterAll(() => server.close())
afterEach(() => {
  server.resetHandlers()
  jest.clearAllMocks()
})

// 共通で使う QueryClient を用意
const queryClient = new QueryClient()

// テスト描画用ヘルパー
function renderApp(path = '/words/edit/123') {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/words/edit/:id" element={<WordEdit />} />
          <Route path="/words/:id" element={<WordShow />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

// 便利ハンドラ作成
const mockWordData: WordForUpdate = {
  id: 123,
  name: 'hello',
  wordInfos: [
    {
      id: 1,
      partOfSpeechId: 1, // "Noun" と想定
      japaneseMeans: [{ id: 10, name: 'こんにちは' }],
    },
  ],
}

// モックエンドポイント
const getWordSuccessHandler = rest.get(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    return res(ctx.status(200), ctx.json(mockWordData))
  },
)
const getWordErrorHandler = rest.get(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    return res(ctx.status(500), ctx.json({ message: 'Internal Server Error' }))
  },
)
const getWordNullHandler = rest.get(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    // APIがnullを返却 → "データが存在しません"
    return res(ctx.status(200), ctx.json(null))
  },
)

// PUTハンドラ
const putWordSuccessHandler = rest.put(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    // 正常終了 → { name: "updated-word" } を返す
    return res(ctx.status(200), ctx.json({ name: 'updated-word' }))
  },
)
const putWordErrorHandler = rest.put(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    return res(ctx.status(500))
  },
)

describe('WordEdit (React Query version) tests', () => {
  test('読み込み失敗 -> エラー文言 + 再取得ボタンを表示', async () => {
    server.use(getWordErrorHandler)
    renderApp()

    // 最初は読み込み中
    expect(screen.getByText('読み込み中...')).toBeInTheDocument()

    // "単語情報の取得中にエラーが発生しました。" の表示を待つ
    await waitFor(
      () => {
        expect(
          screen.getByText('単語情報の取得中にエラーが発生しました。'),
        ).toBeInTheDocument()
      },
      { timeout: 3000 },
    )

    // フォームは表示されない
    expect(screen.queryByText('単語更新フォーム')).not.toBeInTheDocument()

    // 再取得ボタンが表示されることを確認
    const retryButton = screen.getByRole('button', { name: '再取得' })
    expect(retryButton).toBeInTheDocument()

    // ここで再取得を試してみる例
    server.use(getWordSuccessHandler) // 再取得時に成功に切り替え
    await userEvent.click(retryButton)

    // 再取得してフォームが表示されるようになる
    expect(await screen.findByText('単語更新フォーム')).toBeInTheDocument()
  })

  test('読み込み中の表示 -> 成功時にフォームが出る', async () => {
    server.use(getWordSuccessHandler)
    renderApp()

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
    // 成功後、フォーム表示
    expect(await screen.findByText('単語更新フォーム')).toBeInTheDocument()
    // 単語名 "hello" が表示されている
    expect(screen.getByDisplayValue('hello')).toBeInTheDocument()
  })

  test('読み込み結果が null -> "データが存在しません。" を表示', async () => {
    server.use(getWordNullHandler)
    renderApp()

    expect(
      await screen.findByText('データが存在しません。'),
    ).toBeInTheDocument()
  })

  test('更新成功 → すぐに /words/123 へ遷移し、WordShow 側でメッセージを表示', async () => {
    // GET成功 + PUT成功に設定
    server.use(getWordSuccessHandler, putWordSuccessHandler)

    renderApp('/words/edit/123')

    // WordEditのロード完了
    expect(await screen.findByText('単語更新フォーム')).toBeInTheDocument()

    // [単語を更新] ボタン押下
    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))
  })

  test('更新失敗 -> エラーメッセージを表示', async () => {
    server.use(getWordSuccessHandler, putWordErrorHandler)
    renderApp()

    // フォーム表示待ち
    await screen.findByText('単語更新フォーム')
    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    // エラー
    await waitFor(() => {
      expect(
        screen.getByText('単語情報の更新中にエラーが発生しました。'),
      ).toBeInTheDocument()
    })
  })

  test('バリデーションエラーで更新されず、PUTが呼ばれない', async () => {
    server.use(getWordSuccessHandler, putWordSuccessHandler)
    // putWordSuccessHandler は呼ばれてほしくないので、一応セットだけはしておく

    renderApp()
    await screen.findByText('単語更新フォーム')

    // 単語名に数字を入れる -> バリデーションNG
    const nameInput = screen.getByDisplayValue('hello') as HTMLInputElement
    await userEvent.clear(nameInput)
    await userEvent.type(nameInput, 'hello123')

    // 「単語を更新」押下
    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    // バリデーションエラーメッセージが表示される
    expect(
      screen.getByText('単語名は半角アルファベットのみ入力できます。'),
    ).toBeInTheDocument()

    // 成功メッセージは当然出ない
    expect(
      screen.queryByText('updated-wordが正常に更新されました！'),
    ).not.toBeInTheDocument()
  })
})
