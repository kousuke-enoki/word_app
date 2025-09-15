import '@/styles/components/word/WordList.css'

import React, { useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import Pagination from '@/components/common/Pagination'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { getPartOfSpeech } from '@/service/word/GetPartOfSpeech'
import { registerWord } from '@/service/word/RegisterWord'
import type { JapaneseMean, Word, WordInfo } from '@/types/wordTypes'

import PageBottomNav from '../common/PageBottomNav'
import PageTitle from '../common/PageTitle'

const WordList: React.FC = () => {
  const [words, setWords] = useState<Word[]>([])
  const [search, setSearch] = useState('')
  const [sortBy, setSortBy] = useState('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')
  const location = useLocation()
  const [page, setPage] = useState<number>(location.state?.page || 1)
  const [totalPages, setTotalPages] = useState(1)
  const [limit, setLimit] = useState(10)
  const [isInitialized, setIsInitialized] = useState(false)
  const [successMessage, setSuccessMessage] = useState('')

  const getPartOfSpeechName = (id: number) =>
    getPartOfSpeech.find((pos) => pos.id === id)?.name ?? '未定義'

  useEffect(() => {
    if (location.state) {
      setSearch(location.state.search || '')
      setSortBy(location.state.sortBy || 'name')
      setOrder((location.state.order as 'asc' | 'desc') || 'asc')
      setPage(location.state.page || 1)
      setLimit(location.state.limit || 10)
    }
    setIsInitialized(true)
  }, [location.state])

  useEffect(() => {
    if (!isInitialized) return
    const fetchWords = async () => {
      try {
        const { data } = await axiosInstance.get('words', {
          params: { search, sortBy, order, page, limit },
        })
        setWords(data.words)
        setTotalPages(data.totalPages)
      } catch (e) {
        console.error('Failed to fetch words:', e)
      }
    }
    fetchWords()
  }, [search, sortBy, order, page, limit, isInitialized])

  const handleRegister = async (word: Word) => {
    try {
      const updated = await registerWord(word.id, !word.isRegistered)
      setWords((prev) =>
        prev.map((w) =>
          w.id === word.id
            ? {
                ...w,
                isRegistered: updated.isRegistered,
                registrationCount: updated.registrationCount,
              }
            : w,
        ),
      )
      setSuccessMessage(
        `${updated.name} を${updated.isRegistered ? '登録しました' : '登録解除しました'}。`,
      )
      setTimeout(() => setSuccessMessage(''), 3000)
    } catch (e) {
      console.error('Error registering word:', e)
    }
  }

  const Toolbar = (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <div className="flex-1 min-w-[200px]">
        <Input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="単語検索"
        />
      </div>
      <select
        className="rounded-xl border border-[var(--input_bd)] bg-[var(--select)] px-3 py-2 text-[var(--select_c)]"
        value={sortBy}
        onChange={(e) => {
          const v = e.target.value
          if (v === 'register' && sortBy !== 'register') setPage(1)
          setSortBy(v)
        }}
      >
        <option value="name">単語名</option>
        <option value="registrationCount">登録数</option>
        <option value="register">登録</option>
      </select>
      <Button
        variant="outline"
        onClick={() => setOrder(order === 'asc' ? 'desc' : 'asc')}
      >
        {order === 'asc' ? '昇順' : '降順'}
      </Button>
      <Badge>総ページ: {totalPages}</Badge>
    </div>
  )

  // const Pagination = (
  //   <div className="mt-4 flex flex-wrap items-center gap-2">
  //     <select
  //       className="rounded-xl border border-[var(--input_bd)] bg-[var(--select)] px-3 py-2 text-[var(--select_c)]"
  //       value={limit}
  //       onChange={(e) => setLimit(Number(e.target.value))}
  //     >
  //       {[10, 20, 30, 50].map((n) => (
  //         <option key={n} value={n}>
  //           {n}
  //         </option>
  //       ))}
  //     </select>
  //     <Button onClick={() => setPage(1)} disabled={page === 1}>
  //       最初へ
  //     </Button>
  //     <Button onClick={() => setPage(page - 1)} disabled={page === 1}>
  //       前へ
  //     </Button>
  //     <span className="px-2 text-sm opacity-80">
  //       ページ {page} / {totalPages}
  //     </span>
  //     <Button onClick={() => setPage(page + 1)} disabled={page === totalPages}>
  //       次へ
  //     </Button>
  //     <Button
  //       onClick={() => setPage(totalPages)}
  //       disabled={page === totalPages}
  //     >
  //       最後へ
  //     </Button>
  //   </div>
  // )

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <PageTitle title="単語一覧" />
        <Link to="/words/new">
          <Button>新規登録</Button>
        </Link>
      </div>

      {successMessage && (
        <div className="mb-4 rounded-xl border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-4 py-3 text-sm">
          {successMessage}
        </div>
      )}

      <Card className="p-4">
        {Toolbar}
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="bg-[var(--thbc)] text-left">
                {['単語名', '日本語訳', '品詞', '登録数', '登録'].map((th) => (
                  <th
                    key={th}
                    className="border-b border-[var(--thbd)] px-3 py-2 text-[var(--fg)]"
                  >
                    {th}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {words.map((w) => (
                <tr key={w.id} className="even:bg-[var(--table_tr_e)]">
                  <td className="px-3 py-2">
                    <Link
                      to={`/words/${w.id}`}
                      state={{ search, sortBy, order, page, limit }}
                      className="underline"
                    >
                      {w.name}
                    </Link>
                  </td>
                  <td className="px-3 py-2">
                    {w.wordInfos
                      .map((info: WordInfo) =>
                        info.japaneseMeans
                          .map((jm: JapaneseMean) => jm.name)
                          .join(', '),
                      )
                      .join(', ')}
                  </td>
                  <td className="px-3 py-2">
                    {w.wordInfos
                      .map((info: WordInfo) =>
                        getPartOfSpeechName(info.partOfSpeechId),
                      )
                      .join(', ')}
                  </td>
                  <td className="px-3 py-2">{w.registrationCount}</td>
                  <td className="px-3 py-2">
                    <Button
                      className="min-w-[80px]"
                      variant={w.isRegistered ? 'outline' : 'primary'}
                      onClick={() => handleRegister(w)}
                    >
                      {w.isRegistered ? '解除' : '登録'}
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <Pagination
          className="mt-4"
          page={page}
          totalPages={totalPages}
          onPageChange={(p) => setPage(p)}
          pageSize={limit}
          onPageSizeChange={(n) => setLimit(n)}
          pageSizeOptions={[10, 20, 30, 50]}
        />
      </Card>
      <Card className="mt1 p-2">
        <PageBottomNav className="mt-1" showHome inline compact />
      </Card>
    </div>
  )
}

export default WordList
