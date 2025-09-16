import '@/styles/components/word/WordShow.css'

import React, { useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { RegisterToggle } from '@/components/common/RegisterToggle'
import { Badge, Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { deleteWord } from '@/service/word/DeleteWord'
import { getPartOfSpeech } from '@/service/word/GetPartOfSpeech'
import { registerWord } from '@/service/word/RegisterWord'
import { saveMemo } from '@/service/word/SaveMemo'
import type { JapaneseMean, Word, WordInfo } from '@/types/wordTypes'

import PageBottomNav from '../common/PageBottomNav'

const WordShow: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const [word, setWord] = useState<Word | null>(null)
  const [loading, setLoading] = useState(true)
  const [memo, setMemo] = useState('')
  const [successMessage, setSuccessMessage] = useState<string>(
    () => location.state?.successMessage || '',
  )

  useEffect(() => {
    const fetchWord = async () => {
      try {
        const { data } = await axiosInstance.get(`/words/${id}`)
        setWord(data)
        setMemo(data.memo || '')
      } catch (e) {
        console.error(e)
        alert('単語情報の取得中にエラーが発生しました。')
      } finally {
        setLoading(false)
      }
    }
    fetchWord()
  }, [id])

  if (loading) return <p>Loading...</p>
  if (!word) return <p>No word details found.</p>

  const getPartOfSpeechName = (id: number) =>
    getPartOfSpeech.find((pos) => pos.id === id)?.name ?? '未定義'

  const handleRegister = async () => {
    if (!word) return
    try {
      const updated = await registerWord(word.id, !word.isRegistered)
      setWord({
        ...word,
        isRegistered: updated.isRegistered,
        registrationCount: updated.registrationCount,
      })
      setSuccessMessage(
        updated.isRegistered ? '登録しました。' : '登録解除しました。',
      )
      setTimeout(() => setSuccessMessage(''), 2500)
    } catch (e) {
      console.error(e)
      alert('単語の登録中にエラーが発生しました。')
    }
  }

  const handleEdit = () => navigate(`/words/edit/${word.id}`)

  const handleDelete = async () => {
    if (!word) return
    if (!window.confirm('本当にこの単語を削除しますか？')) return
    try {
      await deleteWord(word.id)
      setSuccessMessage('単語を削除しました。')
      setTimeout(() => {
        navigate('/words', {
          state: {
            search: location.state?.search || '',
            sortBy: location.state?.sortBy || 'name',
            order: location.state?.order || 'asc',
            page: location.state?.page || 1,
            limit: location.state?.limit || 10,
          },
        })
      }, 1200)
    } catch (e) {
      console.error(e)
      setSuccessMessage('単語の削除に失敗しました。')
      setTimeout(() => setSuccessMessage(''), 3000)
    }
  }

  const handleSaveMemo = async () => {
    if (!word) return
    try {
      await saveMemo(word.id, memo || '')
      setSuccessMessage('メモを保存しました！')
      setTimeout(() => setSuccessMessage(''), 2500)
    } catch (e) {
      console.error(e)
      alert('メモの保存中にエラーが発生しました。')
    }
  }

  return (
    <div>
      {successMessage && (
        <div className="mb-4 rounded-xl border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-4 py-3 text-sm">
          {successMessage}
        </div>
      )}
      <div
        className="
            mb-4 grid items-center gap-2
            grid-cols-1
            sm:grid-cols-[1fr_auto_1fr]   /* ← 左右1fr, 中央auto：常に中央寄せ */
          "
      >
        {/* タイトル：中央カラムに配置 */}
        <h1
          className="
              text-2xl font-bold text-[var(--h1_fg)] text-center
              sm:col-start-2
            "
        >
          {word.name}
        </h1>
        <div className="justify-self-center sm:justify-self-end">
          <RegisterToggle
            isRegistered={word.isRegistered}
            onToggle={handleRegister}
            variant="split"
            widthClass="w-36 sm:w-40"
          />
        </div>
      </div>
      <Card className="mb-4 p-5">
        {word.wordInfos.map((info: WordInfo) => (
          <div key={info.id} className="mb-3">
            <p>
              <span className="opacity-70">日本語訳:</span>{' '}
              {info.japaneseMeans.map((jm: JapaneseMean) => jm.name).join(', ')}
            </p>
            <p className="opacity-80">
              品詞: {getPartOfSpeechName(info.partOfSpeechId)}
            </p>
          </div>
        ))}
        <div className="mt-2 flex flex-wrap gap-2 text-sm opacity-80">
          <Badge>全登録数: {word.registrationCount}</Badge>
          <Badge>注意レベル: {word.attentionLevel}</Badge>
          <Badge>テスト回数: {word.quizCount}</Badge>
          <Badge>チェック回数: {word.correctCount}</Badge>
        </div>{' '}
        <div className="flex flex-wrap items-center gap-3">
          <Button
            variant="outline"
            onClick={() =>
              navigate('/words', {
                state: {
                  search: location.state?.search || '',
                  sortBy: location.state?.sortBy || 'name',
                  order: location.state?.order || 'asc',
                  page: location.state?.page || 1,
                  limit: location.state?.limit || 10,
                },
              })
            }
          >
            一覧に戻る
          </Button>
          <Button onClick={handleEdit}>編集する</Button>
          <Button variant="outline" onClick={handleDelete}>
            削除する
          </Button>
        </div>{' '}
      </Card>
      <Card className="mb-6 p-5">
        <div className="mb-2 text-sm font-medium">メモ（200文字まで）</div>
        <textarea
          className="w-full rounded-xl border border-[var(--textarea_bd, var(--input_bd))] bg-[var(--textarea_bg)] p-3 text-[var(--textarea_c)] outline-none focus:ring-2 ring-[var(--button_bg)]"
          value={memo}
          onChange={(e) => setMemo(e.target.value)}
          rows={6}
          maxLength={200}
        />
        <div className="mt-3">
          <Button onClick={handleSaveMemo}>保存する</Button>
        </div>
      </Card>
      <Card className="mt1 p-2">
        <PageBottomNav className="mt-1" showHome inline compact />
      </Card>
    </div>
  )
}

export default WordShow
