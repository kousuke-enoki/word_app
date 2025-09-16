// src/components/quiz/__tests__/QuizQuestionView.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import QuizQuestionView from '../../quiz/QuizQuestionView'

/* ====== 依存モック ====== */
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

// UI ラッパは極薄に
vi.mock('@/components/ui/card', () => ({
  Card: ({ children, ...rest }: any) => (
    <div data-testid="Card" {...rest}>
      {children}
    </div>
  ),
  Badge: ({ children, ...rest }: any) => (
    <span data-testid="Badge" {...rest}>
      {children}
    </span>
  ),
}))
vi.mock('@/components/ui/ui', () => ({
  Button: ({
    children,
    onClick,
    disabled,
    className,
    variant,
    full,
    ...rest
  }: any) => (
    <button
      data-testid="UIButton"
      data-variant={variant || ''}
      data-full={String(!!full)}
      onClick={onClick}
      disabled={disabled}
      className={className}
      {...rest}
    >
      {children}
    </button>
  ),
}))
vi.mock('@/components/common/PageBottomNav', () => ({
  default: (p: any) => <div data-testid="PageBottomNav" {...p} />,
}))
vi.mock('@/components/common/PageTitle', () => ({
  default: ({ title }: any) => <h1 data-testid="PageTitle">{title}</h1>,
}))

/* ====== ダミー生成 ====== */
type ChoiceJpm = import('@/types/quiz').ChoiceJpm
type QuizQuestion = import('@/types/quiz').QuizQuestion
// type AnswerRouteRes = import('@/types/quiz').AnswerRouteRes

const q = (over: Partial<QuizQuestion> = {}): QuizQuestion => ({
  quizID: over.quizID ?? 1,
  questionNumber: over.questionNumber ?? 1,
  wordName: over.wordName ?? 'apple',
  choicesJpms:
    over.choicesJpms ??
    ([
      { japaneseMeanID: 11, name: '意味A' },
      { japaneseMeanID: 22, name: '意味B' },
    ] as ChoiceJpm[]),
})

const renderWith = (props: Parameters<typeof QuizQuestionView>[0]) =>
  render(<QuizQuestionView {...props} />)

beforeEach(() => {
  vi.clearAllMocks()
})

