/**
 * WordNew コンポーネントの統合テスト
 * MSW（Mock Service Worker）を使用して実際のHTTPリクエストフローとナビゲーションを含むエンドツーエンドのシナリオをテスト
 */
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'

import { server } from '@/__tests__/mswServer'
import { renderWithClient } from '@/__tests__/testUtils'

import WordNew from '../WordNew'

/**
 * リダイレクト先のモックコンポーネント
 * location.stateから受け取ったsuccessMessageを表示し、リダイレクト後の状態を検証するために使用
 */
const WordShowSpy = () => {
  const location = useLocation()
  return (
    <div>
      <p>Word Show Dummy</p>
      <p data-testid="success-message">
        {(location.state as { successMessage?: string } | null)
          ?.successMessage ?? ''}
      </p>
    </div>
  )
}

/**
 * WordNewコンポーネントをテスト用にレンダリング
 * リダイレクト先（/words/:id）も含めたルーティングを設定
 */
const renderWordNew = () =>
  renderWithClient(<WordNew />, {
    router: {
      initialEntries: ['/words/new'],
      routes: [
        { path: '/words/new', element: <WordNew /> },
        { path: '/words/:id', element: <WordShowSpy /> },
      ],
    },
  })

/**
 * 有効なフォームデータを入力して送信するヘルパー関数
 * 単語名: 'apple', 品詞ID: 1, 日本語訳: 'りんご'
 */
const submitValidForm = async () => {
  const user = userEvent.setup()
  await user.type(screen.getByPlaceholderText('example'), 'apple')
  await user.selectOptions(screen.getByRole('combobox'), '1')
  await user.type(screen.getByPlaceholderText('意味'), 'りんご')
  await user.click(screen.getByRole('button', { name: '単語を登録する' }))
}

describe('WordNew integration (MSW)', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  /**
   * 成功フロー: フォーム送信 → API呼び出し → リダイレクト → 成功メッセージ表示
   * - 送信されたペイロードが正しいことを検証
   * - リダイレクト先で成功メッセージが表示されることを検証
   */
  it('success: posts payload then redirects with success message', async () => {
    server.use(
      rest.post('http://localhost:8080/words/new', async (req, res, ctx) => {
        const body = (await req.json()) as {
          name: string
          wordInfos: Array<{
            partOfSpeechId: number
            japaneseMeans: Array<{ name: string }>
          }>
        }
        // 送信されたペイロードが期待値と一致することを検証
        expect(body).toEqual({
          name: 'apple',
          wordInfos: [
            { partOfSpeechId: 1, japaneseMeans: [{ name: 'りんご' }] },
          ],
        })
        return res(ctx.status(200), ctx.json({ id: 21, name: 'apple' }))
      }),
    )

    renderWordNew()
    await submitValidForm()

    // リダイレクト先で成功メッセージが表示されることを検証
    expect(await screen.findByTestId('success-message')).toHaveTextContent(
      'appleが正常に登録されました！',
    )
  })

  /**
   * 422エラー（バリデーションエラー）時のフロー
   * - エラーメッセージが表示されることを検証
   * - リダイレクトされず、フォーム画面に留まることを検証
   */
  it('422: shows error message and stays on form', async () => {
    server.use(
      rest.post('http://localhost:8080/words/new', (_, res, ctx) =>
        res(
          ctx.status(422),
          ctx.json({
            errors: [{ field: 'name', message: 'name is invalid' }],
          }),
        ),
      ),
    )

    renderWordNew()
    await submitValidForm()

    // エラーメッセージが表示されることを検証
    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'),
    ).toBeInTheDocument()
    // リダイレクトされていないことを検証（成功メッセージは表示されない）
    expect(screen.queryByTestId('success-message')).toBeNull()
    // フォーム画面に留まっていることを検証
    expect(
      screen.getByRole('heading', { name: '単語登録' }),
    ).toBeInTheDocument()
  })

  /**
   * 500エラー（サーバーエラー）時のフロー
   * - フォールバックエラーメッセージが表示されることを検証
   * - リダイレクトされないことを検証
   */
  it('500: fallback error message is shown', async () => {
    server.use(
      rest.post('http://localhost:8080/words/new', (_, res, ctx) =>
        res(ctx.status(500)),
      ),
    )

    renderWordNew()
    await submitValidForm()

    // フォールバックエラーメッセージが表示されることを検証
    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'),
    ).toBeInTheDocument()
    // リダイレクトされていないことを検証
    expect(screen.queryByTestId('success-message')).toBeNull()
  })
})
