import { useMutation,useQuery } from '@tanstack/react-query'
import React, { useEffect, useState } from 'react'
import { useNavigate,useParams } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import {
  getPartOfSpeech,
  PartOfSpeechOption,
} from '@/service/word/GetPartOfSpeech'

// 型定義
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

// バリデーションエラーをフィールドごとに管理するための型
type ValidationErrors = {
  name?: string
  wordInfos?: Array<{
    partOfSpeech?: string
    japaneseMeans?: string[]
  }>
}

const isTestEnv = process.env.NODE_ENV === 'test'

const WordEdit: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()

  // React Queryでデータ取得
  // --------------------------------
  const {
    data: fetchedWord,
    isLoading,
    isError,
    refetch, // 再取得を行う関数
  } = useQuery<WordForUpdate | null>({
    queryKey: ['word', id],
    queryFn: async () => {
      const res = await axiosInstance.get(`/words/${id}`)
      return res.data
    },
    retry: isTestEnv ? false : 3,
    enabled: Boolean(id),
  })

  // フォームで編集するためのローカルステート
  const [word, setWord] = useState<WordForUpdate | null>(null)

  // バリデーションエラー & 成功/失敗メッセージ
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>({})
  const [errorMessage, setErrorMessage] = useState('')

  // 取得後にローカルステートへコピー
  useEffect(() => {
    if (fetchedWord) {
      setWord(fetchedWord)
    }
  }, [fetchedWord])

  // バリデーション用正規表現
  const wordNameRegex = /^[A-Za-z]+$/ // 半角アルファベットのみ
  // eslint-disable-next-line no-control-regex
  const japaneseMeanRegex = /^[^\x01-\x7E\uFF61-\uFF9F~]+$/

  // ★ バリデーションロジック（フィールドごとにエラーメッセージをセット）
  const validateWord = (targetWord: WordForUpdate) => {
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

  // フィールド変更時のハンドラ
  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!word) return
    setWord({ ...word, name: e.target.value })
  }

  const handlePartOfSpeechChange = (index: number, value: string) => {
    if (!word) return
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[index].partOfSpeechId = parseInt(value, 10)
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  const handleJapaneseMeanChange = (
    wordInfoIndex: number,
    meanIndex: number,
    value: string,
  ) => {
    if (!word) return
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[wordInfoIndex].japaneseMeans[meanIndex].name = value
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 更新APIをreact-queryのMutationで管理
  // ---------------------------------------
  const updateWordMutation = useMutation<
    { name: string },
    unknown,
    WordForUpdate
  >({
    // ミューテーション関数
    mutationFn: async (updatedWord: WordForUpdate) => {
      const response = await axiosInstance.put(`/words/${id}`, updatedWord)
      return response.data
    },
    onSuccess: (data) => {
      const newName = data.name
      // メッセージを表示せず、遷移先へ渡す
      navigate(`/words/${id}`, {
        state: {
          successMessage: `${newName}が正常に更新されました！`,
        },
      })
    },
    onError: () => {
      setErrorMessage('単語情報の更新中にエラーが発生しました。')
    },
  })

  // フォーム送信
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

    updateWordMutation.mutate(word)
  }

  // 他の品詞と重複しないようにフィルタ
  const getAvailablePartOfSpeechOptions = (
    currentIndex: number,
  ): PartOfSpeechOption[] => {
    if (!word) return []
    const selectedIds = word.wordInfos
      .filter((_, i) => i !== currentIndex)
      .map((info) => info.partOfSpeechId)

    return getPartOfSpeech.filter((option) => !selectedIds.includes(option.id))
  }

  // --- UI出力 ---
  // ---------------------------------------
  // ローディング中
  if (isLoading) {
    return <p>読み込み中...</p>
  }

  // 取得エラー時
  if (isError) {
    return (
      <div style={{ color: 'red' }}>
        <p>単語情報の取得中にエラーが発生しました。</p>
        <button onClick={() => refetch()}>再取得</button>
      </div>
    )
  }

  // fetchedWordがnullのケース
  if (!fetchedWord) {
    return <p>データが存在しません。</p>
  }

  // wordがnullも同義
  if (!word) {
    return <p>データが存在しません。</p>
  }

  return (
    <div className="word-update-container">
      <h1>単語更新フォーム</h1>
      {/* エラーメッセージ */}
      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      <form className="word-update-form" onSubmit={handleSubmit}>
        {/* 単語名: フィールドごとのエラーメッセージ */}
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

        {/* wordInfos */}
        {word.wordInfos.map((wordInfo, wordInfoIndex) => {
          const infoError = validationErrors.wordInfos?.[wordInfoIndex]
          return (
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
                {infoError?.partOfSpeech && (
                  <p style={{ color: 'red' }}>{infoError.partOfSpeech}</p>
                )}
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
                  {infoError?.japaneseMeans?.[meanIndex] && (
                    <p style={{ color: 'red' }}>
                      {infoError.japaneseMeans[meanIndex]}
                    </p>
                  )}
                </div>
              ))}
            </div>
          )
        })}

        {/* ボタン */}
        <div className="submit-button">
          <button type="submit" disabled={isLoading}>
            単語を更新
          </button>
        </div>
        <div>
          <button
            type="button"
            className="back-button"
            onClick={() => navigate(`/words/${word.id}`)}
          >
            単語詳細に戻る
          </button>
        </div>
      </form>
    </div>
  )
}

export default WordEdit
