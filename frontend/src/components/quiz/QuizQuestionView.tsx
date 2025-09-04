import React, { useEffect, useState } from 'react'

import axiosInstance from '@/axiosConfig'
import { Badge, Card } from '@/components/card'
import { Button } from '@/components/ui'
import type { AnswerRouteRes, ChoiceJpm, QuizQuestion } from '@/types/quiz'

type Props = {
  question: QuizQuestion
  onAnswered: (res: AnswerRouteRes) => void
}

const QuizQuestionView: React.FC<Props> = ({ question, onAnswered }) => {
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [posting, setPosting] = useState(false)

  useEffect(() => {
    setSelectedId(null)
    setPosting(false)
  }, [question.quizID, question.questionNumber])

  const handleSubmit = async () => {
    if (selectedId == null || posting) return
    setPosting(true)
    try {
      const payload = {
        quizID: question.quizID,
        answerJpmID: selectedId,
        questionNumber: question.questionNumber,
      }
      const res = await axiosInstance.post<AnswerRouteRes>(
        `/quizzes/answers/${question.quizID}`,
        payload,
      )
      onAnswered(res.data)
    } catch (e) {
      console.error(e)
      alert('回答送信に失敗しました')
    } finally {
      setPosting(false)
    }
  }

  // 2×2に整形（空を-1で埋める）
  const gridChoices: ChoiceJpm[] = [...question.choicesJpms]
    .concat(new Array(4).fill({ japaneseMeanID: -1, name: '' }))
    .slice(0, 4)

  // ---- keyboard: 1/2/3/4 で選択、Enter で送信 ----
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (posting) return
      // 1..4
      if (e.key >= '1' && e.key <= '4') {
        const idx = Number(e.key) - 1
        const c = gridChoices[idx]
        if (c && c.japaneseMeanID !== -1) {
          setSelectedId(c.japaneseMeanID)
        }
      }
      // Enter
      if (e.key === 'Enter') {
        if (selectedId != null) {
          e.preventDefault()
          handleSubmit()
        }
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [gridChoices, selectedId, posting])

  return (
    <div className="mx-auto max-w-2xl">
      <Card className="p-6">
        <div className="mb-4 flex items-center justify-between">
          <Badge>Q{question.questionNumber}</Badge>
        </div>

        <h1 className="mb-6 text-center text-4xl font-extrabold">
          {question.wordName}
        </h1>

        {/* 選択肢（ボタンを大きく） */}
        <div className="grid grid-cols-2 gap-3">
          {gridChoices.map((c, i) => {
            const disabled = c.japaneseMeanID === -1
            const isSelected = selectedId === c.japaneseMeanID
            const hotkey = String(i + 1) // 1..4
            return (
              <div key={`${question.quizID}-${i}`} className="relative">
                {/* 左上にホットキーの丸ラベル（任意） */}
                <span className="pointer-events-none absolute left-2 top-2 inline-flex h-6 w-6 items-center justify-center rounded-full border border-[var(--btn-subtle-bd)] bg-[var(--btn-subtle-bg)] text-xs">
                  {hotkey}
                </span>
                <Button
                  variant={isSelected ? 'primary' : 'outline'}
                  full
                  disabled={disabled}
                  onClick={() => setSelectedId(c.japaneseMeanID)}
                  className="
                    h-20 sm:h-24                /* ← 縦を拡大 */
                    text-base sm:text-lg         /* ← 文字も少し大きく */
                    px-4                         /* 横paddingは控えめ（full幅） */
                  "
                >
                  {c.name || '　'}
                </Button>
              </div>
            )
          })}
        </div>

        {/* OKボタン（縦も横も拡大） */}
        <div className="mt-8 text-center">
          <Button
            onClick={handleSubmit}
            disabled={selectedId == null || posting}
            className="
              h-14 px-12                 /* ← 大きめ */
              text-lg                    /* ← 文字も大きく */
              min-w-[200px]              /* ← 横幅の基準を確保 */
            "
          >
            {posting ? '送信中…' : 'OK'}
          </Button>
          <p className="mt-2 text-sm opacity-70">
            ショートカット： 1 / 2 / 3 / 4 で選択、Enter で送信
          </p>
        </div>
      </Card>
    </div>
  )
}

export default QuizQuestionView
