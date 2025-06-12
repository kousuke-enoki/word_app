import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import axiosInstance from '@/axiosConfig'
import {
  getPartOfSpeech,
  PartOfSpeechOption,
} from '@/service/word/GetPartOfSpeech'
import '@/styles/components/word/WordNew.css'

export type WordForNew = {
  name: string
  wordInfos: WordInfoForNew[]
}

export type WordInfoForNew = {
  partOfSpeechId: number
  japaneseMeans: japaneseMeansForNew[]
}

export type japaneseMeansForNew = {
  name: string
}

// バリデーションエラーをフィールドごとに管理するための型
export type ValidationErrors = {
  name?: string
  wordInfos?: Array<{
    partOfSpeech?: string
    japaneseMeans?: string[]
  }>
}

const WordNew: React.FC = () => {
  const [word, setWord] = useState<WordForNew>({
    name: '',
    wordInfos: [
      {
        partOfSpeechId: 0,
        japaneseMeans: [{ name: '' }],
      },
    ],
  })
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>({})
  const [errorMessage, setErrorMessage] = useState('')
  const navigate = useNavigate()

  const MAX_PART_OF_SPEECH = 10
  const MAX_JAPANESE_MEANS = 10

  // バリデーション用の正規表現
  const wordNameRegex = /^[A-Za-z]+$/ // 半角アルファベットのみ
  // eslint-disable-next-line no-control-regex
  const japaneseMeanRegex = /^[^\x01-\x7E\uFF61-\uFF9F~]+$/ // 日本語 (ひらがな、カタカナ、漢字) と記号「~」のみ

  // 単語名の変更ハンドラー (修正版)
  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    if (value === '' || wordNameRegex.test(value)) {
      setWord({ ...word, name: value })
    } else {
      alert('単語名は半角アルファベットのみ入力できます。')
    }
  }

  // 品詞の変更ハンドラー
  const handlePartOfSpeechChange = (index: number, value: string) => {
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[index].partOfSpeechId = parseInt(value, 10)
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 日本語訳の変更ハンドラー (修正版)
  const handleJapaneseMeanChange = (
    wordInfoIndex: number,
    japaneseMeanIndex: number,
    value: string,
  ) => {
    if (value === '' || japaneseMeanRegex.test(value)) {
      const updatedWordInfos = [...word.wordInfos]
      updatedWordInfos[wordInfoIndex].japaneseMeans[japaneseMeanIndex].name =
        value
      setWord({ ...word, wordInfos: updatedWordInfos })
    } else {
      alert(
        '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      )
    }
  }

  // 品詞を追加
  const addPartOfSpeech = () => {
    setWord({
      ...word,
      wordInfos: [
        ...word.wordInfos,
        { partOfSpeechId: 0, japaneseMeans: [{ name: '' }] },
      ],
    })
  }

  // 日本語訳を追加
  const addJapaneseMean = (index: number) => {
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[index].japaneseMeans.push({ name: '' })
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 品詞を削除
  const removePartOfSpeech = (index: number) => {
    const updatedWordInfos = word.wordInfos.filter((_, i) => i !== index)
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 日本語訳を削除
  const removeJapaneseMean = (
    wordInfoIndex: number,
    japaneseMeanIndex: number,
  ) => {
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[wordInfoIndex].japaneseMeans = updatedWordInfos[
      wordInfoIndex
    ].japaneseMeans.filter((_, i) => i !== japaneseMeanIndex)
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 使用可能な品詞の選択肢を取得
  const getAvailablePartOfSpeechOptions = (
    currentIndex: number,
  ): PartOfSpeechOption[] => {
    const selectedIds = word.wordInfos
      .filter((_, index) => index !== currentIndex) // 現在のフォーム以外を対象
      .map((info) => info.partOfSpeechId)

    return getPartOfSpeech.filter((option) => !selectedIds.includes(option.id))
  }

  // ★ バリデーションロジック（フィールドごとにエラーメッセージをセット）
  const validateWord = (targetWord: WordForNew) => {
    const newErrors: ValidationErrors = {}

    // 単語名
    if (!wordNameRegex.test(targetWord.name)) {
      newErrors.name = '単語名は半角アルファベットのみ入力できます。'
    }

    // wordInfos
    const wordInfoErrors = targetWord.wordInfos.map((info) => {
      const infoError: {
        partOfSpeech?: string
        japaneseMeans?: string[]
      } = {}
      if (info.partOfSpeechId === 0) {
        infoError.partOfSpeech = '品詞を選択してください。'
      }
      // 日本語訳
      const meansErrors = info.japaneseMeans.map((mean) => {
        if (!japaneseMeanRegex.test(mean.name)) {
          return '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。'
        }
        return '' // 問題なし
      })
      // 空文字列以外があればエラー
      if (meansErrors.some((err) => err !== '')) {
        infoError.japaneseMeans = meansErrors
      }
      return infoError
    })

    // wordInfoErrors のいずれかにエラーがある場合のみ格納
    if (
      wordInfoErrors.some(
        (infoError) =>
          infoError.partOfSpeech ||
          (infoError.japaneseMeans && infoError.japaneseMeans.length > 0),
      )
    ) {
      newErrors.wordInfos = wordInfoErrors
    }

    return newErrors
  }

  const NewWordMutation = useMutation<
    { name: string; id: number },
    unknown,
    WordForNew
  >({
    // ミューテーション関数
    mutationFn: async (newWord: WordForNew) => {
      const response = await axiosInstance.post(`/words/new`, newWord)
      return response.data
    },
    onSuccess: (data) => {
      const newName = data.name
      // メッセージを表示せず、遷移先へ渡す
      navigate(`/words/${data.id}`, {
        state: {
          successMessage: `${newName}が正常に登録されました！`,
        },
      })
    },
    onError: () => {
      setErrorMessage('単語の登録中にエラーが発生しました。')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!word) return

    // バリデーション
    const errors = validateWord(word)
    setValidationErrors(errors)

    // errors オブジェクトに何かしらエラーがあれば送信中断
    if (Object.keys(errors).length > 0) {
      return
    }

    NewWordMutation.mutate(word)
  }

  return (
    <div className="word-create-container">
      <h1>単語登録フォーム</h1>
      {/* エラーメッセージ */}
      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      <form 
        aria-label="word-create-form"
        className="word-create-form"
        onSubmit={handleSubmit}
      >
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
          {validationErrors.name && (
            <p style={{ color: 'red' }}>{validationErrors.name}</p>
          )}
        </div>

        {word.wordInfos.map((wordInfo, wordInfoIndex) => {
          const infoError = validationErrors.wordInfos?.[wordInfoIndex]
          return (
            <div key={wordInfoIndex} className="word-info-section">
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
                {infoError?.partOfSpeech && (
                  <p style={{ color: 'red' }}>{infoError.partOfSpeech}</p>
                )}
                {word.wordInfos.length > 1 && (
                  <button 
                    className="part-of-speech-section"
                    type="button"
                    onClick={() => removePartOfSpeech(wordInfoIndex)}
                  >
                    品詞を削除
                  </button>
                )}
              </div>

              {wordInfo.japaneseMeans.map((mean, meanIndex) => (
                <div key={meanIndex} className="japanese-mean-section">
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
                  {infoError?.japaneseMeans?.[meanIndex] && (
                    <p style={{ color: 'red' }}>
                      {infoError.japaneseMeans[meanIndex]}
                    </p>
                  )}
                  {wordInfo.japaneseMeans.length > 1 && (
                    <button
                      type="button"
                      onClick={() =>
                        removeJapaneseMean(wordInfoIndex, meanIndex)
                      }
                    >
                      削除
                    </button>
                  )}
                </div>
              ))}
              {wordInfo.japaneseMeans.length < MAX_JAPANESE_MEANS && (
                <button
                  type="button"
                  onClick={() => addJapaneseMean(wordInfoIndex)}
                >
                  日本語訳を追加
                </button>
              )}
            </div>
          )
        })}
        {word.wordInfos.length < MAX_PART_OF_SPEECH && (
          <button type="button" onClick={addPartOfSpeech}>
            品詞を追加
          </button>
        )}

        <div className="submit-button">
          <button type="submit">単語を登録</button>
        </div>
        <div>
          <button className="back-button" onClick={() => navigate('/', {})}>
            mypageに戻る
          </button>
        </div>
      </form>
    </div>
  )
}

export default WordNew
