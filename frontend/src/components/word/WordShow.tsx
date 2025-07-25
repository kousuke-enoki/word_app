import '@/styles/components/word/WordShow.css'

import React, { useEffect, useState } from 'react'
import { useLocation,useNavigate, useParams } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { deleteWord } from '@/service/word/DeleteWord'
import { getPartOfSpeech } from '@/service/word/GetPartOfSpeech'
import { registerWord } from '@/service/word/RegisterWord'
import { saveMemo } from '@/service/word/SaveMemo'
import { JapaneseMean,Word, WordInfo } from '@/types/wordTypes'

const WordShow: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const [word, setWord] = useState<Word | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const [memo, setMemo] = useState<string>('')
  const [successMessage, setSuccessMessage] = useState<string>(
    () => location.state?.successMessage || '',
  )

  useEffect(() => {
    const fetchWord = async () => {
      try {
        const response = await axiosInstance.get(`/words/${id}`)
        setWord(response.data)
        setMemo(response.data.memo || '')
      } catch (error) {
        console.error(error)
        alert('単語情報の取得中にエラーが発生しました。')
      } finally {
        setLoading(false)
      }
    }
    fetchWord()
  }, [id])

  // useEffect(() => {
  //   if (successMessage) {
  //     const timer = setTimeout(() => {
  //       setSuccessMessage('')
  //     }, 3000)
  //     return () => clearTimeout(timer)
  //   }
  // }, [successMessage])

  if (loading) {
    return <p>Loading...</p>
  }

  if (!word) {
    return <p>No word details found.</p>
  }
  if (loading) {
    return <p>Loading...</p>
  }

  if (!word) {
    return <p>No word details found.</p>
  }

  // IDから品詞名を取得するヘルパー関数
  const getPartOfSpeechName = (id: number): string => {
    const partOfSpeech = getPartOfSpeech.find((pos) => pos.id === id)
    return partOfSpeech ? partOfSpeech.name : '未定義'
  }

  const handleRegister = async () => {
    if (!word) return
    try {
      // API呼び出しから新しい登録状態と登録数を取得
      const updatedWord = await registerWord(word.id, !word.isRegistered)

      // 登録状態と登録数を更新
      setWord({
        ...word,
        isRegistered: updatedWord.isRegistered,
        registrationCount: updatedWord.registrationCount,
      })
      if (updatedWord.isRegistered) {
        setSuccessMessage('登録しました。')
      } else {
        setSuccessMessage('登録解除しました。')
      }
      setTimeout(() => setSuccessMessage(''), 3000)
    } catch (error) {
      console.error(error)
      alert('単語の登録中にエラーが発生しました。')
    }
  }

  const handleEdit = async () => {
    window.location.href = '/words/edit/' + word.id
  }

  const handleDelete = async () => {
    if (!word) return
    const confirmDelete = window.confirm('本当にこの単語を削除しますか？')
    if (!confirmDelete) return

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
      }, 1500)
    } catch (error) {
      console.error(error)
      setSuccessMessage('単語の削除に失敗しました。')
      setTimeout(() => setSuccessMessage(''), 3000)
    }
  }

  const handleSaveMemo = async () => {
    if (!word) return
    try {
      await saveMemo(word.id, memo || '')
      setSuccessMessage('メモを保存しました！')
      setTimeout(() => setSuccessMessage(''), 3000)
    } catch (error) {
      console.error(error)
      alert('メモの保存中にエラーが発生しました。')
    }
  }

  const handleMemoChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMemo(e.target.value)
  }

  return (
    <div className="container">
      {successMessage && (
        <div className="success-message">{successMessage}</div>
      )}
      <h1>{word.name}</h1>
      {word.wordInfos.map((info: WordInfo) => (
        <div key={info.id}>
          <p>
            日本語訳:{' '}
            {info.japaneseMeans
              .map((japaneseMean: JapaneseMean) => japaneseMean.name)
              .join(', ')}
          </p>
          <p>品詞: {getPartOfSpeechName(info.partOfSpeechId)}</p>
        </div>
      ))}
      <p>{word.isRegistered ? '登録済み' : '未登録'}</p>
      <p>全ユーザーの登録数: {word.registrationCount}</p>
      <p>単語注意レベル: {word.attentionLevel}</p>
      <p>テスト回数: {word.quizCount}</p>
      <p>チェック回数: {word.correctCount}</p>
      <div>
        <label>
          メモ:
          <textarea
            value={memo}
            onChange={handleMemoChange}
            rows={6} // 高さ5行
            cols={35} // 横幅30文字
            maxLength={200} // 200文字制限
          />
        </label>
        <button className="save-button" onClick={handleSaveMemo}>
          保存する
        </button>
      </div>
      {/* {successMessage && <div className="success-popup">{successMessage}</div>} */}
      <div>
        <button
          className={`register-button ${word.isRegistered ? 'registered' : ''}`}
          onClick={handleRegister}
        >
          {word.isRegistered ? '登録解除' : '登録する'}
        </button>
      </div>
      <div>
        <button className="delete-button" onClick={handleEdit}>
          編集する
        </button>
        <button className="delete-button" onClick={handleDelete}>
          削除する
        </button>
      </div>
      <button
        className="back-button"
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
      </button>
    </div>
  )
}

export default WordShow
