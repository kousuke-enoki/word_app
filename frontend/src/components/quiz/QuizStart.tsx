import React, { useEffect, useRef, useState } from 'react'

import axiosInstance from '../../axiosConfig'
import {
  CreateQuizResponse,
  QuizQuestion,
  QuizSettingsType,
} from '../../types/quiz'

interface Props {
  settings: QuizSettingsType
  /** 成功したら親にクイズ ID と 1 問目を渡す */
  onSuccess: (quizId: number, firstQuestion: QuizQuestion) => void
  onFail?: (msg: string) => void
}

const QuizStart: React.FC<Props> = ({ settings, onSuccess, onFail }) => {
  const [loading, setLoading] = useState(true)
  const called = useRef(false)

  useEffect(() => {
    if (called.current) return // ← 開発モード用の 2 回目はスキップ
    called.current = true
    /** 設定をそのまま POST 用 JSON に変換 */
    const payload = {
      questionCount: settings.questionCount,
      isSaveResult: settings.isSaveResult,
      isRegisteredWords: settings.isRegisteredWords,
      correctRate: settings.correctRate,
      attentionLevelList: settings.attentionLevelList,
      partsOfSpeeches: settings.partsOfSpeeches,
      isIdioms: settings.isIdioms,
      isSpecialCharacters: settings.isSpecialCharacters,
    }

    axiosInstance
      .post<CreateQuizResponse>('/quizzes/new', payload)
      .then((res) => {
        const { quizID, nextQuestion } = res.data
        /** フロント内部型に合わせて詰め替え */
        const first: QuizQuestion = {
          quizID: nextQuestion.quizID,
          questionNumber: nextQuestion.questionNumber,
          wordName: nextQuestion.wordName,
          choicesJpms: nextQuestion.choicesJpms.map((c) => ({
            japaneseMeanID: c.japaneseMeanID,
            name: c.name,
          })),
        }
        onSuccess(quizID, first)
      })
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      .catch((err: any) => {
        console.error(err)
        if (err?.response?.status === 429) {
          onFail?.('クイズ生成上限に達しました')
        } else {
          onFail?.('クイズ生成に失敗しました')
        }
      })
      .finally(() => setLoading(false))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []) // settings は確定後のみ呼ばれる想定

  if (loading) return <p>クイズを生成中です…</p>
  return null // 成功時は親側で state が切り替わり、ここは描画されなくなる
}

export default QuizStart
