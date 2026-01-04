import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { rest } from 'msw'
import React from 'react'
import { useLocation } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'

import { server } from '@/__tests__/mswServer'
import { renderWithClient } from '@/__tests__/testUtils'

import WordNew from '../WordNew'

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
        expect(body).toEqual({
          name: 'apple',
          wordInfos: [{ partOfSpeechId: 1, japaneseMeans: [{ name: 'りんご' }] }],
        })
        return res(ctx.status(200), ctx.json({ id: 21, name: 'apple' }))
      }),
    )

    renderWordNew()
    await submitValidForm()

    expect(
      await screen.findByTestId('success-message'),
    ).toHaveTextContent('appleが正常に登録されました！')
  })

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

    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'),
    ).toBeInTheDocument()
    expect(screen.queryByTestId('success-message')).toBeNull()
    expect(
      screen.getByRole('heading', { name: '単語登録' }),
    ).toBeInTheDocument()
  })

  it('500: fallback error message is shown', async () => {
    server.use(
      rest.post('http://localhost:8080/words/new', (_, res, ctx) =>
        res(ctx.status(500)),
      ),
    )

    renderWordNew()
    await submitValidForm()

    expect(
      await screen.findByText('単語の登録中にエラーが発生しました。'),
    ).toBeInTheDocument()
    expect(screen.queryByTestId('success-message')).toBeNull()
  })
})
