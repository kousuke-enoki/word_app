import { render, screen, waitFor } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import AppRouter from './AppRouter'
import { WordForUpdate } from '../components/word/WordEdit'
import userEvent from '@testing-library/user-event'

// (useAuth をモックして isLoggedIn=true に)
jest.mock('../hooks/useAuth', () => ({
  useAuth: () => ({
    isLoggedIn: true,
    isLoading: false,
  }),
}))

const server = setupServer()
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
beforeEach(() => localStorage.setItem('token', 'dummy'))
afterAll(() => server.close())
afterEach(() => {
  server.resetHandlers()
  jest.clearAllMocks()
})
const queryClient = new QueryClient()

function renderApp() {
  return render(
    <QueryClientProvider client={queryClient}>
      <AppRouter />
    </QueryClientProvider>,
  )
}

const mockWordData: WordForUpdate = {
  id: 0,
  name: '',
  wordInfos: [],
}

const getWordHandlerSuccess = rest.get(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    return res(ctx.status(200), ctx.json(mockWordData))
  },
)

const putWordHandlerSuccess = rest.put(
  'http://localhost:8080/words/:id',
  (req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ name: 'updated-word' }))
  },
)

describe('AppRouter Integration Tests', () => {
  test('WordEdit更新成功 → WordShow でメッセージ表示', async () => {
    server.use(getWordHandlerSuccess, putWordHandlerSuccess)

    renderApp()

    // テストで "/words/edit/123" に移動したい
    // BrowserRouter だとデフォルトで "/" から始まるため、
    // #1. いきなり Link/ボタンを押して移動する or
    // #2. history.push("/words/edit/123") する(歴史操作)

    // 例: テスト内で window.history
    window.history.pushState({}, 'TestTitle', '/words/edit/123')

    // -> AppRouter が再レンダーして WordEdit が表示される
    // 以降、同様に "単語更新フォーム" 確認 → ボタンクリック → /words/123
    expect(await screen.findByText('単語更新フォーム')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: '単語を更新' }))

    // navigate("/words/123") → BrowserRouter => actual path changes to "/words/123"
    // テストで window.location.pathname may or may not track this automatically,
    // but let's waitFor page to update:
    await waitFor(() => {
      // React Testing Library doesn't always keep "window.location.pathname" in sync
      // But let's see if we can read it
      expect(window.location.pathname).toBe('/words/123')
    })

    // WordShow で "updated-wordが正常に更新されました！" 表示
    expect(
      await screen.findByText('updated-wordが正常に更新されました！'),
    ).toBeInTheDocument()
  })
})
