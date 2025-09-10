// src/components/word/__tests__/WordList.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import WordList from '../WordList'
/* ---- UI を薄くモック（安定化） ---- */
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
  Input: ({ ...rest }: any) => <input {...rest} />,
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({ children, ...rest }: any) => <button {...rest}>{children}</button>,
}))

/* ---- 品詞データ固定 ---- */
vi.mock('@/service/word/GetPartOfSpeech', () => ({
  getPartOfSpeech: [
    { id: 1, name: '名詞' },
    { id: 2, name: '動詞' },
    { id: 3, name: '形容詞' },
  ],
}))

/* ---- registerWord をモック ---- */
vi.mock('@/service/word/RegisterWord', () => ({
  registerWord: vi.fn(),
}))

/* ---- axios をモック ---- */
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

// beforeEach(() => {
//   vi.useFakeTimers()
//   vi.resetAllMocks()
// })
beforeEach(() => {
  vi.useRealTimers()
  vi.resetAllMocks()
})

/* ---- テスト用データ ---- */
const page1 = {
  words: [
    {
      id: 1,
      name: 'apple',
      wordInfos: [
        {
          partOfSpeechId: 1,
          japaneseMeans: [{ name: 'りんご' }, { name: '林檎' }],
        },
      ],
      registrationCount: 2,
      isRegistered: false,
    },
    {
      id: 2,
      name: 'run',
      wordInfos: [{ partOfSpeechId: 2, japaneseMeans: [{ name: '走る' }] }],
      registrationCount: 5,
      isRegistered: true,
    },
  ],
  totalPages: 5,
}
const page2 = {
  words: [
    {
      id: 3,
      name: 'blue',
      wordInfos: [{ partOfSpeechId: 3, japaneseMeans: [{ name: '青い' }] }],
      registrationCount: 1,
      isRegistered: false,
    },
  ],
  totalPages: 5,
}
const filtered = {
  words: [
    {
      id: 4,
      name: 'apricot',
      wordInfos: [{ partOfSpeechId: 1, japaneseMeans: [{ name: 'あんず' }] }],
      registrationCount: 3,
      isRegistered: false,
    },
  ],
  totalPages: 2,
}

const renderList = (state?: any) =>
  render(
    <MemoryRouter initialEntries={[{ pathname: '/words', state }]}>
      <WordList />
    </MemoryRouter>,
  )