/* ====== テスト ====== */
describe('QuizQuestionView', () => {
  it('初期表示: タイトル/バッジ/単語/4ボタン（不足分は disabled）/OK は disabled', () => {
    const onAnswered = vi.fn()
    renderWith({ question: q(), onAnswered })

    // タイトルと Q番号
    expect(screen.getByTestId('PageTitle')).toHaveTextContent('クイズ')
    expect(screen.getByTestId('Badge')).toHaveTextContent('Q1')
    // 単語
    expect(screen.getByRole('heading', { name: 'apple' })).toBeInTheDocument()

    // 選択肢ボタンは常に 4 つ（2つは空で disabled）
    const choiceButtons = screen.getAllByTestId('UIButton').slice(0, 4)
    expect(choiceButtons).toHaveLength(4)
    expect(choiceButtons[0]).toHaveTextContent('意味A')
    expect(choiceButtons[1]).toHaveTextContent('意味B')
    expect(choiceButtons[2]).toHaveAttribute('disabled')
    expect(choiceButtons[3]).toHaveAttribute('disabled')

    // OK ボタンは最後の UIButton（実装上この順）
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    expect(okBtn).toHaveTextContent('OK')
    expect(okBtn).toBeDisabled()
  })

  it('クリックで選択 → OK enabled → 成功レスポンスで onAnswered 呼び出し', async () => {
    const onAnswered = vi.fn()
    // 手動で解決する Promise にして「送信中…」状態を観測する
    let resolvePost!: (v: any) => void
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolvePost = res
        }),
    )

    renderWith({ question: q(), onAnswered })

    const [, c2] = screen.getAllByTestId('UIButton').slice(0, 2) // ← 2個目を選択
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!

    // 選択して OK enable
    await userEvent.click(c2)
    expect(okBtn).toBeEnabled()

    // 送信中表示 → 戻る
    await userEvent.click(okBtn)
    await waitFor(() => expect(okBtn).toHaveTextContent('送信中…'))
    // 成功レスポンスを解決 → onAnswered 呼び出し＆文言戻る
    resolvePost({
      data: { isFinish: false, nextQuestion: { quizID: 1, questionNumber: 2 } },
    })
    await waitFor(() => expect(okBtn).toHaveTextContent('OK'))

    // 成功でコールされる
    expect(onAnswered).toHaveBeenCalledWith({
      isFinish: false,
      nextQuestion: { quizID: 1, questionNumber: 2 },
    })
    // 投稿後はボタン文言が戻る
    expect(okBtn).toHaveTextContent('OK')

    // POST 先と payload
    expect(axiosInstance.post).toHaveBeenCalledWith(
      '/quizzes/answers/1',
      expect.objectContaining({
        quizID: 1,
        questionNumber: 1,
        answerJpmID: 22, // 2つ目を選択（修正後は本当に2つ目）
      }),
    )
  })

  it('送信失敗: alert 表示・OK が再び押せる状態に戻る', async () => {
    const onAnswered = vi.fn()
    ;(axiosInstance.post as any).mockRejectedValueOnce(new Error('NG'))
    const alertSpy = vi.spyOn(window, 'alert').mockImplementation(() => {})

    renderWith({ question: q(), onAnswered })

    const c1 = screen.getAllByTestId('UIButton')[0]
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    await userEvent.click(c1)
    await userEvent.click(okBtn)

    expect(alertSpy).toHaveBeenCalledWith('回答送信に失敗しました')
    // 再び押せる
    expect(okBtn).toHaveTextContent('OK')
    expect(okBtn).toBeEnabled()

    alertSpy.mockRestore()
  })

  it('キーボード: 1〜4 で選択/Enter で送信、未選択 Enter は無視、posting 中は無視', async () => {
    const onAnswered = vi.fn()
    // 未解決の Promise で posting 中を維持
    let resolvePost: (v: any) => void
    ;(axiosInstance.post as any).mockImplementationOnce(
      () =>
        new Promise((res) => {
          resolvePost = res
        }),
    )

    renderWith({ question: q(), onAnswered })
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!

    // 未選択 Enter は無視
    await userEvent.keyboard('{Enter}')
    expect(axiosInstance.post).not.toHaveBeenCalled()

    // 「2」を押して選択
    await userEvent.keyboard('2')
    expect(okBtn).toBeEnabled()

    // Enter で送信開始（posting 中）
    await userEvent.keyboard('{Enter}')
    expect(axiosInstance.post).toHaveBeenCalledTimes(1)
    expect(okBtn).toHaveTextContent('送信中…')

    // posting 中はさらにキー入力を無視
    await userEvent.keyboard('1')
    await userEvent.keyboard('{Enter}')
    expect(axiosInstance.post).toHaveBeenCalledTimes(1)

    // 解決させて終了
    resolvePost!({ data: { isFinish: false } })
  })

  it('パディングされた空の選択肢はホットキー/クリックともに選択されない（OKは disabled のまま）', async () => {
    const onAnswered = vi.fn()
    renderWith({
      question: q({ choicesJpms: [{ japaneseMeanID: 1, name: 'A' }] }),
      onAnswered,
    })

    const choiceButtons = screen.getAllByTestId('UIButton').slice(0, 4)
    // 1つ目のみ有効、それ以外 disabled
    expect(choiceButtons[0]).not.toBeDisabled()
    expect(choiceButtons[1]).toBeDisabled()
    expect(choiceButtons[2]).toBeDisabled()
    expect(choiceButtons[3]).toBeDisabled()

    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    expect(okBtn).toBeDisabled()

    // キーボード「3」「4」は無視
    await userEvent.keyboard('3')
    await userEvent.keyboard('4')
    expect(okBtn).toBeDisabled()
  })

  it('質問が変わると（quizID or questionNumber）選択/投稿状態がリセットされる', async () => {
    const onAnswered = vi.fn()
    const { rerender } = renderWith({
      question: q({ quizID: 1, questionNumber: 1 }),
      onAnswered,
    })

    const c1 = screen.getAllByTestId('UIButton')[0]
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    await userEvent.click(c1)
    expect(okBtn).toBeEnabled()

    // questionNumber が変わる → リセット
    rerender(
      <QuizQuestionView
        question={q({ quizID: 1, questionNumber: 2 })}
        onAnswered={onAnswered}
      />,
    )
    const okBtn2 = screen.getAllByTestId('UIButton').at(-1)!
    expect(okBtn2).toBeDisabled()

    // quizID が変わる → リセット
    rerender(
      <QuizQuestionView
        question={q({ quizID: 2, questionNumber: 2 })}
        onAnswered={onAnswered}
      />,
    )
    const okBtn3 = screen.getAllByTestId('UIButton').at(-1)!
    expect(okBtn3).toBeDisabled()
  })

  it('質問のラベルだけ変わっても（wordName 変更）選択は維持される', async () => {
    const onAnswered = vi.fn()
    const { rerender } = renderWith({
      question: q({ wordName: 'alpha' }),
      onAnswered,
    })
    const c1 = screen.getAllByTestId('UIButton')[0]
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    await userEvent.click(c1)
    expect(okBtn).toBeEnabled()

    // IDs は同じでラベルだけ変更 → effect の依存に当たらず維持
    rerender(
      <QuizQuestionView
        question={q({ wordName: 'beta' })}
        onAnswered={onAnswered}
      />,
    )
    const okBtn2 = screen.getAllByTestId('UIButton').at(-1)!
    expect(okBtn2).toBeEnabled()
    expect(screen.getByRole('heading', { name: 'beta' })).toBeInTheDocument()
  })

  it('POST のペイロード: 選択肢3/4が存在するケースでも正しく送信', async () => {
    const onAnswered = vi.fn()
    ;(axiosInstance.post as any).mockResolvedValueOnce({
      data: { isFinish: true, quizNumber: 999 },
    })
    // 4つ埋まっているケース
    const question = q({
      choicesJpms: [
        { japaneseMeanID: 10, name: 'A' },
        { japaneseMeanID: 20, name: 'B' },
        { japaneseMeanID: 30, name: 'C' },
        { japaneseMeanID: 40, name: 'D' },
      ],
    })

    renderWith({ question, onAnswered })
    const choiceButtons = screen.getAllByTestId('UIButton').slice(0, 4)
    await userEvent.click(choiceButtons[3]) // D を選択
    const okBtn = screen.getAllByTestId('UIButton').at(-1)!
    await userEvent.click(okBtn)

    expect(axiosInstance.post).toHaveBeenCalledWith(
      '/quizzes/answers/1',
      expect.objectContaining({ answerJpmID: 40 }),
    )
    expect(onAnswered).toHaveBeenCalledWith({ isFinish: true, quizNumber: 999 })
  })
})
