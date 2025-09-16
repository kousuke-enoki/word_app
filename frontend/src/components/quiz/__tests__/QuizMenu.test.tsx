// src/components/quiz/__tests__/QuizMenu.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import QuizMenu from '../../quiz/QuizMenu' // 実際の相対パスに合わせて調整

/* ========= 依存モック ========= */
// axios
vi.mock('@/axiosConfig', () => ({
  default: { get: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

// useNavigate
const navigateMock = vi.fn()
vi.mock('react-router-dom', async () => {
  const actual =
    await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => navigateMock }
})

// 子コンポーネント（機能だけ持つ薄いモック）
vi.mock('../../quiz/QuizSettings', () => ({
  default: ({ settings, onSaveSettings }: any) => (
    <div data-testid="QuizSettings" data-qcount={settings?.questionCount}>
      <button
        onClick={() =>
          onSaveSettings({ ...settings, quizSettingCompleted: true })
        }
      >
        save-settings
      </button>
    </div>
  ),
}))
vi.mock('../../quiz/QuizStart', () => ({
  default: ({ onSuccess, onFail, settings }: any) => (
    <div data-testid="QuizStart" data-qcount={settings?.questionCount}>
      <button
        onClick={() =>
          onSuccess(123, {
            quizID: 123,
            questionNumber: 1,
            wordName: 'first',
            choicesJpms: [{ japaneseMeanID: 1, name: 'A' }],
          })
        }
      >
        create-success
      </button>
      <button onClick={() => onFail('作成に失敗しました')}>create-fail</button>
    </div>
  ),
}))
vi.mock('../../quiz/QuizQuestionView', () => ({
  default: ({ question, onAnswered }: any) => (
    <div data-testid="QuizQuestionView">
      <span data-testid="q-word">{question?.wordName}</span>
      <button
        onClick={() =>
          onAnswered({
            isFinish: false,
            nextQuestion: {
              quizID: question.quizID,
              questionNumber: question.questionNumber + 1,
              wordName: 'next-word',
              choicesJpms: [{ japaneseMeanID: 1, name: 'B' }],
            },
          })
        }
      >
        answer-continue
      </button>
      <button
        onClick={() => onAnswered({ isFinish: true, quizNumber: 777 } as any)}
      >
        answer-finish
      </button>
    </div>
  ),
}))

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter initialEntries={['/quizs']}>{ui}</MemoryRouter>)

beforeEach(() => {
  vi.clearAllMocks()
})

describe('QuizMenu', () => {
  it('初期: 進行中クイズあり → running で QuizQuestionView を表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        isRunningQuiz: true,
        nextQuestion: {
          quizID: 9,
          questionNumber: 3,
          wordName: 'resume-word',
          choicesJpms: [{ japaneseMeanID: 1, name: 'A' }],
        },
      },
    })

    renderWithRouter(<QuizMenu />)

    // QuizQuestionView が出て、現在の単語名が表示される
    const qv = await screen.findByTestId('QuizQuestionView')
    expect(qv).toBeInTheDocument()
    expect(within(qv).getByTestId('q-word')).toHaveTextContent('resume-word')
  })

  it('初期: 進行中クイズなし → setting で QuizSettings を表示', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { isRunningQuiz: false },
    })

    renderWithRouter(<QuizMenu />)
    const set = await screen.findByTestId('QuizSettings')
    expect(set).toBeInTheDocument()
    expect(set).toHaveAttribute('data-qcount', '10') // 初期の questionCount が渡っている
  })

  it('初期: 取得エラー → setting へフォールバック', async () => {
    ;(axiosInstance.get as any).mockRejectedValueOnce(new Error('NG'))
    renderWithRouter(<QuizMenu />)
    expect(await screen.findByTestId('QuizSettings')).toBeInTheDocument()
  })

  it('設定保存 → create 状態で QuizStart 表示 → 成功で running へ', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { isRunningQuiz: false },
    })
    renderWithRouter(<QuizMenu />)

    // setting 表示
    const set = await screen.findByTestId('QuizSettings')
    // 保存クリック → create へ
    await userEvent.click(within(set).getByText('save-settings'))

    // QuizStart 出現
    const start = await screen.findByTestId('QuizStart')
    expect(start).toBeInTheDocument()
    expect(start).toHaveAttribute('data-qcount', '10')

    // 作成成功で running へ
    await userEvent.click(within(start).getByText('create-success'))
    const qv = await screen.findByTestId('QuizQuestionView')
    expect(within(qv).getByTestId('q-word')).toHaveTextContent('first')
  })

  it('作成失敗: alert が呼ばれ setting に戻る', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: { isRunningQuiz: false },
    })
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})

    renderWithRouter(<QuizMenu />)
    const set = await screen.findByTestId('QuizSettings')
    await userEvent.click(within(set).getByText('save-settings'))
    const start = await screen.findByTestId('QuizStart')

    await userEvent.click(within(start).getByText('create-fail'))
    expect(alertSpy).toHaveBeenCalledWith('作成に失敗しました')

    // setting に戻っている
    expect(await screen.findByTestId('QuizSettings')).toBeInTheDocument()
    alertSpy.mockRestore()
  })

  it('回答: 継続なら次問に差し替え / 完了なら結果画面へ navigate', async () => {
    ;(axiosInstance.get as any).mockResolvedValueOnce({
      data: {
        isRunningQuiz: true,
        nextQuestion: {
          quizID: 55,
          questionNumber: 1,
          wordName: 'start-word',
          choicesJpms: [{ japaneseMeanID: 1, name: 'X' }],
        },
      },
    })
    renderWithRouter(<QuizMenu />)

    const qv = await screen.findByTestId('QuizQuestionView')
    expect(within(qv).getByTestId('q-word')).toHaveTextContent('start-word')

    // 継続
    await userEvent.click(within(qv).getByText('answer-continue'))
    const qv2 = await screen.findByTestId('QuizQuestionView')
    expect(within(qv2).getByTestId('q-word')).toHaveTextContent('next-word')

    // 完了
    await userEvent.click(within(qv2).getByText('answer-finish'))
    expect(navigateMock).toHaveBeenCalledWith('/results/777')
  })

  it('（現実装の仕様）loading=true 中は「Loading...」が描画されない', async () => {
    // get を未解決にして pending にする
    let resolver: (v: any) => void
    ;(axiosInstance.get as any).mockImplementationOnce(
      () => new Promise((res) => (resolver = res)),
    )
    renderWithRouter(<QuizMenu />)

    // 「Loading...」は出ない（真偽値 true が描画されるため）
    expect(screen.queryByText('Loading...')).toBeNull()

    // 後片付け: 解決させて setting に遷移
    resolver!({ data: { isRunningQuiz: false } })
    expect(await screen.findByTestId('QuizSettings')).toBeInTheDocument()
  })
})
