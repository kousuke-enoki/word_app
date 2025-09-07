import React, { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import type { ResultSummary } from '@/types/result'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'

/* 品詞 ID → 名称 */
const POS_MAP: Record<number, string> = {
  1: '名詞',
  2: '代名',
  3: '動詞',
  4: '形容',
  5: '副詞',
}

const PAGE_SIZES = [10, 20, 30] as const

const ResultIndex: React.FC = () => {
  const nav = useNavigate()

  const [list, setList] = useState<ResultSummary[]>([])
  const [pageSize, setPageSize] =
    useState<(typeof PAGE_SIZES)[number]>(10)
  const [page, setPage] = useState(0) // 0-based
  const [loading, setLoading] = useState(true)
  const [errMsg, setErrMsg] = useState('')

  useEffect(() => {
    const fetchAll = async () => {
      try {
        const res = await axiosInstance.get<ResultSummary[]>('/results')
        setList(
          res.data.sort(
            (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
          ),
        )
      } catch (e) {
        console.error(e)
        setErrMsg('成績の取得に失敗しました')
      } finally {
        setLoading(false)
      }
    }
    fetchAll()
  }, [])

  const totalPages = Math.max(1, Math.ceil(list.length / pageSize))
  const start = page * pageSize
  const end = Math.min(start + pageSize, list.length)

  const rows = useMemo(() => list.slice(start, end), [list, start, end])

  if (loading) return <p className="p-4">読み込み中…</p>
  if (errMsg) return <p className="p-4 text-red-600">{errMsg}</p>

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-4 flex items-end justify-between gap-3">
        <h1 className="text-2xl font-bold">成績一覧</h1>
        <span className="text-sm opacity-70">全 {list.length} 件</span>
      </div>

      <Card className="p-4 sm:p-6">
        {/* コントロール行 */}
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
          {/* ページサイズ切替（セグメント） */}
          <div className="inline-flex rounded-xl p-0.5 gap-0.5 bg-[var(--btn-subtle-bd)]">
            {PAGE_SIZES.map((s) => (
              <Button
                key={s}
                variant={s === pageSize ? 'primary' : 'ghost'}
                onClick={() => {
                  setPageSize(s)
                  setPage(0)
                }}
                className="rounded-lg px-3 py-1.5"
              >
                {s} 件
              </Button>
            ))}
          </div>

          {/* ページャ */}
          <div className="flex items-center gap-2">
            <span className="text-sm opacity-70">
              {list.length === 0 ? '0' : `${start + 1}`}–{end} / {list.length}
            </span>
            <Button
              variant="outline"
              disabled={page === 0}
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              className="px-3"
            >
              Prev
            </Button>
            <Button
              variant="outline"
              disabled={page + 1 >= totalPages}
              onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
              className="px-3"
            >
              Next
            </Button>
          </div>
        </div>

        {/* テーブル */}
        <div className="overflow-x-auto">
          <table className="min-w-[720px] w-full border-collapse text-sm">
            <thead>
              <tr className="bg-[var(--table_th)] text-[var(--table_th_c)]">
                <th className="px-3 py-2 text-left font-semibold rounded-l-lg">#</th>
                <th className="px-3 py-2 text-left font-semibold whitespace-nowrap">日付</th>
                <th className="px-3 py-2 text-left font-semibold">登録単語</th>
                <th className="px-3 py-2 text-left font-semibold">慣用句</th>
                <th className="px-3 py-2 text-left font-semibold">特殊</th>
                <th className="px-3 py-2 text-left font-semibold">品詞</th>
                <th className="px-3 py-2 text-right font-semibold">問題</th>
                <th className="px-3 py-2 text-right font-semibold">正解</th>
                <th className="px-3 py-2 text-right font-semibold rounded-r-lg">正解率</th>
              </tr>
            </thead>

            <tbody className="divide-y divide-[var(--border)]">
              {rows.length === 0 ? (
                <tr>
                  <td colSpan={9} className="px-3 py-10 text-center opacity-70">
                    成績がありません
                  </td>
                </tr>
              ) : (
                rows.map((r) => (
                  <tr
                    key={r.quizNumber}
                    onClick={() => nav(`/results/${r.quizNumber}`)}
                    tabIndex={0}
                    className="
                      group cursor-pointer transition-colors duration-150
                      even:[&>td]:bg-[var(--table_tr_e)]              /* 偶数行の縞も td に適用 */
                      hover:[&>td]:bg-[var(--table_row_hover)]        /* ← ホバー時は td を塗る */
                      active:[&>td]:bg-[var(--table_row_active)]
                      focus-visible:[&>td]:bg-[var(--table_row_hover)]
                    "
                    onKeyDown={(e) => { if (e.key === 'Enter') nav(`/results/${r.quizNumber}`) }}
                  >
                    <td className="px-3 py-2 group-hover:bg-[var(--table_row_hover)] active:bg-[var(--table_row_active)] even:bg-[var(--table_tr_e)]">
                      {r.quizNumber}
                    </td>


                    <td className="px-3 py-2 whitespace-nowrap">
                      {new Date(r.createdAt).toLocaleString()}
                    </td>
                    <td className="px-3 py-2">{['全', '登録のみ', '未登録のみ'][r.isRegisteredWords]}</td>
                    <td className="px-3 py-2">{['全て', '含む', '含まない'][r.isIdioms]}</td>
                    <td className="px-3 py-2">{['全て', '含む', '含まない'][r.isSpecialCharacters]}</td>
                    <td className="px-3 py-2">
                      {r.choicesPosIds.map((id) => POS_MAP[id] ?? id).join(', ')}
                    </td>
                    <td className="px-3 py-2 text-right">{r.totalQuestionsCount}</td>
                    <td className="px-3 py-2 text-right">{r.correctCount}</td>
                    <td className="px-3 py-2 text-right">
                      {r.resultCorrectRate.toFixed(1)}%
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>


        {/* 下部ページャ（モバイル用補助） */}
        <div className="mt-4 flex items-center justify-center gap-2 sm:justify-end">
          <Button
            variant="outline"
            disabled={page === 0}
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            className="px-3"
          >
            Prev
          </Button>
          <span className="text-sm opacity-70">
            {page + 1} / {totalPages}
          </span>
          <Button
            variant="outline"
            disabled={page + 1 >= totalPages}
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            className="px-3"
          >
            Next
          </Button>
        </div>
      </Card>
    </div>
  )
}

export default ResultIndex