describe('WordList', () => {
  it('(1) デフォルト初期表示：一覧/リンク/ページ情報/パラメータが正しい', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({ data: page1 })

    renderList()

    // 一覧が描画されるまで待つ
    expect(
      await screen.findByRole('heading', { name: '単語一覧' }),
    ).toBeInTheDocument()
    // テーブルの1行目 apple の情報
    expect(screen.getByRole('link', { name: 'apple' })).toHaveAttribute(
      'href',
      '/words/1',
    )
    expect(screen.getByText('りんご, 林檎')).toBeInTheDocument() // 日本語訳 join
    expect(screen.getByText('名詞')).toBeInTheDocument() // 品詞名マッピング
    expect(screen.getByText('2')).toBeInTheDocument() // 登録数

    // ページ表示
    expect(screen.getByText('ページ 1 / 5')).toBeInTheDocument()

    // 初回の API パラメータ
    await waitFor(() => {
      const [path, opts] = (axiosInstance.get as any).mock.calls[0]
      expect(path).toBe('words')
      expect(opts.params).toMatchObject({
        search: '',
        sortBy: 'name',
        order: 'asc',
        page: 1,
        limit: 10,
      })
    })
  })

  it('(2) location.state から初期化される', async () => {
    ;(axiosInstance.get as any).mockResolvedValue({ data: page1 })
    renderList({
      search: 'ap',
      sortBy: 'registrationCount',
      order: 'desc',
      page: 3,
      limit: 20,
    })

    await screen.findByRole('heading', { name: '単語一覧' })

    await waitFor(() => {
      const last = (axiosInstance.get as any).mock.calls.at(-1)
      const [, opts] = last
      expect(opts.params).toMatchObject({
        search: 'ap',
        sortBy: 'registrationCount',
        order: 'desc',
        page: 3,
        limit: 20,
      })
    })
  })

  it('(3) 検索で再取得される（ページは維持）', async () => {
    ;(axiosInstance.get as any)
      .mockResolvedValueOnce({ data: page1 }) // 初回
      .mockResolvedValue({ data: filtered }) // 以降（a と p で2回来てもOK）

    renderList({ page: 2 })

    await screen.findByRole('heading', { name: '単語一覧' })
    await userEvent.type(screen.getByPlaceholderText('単語検索'), 'ap')

    await waitFor(() => {
      const [, opts] = (axiosInstance.get as any).mock.calls.at(-1)
      expect(opts.params).toMatchObject({ search: 'ap', page: 2 })
    })
    expect(
      await screen.findByRole('link', { name: 'apricot' }),
    ).toBeInTheDocument()
  })

  it('(4) ソート変更：register を選ぶと page=1 にリセットされる', async () => {
    ;(axiosInstance.get as any)
      .mockResolvedValueOnce({ data: page1 }) // 初回（page=2 指定で呼ばれる）
      .mockResolvedValueOnce({ data: page1 }) // register 変更後（page=1）

    renderList({ page: 2 })

    await screen.findByRole('heading', { name: '単語一覧' })
    const toolbar = screen
      .getByPlaceholderText('単語検索')
      .closest('div')!.parentElement!
    const sortSelect = within(toolbar).getByRole('combobox')

    // register を選ぶ
    await userEvent.selectOptions(sortSelect, 'register')

    // page=1 で再取得される
    await waitFor(() => {
      const last = (axiosInstance.get as any).mock.calls.at(-1)
      const [, opts] = last
      expect(opts.params).toMatchObject({ sortBy: 'register', page: 1 })
    })
    expect(screen.getByText('ページ 1 / 5')).toBeInTheDocument()
  })

  it('(5) 昇順/降順ボタンで order がトグルし再取得', async () => {
    ;(axiosInstance.get as any)
      .mockResolvedValueOnce({ data: page1 }) // 初回 asc
      .mockResolvedValueOnce({ data: page1 }) // 2回目 desc

    renderList()
    await screen.findByRole('heading', { name: '単語一覧' })

    // 初期は 昇順
    const orderBtn = screen.getByRole('button', { name: '昇順' })
    await userEvent.click(orderBtn)

    await waitFor(() => {
      const last = (axiosInstance.get as any).mock.calls.at(-1)
      const [, opts] = last
      expect(opts.params.order).toBe('desc')
    })
    // ボタン表示も変わる
    expect(screen.getByRole('button', { name: '降順' })).toBeInTheDocument()
  })

  it('(6) ページング：最初/前/次/最後 の活性と動作、limit 変更も再取得', async () => {
    ;(axiosInstance.get as any)
      .mockResolvedValueOnce({ data: page1 }) // 初回 page=1
      .mockResolvedValueOnce({ data: page2 }) // 次へ page=2
      .mockResolvedValueOnce({ data: page1 }) // 最後へ page=5（モック簡略）
      .mockResolvedValueOnce({ data: page1 }) // 最初へ page=1
      .mockResolvedValueOnce({ data: page1 }) // limit 20 に変更

    renderList()
    await screen.findByRole('heading', { name: '単語一覧' })

    const btnFirst = screen.getByRole('button', { name: '最初へ' })
    const btnPrev = screen.getByRole('button', { name: '前へ' })
    const btnNext = screen.getByRole('button', { name: '次へ' })
    const btnLast = screen.getByRole('button', { name: '最後へ' })

    // 初期は最初/前が disabled、次/最後が enabled
    expect(btnFirst).toBeDisabled()
    expect(btnPrev).toBeDisabled()
    expect(btnNext).toBeEnabled()
    expect(btnLast).toBeEnabled()

    // 次へ → page=2
    await userEvent.click(btnNext)
    await waitFor(() =>
      expect(screen.getByText('ページ 2 / 5')).toBeInTheDocument(),
    )

    // 最後へ → page=5
    await userEvent.click(btnLast)
    await waitFor(() => {
      const last = (axiosInstance.get as any).mock.calls.at(-1)
      const [, opts] = last
      expect(opts.params.page).toBe(5)
    })
    expect(screen.getByText('ページ 5 / 5')).toBeInTheDocument()

    // 最初へ → page=1
    await userEvent.click(btnFirst)
    await waitFor(() =>
      expect(screen.getByText('ページ 1 / 5')).toBeInTheDocument(),
    )

    // limit 変更（10→20）
    await userEvent.selectOptions(screen.getAllByRole('combobox')[1], '20')
    await waitFor(() => {
      const last = (axiosInstance.get as any).mock.calls.at(-1)
      const [, opts] = last
      expect(opts.params.limit).toBe(20)
    })
  })

  it('(8) words が空でも表は表示され、行がない（空配列）', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { words: [], totalPages: 1 },
    })
    renderList()
    await screen.findByRole('heading', { name: '単語一覧' })

    const table = screen.getByRole('table')
    expect(
      within(table).getByRole('columnheader', { name: '単語名' }),
    ).toBeInTheDocument()
    expect(within(table).queryByRole('link')).not.toBeInTheDocument() // ← テーブル内にリンク無し
    // ヘッダーの「新規登録」リンクは存在してOK
    expect(screen.getByRole('link', { name: '新規登録' })).toBeInTheDocument()
    expect(screen.getByText('ページ 1 / 1')).toBeInTheDocument()
  })
})
