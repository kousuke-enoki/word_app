// src/components/quiz/__tests__/QuizStart.test.tsx
/* eslint-disable @typescript-eslint/no-explicit-any */
import { render, screen, waitFor } from '@testing-library/react'
import React from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import type { CreateQuizResponse, QuizSettingsType } from '@/types/quiz'

import QuizStart from '../QuizStart'

// axios をモック
vi.mock('@/axiosConfig', () => ({
  default: { post: vi.fn() },
}))
import axiosInstance from '@/axiosConfig'

const baseSettings: QuizSettingsType = {
  quizSettingCompleted: true, // QuizStart 側では送らない点も後で検証
  questionCount: 10,
  isSaveResult: true,
  isRegisteredWords: 0,
  correctRate: 100,
  attentionLevelList: [1, 2, 3, 4, 5],
  partsOfSpeeches: [1, 3, 4, 5],
  isIdioms: 0,
  isSpecialCharacters: 0,
}

const serverOk = (
  over: Partial<CreateQuizResponse> = {},
): CreateQuizResponse => ({
  quizID: 777,
  totalCreateQuestion: 10,
  nextQuestion: {
    quizID: 777,
    questionNumber: 1,
    wordName: 'apple',
    choicesJpms: [
      { japaneseMeanID: 10, name: '意味A' },
      { japaneseMeanID: 20, name: '意味B' },
      { japaneseMeanID: 30, name: '意味C' },
    ],
  },
  ...over,
})

beforeEach(() => {
  vi.clearAllMocks()
})

/** 遅延解決する Promise を作って、ローディング表示の検証をしやすくする */
const deferred = <T,>() => {
  let resolve!: (v: T) => void
  let reject!: (e: unknown) => void
  const p = new Promise<T>((res, rej) => {
    resolve = res
    reject = rej
  })
  return { p, resolve, reject }
}

describe('QuizStart', () => {
  it('初期表示はローディング。「クイズを生成中です…」が出る', () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({ data: serverOk() })
    render(
      <QuizStart
        settings={baseSettings}
        onSuccess={() => {}}
        onFail={() => {}}
      />,
    )
    expect(screen.getByText('クイズを生成中です…')).toBeInTheDocument()
  })

  it('成功: 正しいペイロードで POST → onSuccess(quizID, firstQuestion) が呼ばれ、ローディングは消える', async () => {
    const def = deferred<{ data: CreateQuizResponse }>()
    ;(axiosInstance.post as any).mockImplementationOnce(() => def.p)

    const onSuccess = vi.fn()
    render(<QuizStart settings={baseSettings} onSuccess={onSuccess} />)

    // まずはローディングが出ている
    expect(screen.getByText('クイズを生成中です…')).toBeInTheDocument()

    // POST の引数（エンドポイントとペイロード）をチェック
    await waitFor(() => expect(axiosInstance.post).toHaveBeenCalledTimes(1))
    const [url, payload] = (axiosInstance.post as any).mock.calls[0]
    expect(url).toBe('/quizzes/new')
    // settings から quizSettingCompleted を除いたフィールドのみ送っていること
    expect(payload).toEqual({
      questionCount: 10,
      isSaveResult: true,
      isRegisteredWords: 0,
      correctRate: 100,
      attentionLevelList: [1, 2, 3, 4, 5],
      partsOfSpeeches: [1, 3, 4, 5],
      isIdioms: 0,
      isSpecialCharacters: 0,
    })

    // サーバ成功を解決
    def.resolve({ data: serverOk() })

    // onSuccess が正しい形で呼ばれる
    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1))
    const [quizID, first] = onSuccess.mock.calls[0]
    expect(quizID).toBe(777)
    expect(first).toEqual({
      quizID: 777,
      questionNumber: 1,
      wordName: 'apple',
      choicesJpms: [
        { japaneseMeanID: 10, name: '意味A' },
        { japaneseMeanID: 20, name: '意味B' },
        { japaneseMeanID: 30, name: '意味C' },
      ],
    })

    // ローディングは消える（コンポーネントは null を返す想定）
    await waitFor(() =>
      expect(screen.queryByText('クイズを生成中です…')).toBeNull(),
    )
  })

  it('失敗: onFail があれば呼ばれる。onSuccess は呼ばれず、ローディングは消える', async () => {
    const def = deferred<any>()
    ;(axiosInstance.post as any).mockImplementationOnce(() => def.p)

    const onSuccess = vi.fn()
    const onFail = vi.fn()
    const spyErr = vi.spyOn(console, 'error').mockImplementation(() => {})

    render(
      <QuizStart
        settings={baseSettings}
        onSuccess={onSuccess}
        onFail={onFail}
      />,
    )

    // reject（通信エラー）
    def.reject(new Error('boom'))

    await waitFor(() =>
      expect(onFail).toHaveBeenCalledWith('クイズ生成に失敗しました'),
    )
    expect(onSuccess).not.toHaveBeenCalled()

    await waitFor(() =>
      expect(screen.queryByText('クイズを生成中です…')).toBeNull(),
    )
    spyErr.mockRestore()
  })

  it('失敗: onFail が無くても落ちない（console.error は呼ばれるが例外なし）', async () => {
    const def = deferred<any>()
    ;(axiosInstance.post as any).mockImplementationOnce(() => def.p)

    const onSuccess = vi.fn()
    const spyErr = vi.spyOn(console, 'error').mockImplementation(() => {})

    render(<QuizStart settings={baseSettings} onSuccess={onSuccess} />)

    def.reject(new Error('boom'))

    // onSuccess は呼ばれない
    await waitFor(() => expect(onSuccess).not.toHaveBeenCalled())
    // ローディングが消える
    await waitFor(() =>
      expect(screen.queryByText('クイズを生成中です…')).toBeNull(),
    )
    spyErr.mockRestore()
  })

  it('React.StrictMode 下でも POST は 1 回のみ（副作用二重実行のガードが効く）', async () => {
    ;(axiosInstance.post as any).mockResolvedValueOnce({ data: serverOk() })

    const onSuccess = vi.fn()
    render(
      <React.StrictMode>
        <QuizStart settings={baseSettings} onSuccess={onSuccess} />
      </React.StrictMode>,
    )

    await waitFor(() => expect(axiosInstance.post).toHaveBeenCalledTimes(1))
    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1))
  })

  it('サーバーからの nextQuestion が任意に変わっても詰め替えロジックが守られる', async () => {
    // choices に4件＆名前違い
    const data = serverOk({
      nextQuestion: {
        quizID: 999,
        questionNumber: 5,
        wordName: 'banana',
        choicesJpms: [
          { japaneseMeanID: 1, name: '意味1' },
          { japaneseMeanID: 2, name: '意味2' },
          { japaneseMeanID: 3, name: '意味3' },
          { japaneseMeanID: 4, name: '意味4' },
        ],
      },
    })

    ;(axiosInstance.post as any).mockResolvedValueOnce({ data })

    const onSuccess = vi.fn()
    render(
      <QuizStart
        settings={{ ...baseSettings, questionCount: 20 }}
        onSuccess={onSuccess}
      />,
    )

    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1))
    const [qid, first] = onSuccess.mock.calls[0]
    expect(qid).toBe(777) // quizID は top-level の方
    expect(first).toEqual({
      quizID: 999, // nextQuestion の方
      questionNumber: 5,
      wordName: 'banana',
      choicesJpms: [
        { japaneseMeanID: 1, name: '意味1' },
        { japaneseMeanID: 2, name: '意味2' },
        { japaneseMeanID: 3, name: '意味3' },
        { japaneseMeanID: 4, name: '意味4' },
      ],
    })
  })
})
