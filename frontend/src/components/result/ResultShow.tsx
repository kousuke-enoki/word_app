import React, { useCallback, useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'

import PageBottomNav from '@/components/common/PageBottomNav'
import PageTitle from '@/components/common/PageTitle'
import ResultSettingCard from '@/components/result/ResultShow/ResultSettingCard'
import ResultTable from '@/components/result/ResultShow/ResultTable'
import { Badge } from '@/components/ui/card'
import { Card } from '@/components/ui/card'
import { useQuizResult } from '@/hooks/result/useQuizResult'
import { registerWord } from '@/service/word/RegisterWord'
import { ResultQuestion } from '@/types/quiz'

export default function ResultShow() {
  const { quizNo } = useParams<{ quizNo?: string }>()
  // const nav = useNavigate()
  const { loading, error, result } = useQuizResult(quizNo)

  const [size, setSize] = useState(10)
  const [page, setPage] = useState(0)
  const [view, setView] = useState<typeof result | null>(null)

  useEffect(() => setView(result ?? null), [result])

  // 総ページ数（>=1）
  const totalPages = useMemo(() => {
    if (!result) return 1
    return Math.max(1, Math.ceil(result.resultQuestions.length / size))
  }, [result, size])

  // ページサイズ変更やデータ到着で現在ページがはみ出したら補正
  useEffect(() => {
    if (page + 1 > totalPages) setPage(totalPages - 1)
  }, [totalPages, page])

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

  if (loading) return <p className="py-10 text-center">Loading...</p>
  if (error)
    return <p className="py-10 text-center text-red-500">通信に失敗しました</p>
  if (!result) return null

  // ページサイズ候補（問題数を超えるサイズは除外）
  const pageSizeOptions = [10, 20, 30, 50, 100].filter(
    (n) => n <= result.totalQuestionsCount,
  )

  return (
    <div className="mx-auto max-w-5xl">
      {/* タイトルは共通コンポーネント、rate等はこの画面で表示 */}
      <PageTitle title="クイズ結果" />

      {/* ここで正解数＆正解率を自由に配置可能 */}
      <div className="mb-4 flex items-center gap-2">
        <Badge>
          正解 {result.correctCount}/{result.totalQuestionsCount}
        </Badge>
        <Badge>{result.resultCorrectRate.toFixed(1)}%</Badge>
      </div>

      <ResultSettingCard setting={result.resultSetting} />

      <ResultTable
        rows={rows}
        onToggleRegister={toggleRegister}
        pager={{
          page: page + 1,
          totalPages,
          onPageChange: (p) => setPage(p - 1),
          pageSize: size,
          onPageSizeChange: (n) => {
            setSize(n)
            setPage(0)
          },
          pageSizeOptions,
          align: 'center',
        }}
      />

      <Card className="mt1 p-2">
        <PageBottomNav
          className="mt-1"
          actions={[{ label: 'クイズメニュー', to: '/quizs' }]}
          back={{ label: '成績一覧', to: '/results', variant: 'outline' }}
          showHome
          inline
          compact
        />
      </Card>
    </div>
  )
}
