import {
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it } from 'vitest'

import { server } from '@/__tests__/mswServer'
import UserList from '@/components/user/UserList'

type UserDetail = import('@/components/user/UserList').UserDetail

const makeUser = (id: number, name: string, overrides?: Partial<UserDetail>) =>
  ({
    id,
    name,
    email: `${name}@example.com`,
    isAdmin: false,
    isRoot: false,
    isTest: false,
    isLine: false,
    isSettedPassword: true,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
    ...overrides,
  }) as UserDetail

const renderList = () =>
  render(
    <MemoryRouter initialEntries={[{ pathname: '/users' }]}>
      <UserList />
    </MemoryRouter>,
  )

const endpoint = 'http://localhost:8080/users'

describe('UserList integration', () => {
  it('初期ロードでスケルトンが表示され、解決後に一覧が描画される', async () => {
    server.use(
      rest.get(endpoint, (_req, res, ctx) => {
        return res(
          ctx.delay(50),
          ctx.status(200),
          ctx.json({ users: [], totalPages: 1 }),
        )
      }),
    )

    renderList()

    expect(screen.getByTestId('users-loading')).toBeInTheDocument()

    expect(await screen.findByText('総ページ: 1')).toBeInTheDocument()
    await waitFor(() =>
      expect(screen.queryByTestId('users-loading')).not.toBeInTheDocument(),
    )
  })

  it('正常レスポンス: 行数と主要項目が描画される', async () => {
    const list = [
      makeUser(1, 'alice', { isAdmin: true }),
      makeUser(2, 'bob', { isLine: true }),
    ]
    server.use(
      rest.get(endpoint, (_req, res, ctx) =>
        res(
          ctx.status(200),
          ctx.json({
            users: list,
            totalPages: 2,
          }),
        ),
      ),
    )

    renderList()

    expect(await screen.findByText('総ページ: 2')).toBeInTheDocument()
    const rows = screen.getAllByRole('row').slice(1)
    expect(rows).toHaveLength(2)
    expect(screen.getByText('alice')).toBeInTheDocument()
    expect(screen.getByText('bob')).toBeInTheDocument()

    const aliceRow = rows.find((r) => r.textContent?.includes('alice'))!
    const scoped = within(aliceRow)
    scoped.getByText('alice@example.com')
    scoped.getByText('Admin')
  })

  it('検索/ソート/ページング操作でクエリが反映され再取得される', async () => {
    const user = userEvent.setup()
    const page1 = [makeUser(1, 'Alice'), makeUser(2, 'Bob')]
    const page2 = [makeUser(3, 'Carol')]
    const filtered = [makeUser(4, 'Bobcat')]

    let lastParams: Record<string, string | null> = {}
    let callCount = 0
    server.use(
      rest.get(endpoint, (req, res, ctx) => {
        callCount += 1
        lastParams = {
          search: req.url.searchParams.get('search'),
          sortBy: req.url.searchParams.get('sortBy'),
          order: req.url.searchParams.get('order'),
          page: req.url.searchParams.get('page'),
          limit: req.url.searchParams.get('limit'),
        }
        const search = req.url.searchParams.get('search') ?? ''
        const page = Number(req.url.searchParams.get('page') ?? '1')
        const payload =
          search.toLowerCase().includes('bob') && filtered.length > 0
            ? { users: filtered, totalPages: 1 }
            : page === 2
              ? { users: page2, totalPages: 3 }
              : { users: page1, totalPages: 3 }

        return res(ctx.status(200), ctx.json(payload))
      }),
    )

    renderList()

    expect(await screen.findByText('総ページ: 3')).toBeInTheDocument()
    expect(lastParams).toMatchObject({
      search: '',
      sortBy: 'name',
      order: 'asc',
      page: '1',
      limit: '10',
    })

    const searchInput = screen.getByPlaceholderText('ユーザー名・メール検索')
    await user.type(searchInput, 'Bob')
    await waitFor(() =>
      expect(lastParams).toMatchObject({ search: 'Bob', page: '1' }),
    )
    expect(screen.getByText('Bobcat')).toBeInTheDocument()

    await user.clear(searchInput)
    await waitFor(() =>
      expect(lastParams).toMatchObject({ search: '', page: '1' }),
    )
    expect(screen.getByText('ページ 1 / 3')).toBeInTheDocument()

    const selects = screen.getAllByRole('combobox') as HTMLSelectElement[]
    const sortSelect =
      selects.find((sel) =>
        Array.from(sel.options).some((o) => o.value === 'role'),
      ) ?? selects[0]
    const pageSizeSelect = screen
      .getAllByTestId('pagination-page-size')
      .find((el) => !el.closest('[aria-hidden="true"]')) as HTMLSelectElement

    await user.selectOptions(sortSelect, 'role')
    await waitFor(() =>
      expect(lastParams).toMatchObject({ sortBy: 'role', page: '1' }),
    )

    await user.click(screen.getByRole('button', { name: '昇順' }))
    await waitFor(() =>
      expect(lastParams).toMatchObject({ order: 'desc', page: '1' }),
    )

    await user.click(screen.getByRole('button', { name: '次へ' }))
    expect(await screen.findByText('ページ 2 / 3')).toBeInTheDocument()
    expect(lastParams).toMatchObject({ page: '2', limit: '10' })

    fireEvent.change(pageSizeSelect, { target: { value: '20' } })
    expect(pageSizeSelect).toHaveValue('20')
    const beforeRefetch = callCount
    await user.type(searchInput, 'A')
    await waitFor(() => expect(callCount).toBeGreaterThan(beforeRefetch))
    await waitFor(() => expect(lastParams.search).toBe('A'))
  })

  it('空データと 500 エラーでそれぞれの表示になる', async () => {
    // 空データ
    server.use(
      rest.get(endpoint, (_req, res, ctx) =>
        res(ctx.status(200), ctx.json({ users: [], totalPages: 1 })),
      ),
    )

    renderList()

    expect(await screen.findByText('総ページ: 1')).toBeInTheDocument()
    expect(
      screen.getByText('該当するユーザーが見つかりませんでした'),
    ).toBeInTheDocument()

    // 500 エラー
    server.use(rest.get(endpoint, (_req, res, ctx) => res(ctx.status(500))))

    const user = userEvent.setup()
    await user.type(screen.getByPlaceholderText('ユーザー名・メール検索'), 'x')

    await waitFor(
      () =>
        expect(
          screen.getAllByText('ユーザー取得に失敗しました。').length,
        ).toBeGreaterThan(0),
      { timeout: 3000 },
    )
  })
})
