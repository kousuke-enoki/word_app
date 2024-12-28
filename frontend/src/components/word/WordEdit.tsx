import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import axiosInstance from '../../axiosConfig'
import {
  getPartOfSpeech,
  PartOfSpeechOption,
} from '../../service/word/GetPartOfSpeech'

export type WordForUpdate = {
  id: number
  name: string
  wordInfos: WordInfoForUpdate[]
}

export type WordInfoForUpdate = {
  id: number
  partOfSpeechId: number
  japaneseMeans: JapaneseMeansForUpdate[]
}

export type JapaneseMeansForUpdate = {
  id: number
  name: string
}

const WordEdit: React.FC = () => {
  const { id } = useParams()
  const [word, setWord] = useState<WordForUpdate | null>(null)
  const [successMessage, setSuccessMessage] = useState<string>('')

  // const MAX_PART_OF_SPEECH = 10
  // const MAX_JAPANESE_MEANS = 10

  const wordNameRegex = /^[A-Za-z]+$/
  // eslint-disable-next-line no-control-regex
  const japaneseMeanRegex = /^[^\x01-\x7E\uFF61-\uFF9F~]+$/

  // 初期データを取得
  useEffect(() => {
    const fetchWord = async () => {
      try {
        const response = await axiosInstance.get(`/words/${id}`)
        setWord(response.data)
      } catch (error) {
        alert('単語情報の取得中にエラーが発生しました。')
      }
    }
    fetchWord()
  }, [id])

  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    if (value === '' || wordNameRegex.test(value)) {
      setWord((prevWord) => (prevWord ? { ...prevWord, name: value } : null))
    } else {
      alert('単語名は半角アルファベットのみ入力できます。')
    }
  }

  const handlePartOfSpeechChange = (index: number, value: string) => {
    if (word) {
      const updatedWordInfos = [...word.wordInfos]
      updatedWordInfos[index].partOfSpeechId = parseInt(value, 10)
      setWord({ ...word, wordInfos: updatedWordInfos })
    }
  }

  const handleJapaneseMeanChange = (
    wordInfoIndex: number,
    japaneseMeanIndex: number,
    value: string,
  ) => {
    if (value === '' || japaneseMeanRegex.test(value)) {
      if (word) {
        const updatedWordInfos = [...word.wordInfos]
        updatedWordInfos[wordInfoIndex].japaneseMeans[japaneseMeanIndex].name =
          value
        setWord({ ...word, wordInfos: updatedWordInfos })
      }
    } else {
      alert(
        '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      )
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await axiosInstance.put(`/words/${id}`, word)
      setSuccessMessage(response.data.name + 'が正常に更新されました！')
      setTimeout(() => setSuccessMessage(''), 3000)
      setTimeout(() => {
        window.location.href = '/words/' + id
      }, 1500)
    } catch (error) {
      alert('単語情報の更新中にエラーが発生しました。')
    }
  }

  const getAvailablePartOfSpeechOptions = (
    currentIndex: number,
  ): PartOfSpeechOption[] => {
    if (!word) return []
    const selectedIds = word.wordInfos
      .filter((_, index) => index !== currentIndex)
      .map((info) => info.partOfSpeechId)

    return getPartOfSpeech.filter((option) => !selectedIds.includes(option.id))
  }

  return (
    <div className="word-update-container">
      <h1>単語更新フォーム</h1>
      {word ? (
        <form className="word-update-form" onSubmit={handleSubmit}>
          {successMessage && (
            <div className="success-popup">{successMessage}</div>
          )}
          <div>
            <label>
              単語名:
              <input
                type="text"
                value={word.name}
                onChange={handleWordNameChange}
                required
              />
            </label>
          </div>

          {word.wordInfos.map((wordInfo, wordInfoIndex) => (
            <div key={wordInfo.id} className="word-info-section">
              <div>
                <label>
                  品詞:
                  <select
                    value={wordInfo.partOfSpeechId}
                    onChange={(e) =>
                      handlePartOfSpeechChange(wordInfoIndex, e.target.value)
                    }
                    required
                  >
                    <option value={0}>選択してください</option>
                    {getAvailablePartOfSpeechOptions(wordInfoIndex).map(
                      (option) => (
                        <option key={option.id} value={option.id}>
                          {option.name}
                        </option>
                      ),
                    )}
                  </select>
                </label>
              </div>

              {wordInfo.japaneseMeans.map((mean, meanIndex) => (
                <div key={mean.id} className="japanese-mean-section">
                  <label>
                    日本語訳:
                    <input
                      type="text"
                      value={mean.name}
                      onChange={(e) =>
                        handleJapaneseMeanChange(
                          wordInfoIndex,
                          meanIndex,
                          e.target.value,
                        )
                      }
                      required
                    />
                  </label>
                </div>
              ))}
            </div>
          ))}

          <div className="submit-button">
            <button type="submit">単語を更新</button>
          </div>
          <div>
            <button
              className="back-button"
              onClick={() => (window.location.href = '/words/' + word.id)}
            >
              単語詳細に戻る
            </button>
          </div>
        </form>
      ) : (
        <p>データを読み込み中...</p>
      )}
    </div>
  )
}

export default WordEdit
