import React, { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import Pagination from '@/components/common/Pagination'
import ResultHeader from '@/components/result/ResultShow/ResultHeader'
import ResultSettingCard from '@/components/result/ResultShow/ResultSettingCard'
import ResultTable from '@/components/result/ResultShow/ResultTable'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { useQuizResult } from '@/hooks/result/useQuizResult'
import { registerWord } from '@/service/word/RegisterWord'
import { ResultQuestion } from '@/types/quiz'

const pageSizes = [10, 30, 50, 100] as const
type PageSize = (typeof pageSizes)[number]

export default function ResultShow() {
  const { quizNo } = useParams<{ quizNo?: string }>()
  const nav = useNavigate()
  const { loading, error, result } = useQuizResult(quizNo)

  const [size, setSize] = useState<PageSize>(10)
  const [page, setPage] = useState(0)
  const [view, setView] = useState<typeof result | null>(null)

  useEffect(() => setView(result ?? null), [result])

  const rows = useMemo(() => {
    if (!view) return []
    const start = page * size
    return view.resultQuestions.slice(start, start + size)
  }, [view, page, size])

  const toggleRegister = useCallback(async (row: ResultQuestion) => {
    try {
      const u = await registerWord(row.wordID, !row.registeredWord.isRegistered)
      setView((prev) => {
        if (!prev) return prev
        return {
          ...prev,
          resultQuestions: prev.resultQuestions.map((q) =>
            q.wordID === row.wordID
              ? {
                  ...q,
                  registeredWord: {
                    ...q.registeredWord,
                    isRegistered: u.isRegistered,
                    quizCount: u.quizCount ?? q.registeredWord.quizCount,
                    correctCount:
                      u.correctCount ?? q.registeredWord.correctCount,
                  },
                }
              : q,
          ),
        }
      })
    } catch (e) {
      console.error(e)
    }
  }, [])

  const sizeCandidates: PageSize[] = useMemo(() => {
    if (!result) return []
    return pageSizes.filter(
      (s): s is PageSize => s <= result.totalQuestionsCount,
    )
  }, [result])

  if (loading) return <p className="py-10 text-center">Loading...</p>
  if (error)
    return <p className="py-10 text-center text-red-500">通信に失敗しました</p>
  if (!result) return null

  return (
    <div className="mx-auto max-w-5xl">
      <ResultHeader
        correct={result.correctCount}
        total={result.totalQuestionsCount}
        rate={result.resultCorrectRate}
      />

      <ResultSettingCard setting={result.resultSetting} />

      <ResultTable rows={rows} onToggleRegister={toggleRegister} />

      {/* ページング */}
      <Card className="mt-4 p-4">
        <Pagination
          sizes={sizeCandidates}
          size={size}
          page={page}
          total={result.resultQuestions.length}
          onSize={(s: PageSize) => {
            setSize(s)
            setPage(0)
          }}
          onPrev={() => setPage((p) => Math.max(0, p - 1))}
          onNext={() => setPage((p) => p + 1)}
        />

        {/* 下部ナビ */}
        <div className="mt-6 flex justify-center gap-3">
          <Button onClick={() => nav('/quizs')}>クイズメニュー</Button>
          <Button variant="outline" onClick={() => nav('/results')}>
            成績一覧
          </Button>
          <Button variant="ghost" onClick={() => nav('/')}>
            ホーム
          </Button>
        </div>
      </Card>
    </div>
  )
}
