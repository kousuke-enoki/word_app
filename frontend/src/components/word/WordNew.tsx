import React, { useState } from 'react'
import axiosInstance from '../../axiosConfig'
import { getPartOfSpeech } from '../../service/word/GetPartOfSpeech'

// 使用しないプロパティを除外
export type WordForNew = {
  name: string
  wordInfos: WordInfoForNew[]
}

export type WordInfoForNew = {
  partOfSpeechId: number
  japaneseMeans: japaneseMeansForNew[]
}

export type partOfSpeechForNew = {
  name: string
}

export type japaneseMeansForNew = {
  name: string
}
export type PartOfSpeechOption = {
  id: number
  name: string
}

const WordCreateForm: React.FC = () => {
  const [word, setWord] = useState<WordForNew>({
    name: '',
    wordInfos: [
      {
        partOfSpeechId: 0,
        japaneseMeans: [{ name: '' }],
      },
    ],
  })

  // 単語名の変更ハンドラー
  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setWord({ ...word, name: e.target.value })
  }

  // 品詞の変更ハンドラー
  const handlePartOfSpeechChange = (index: number, value: string) => {
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[index].partOfSpeechId = parseInt(value, 10)
    setWord({ ...word, wordInfos: updatedWordInfos })
  }

  // 日本語訳の変更ハンドラー
  const handleJapaneseMeanChange = (
    wordInfoIndex: number,
    japaneseMeanIndex: number,
    value: string,
  ) => {
    const updatedWordInfos = [...word.wordInfos]
    updatedWordInfos[wordInfoIndex].japaneseMeans[japaneseMeanIndex].name =
      value
    setWord({ ...word, wordInfos: updatedWordInfos })
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

  // フォーム送信ハンドラー
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await axiosInstance.post('/words/new', word)
      console.log('Word created successfully:', response.data.name)
      alert(response.data.name + 'が正常に登録されました！')
    } catch (error) {
      console.error('Error creating word:', error)
      alert('単語の登録中にエラーが発生しました。')
    }
  }

  return (
    <div>
      <h1>単語登録フォーム</h1>
      <form onSubmit={handleSubmit}>
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
          <div key={wordInfoIndex}>
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
                  {getPartOfSpeech.map((option) => (
                    <option key={option.id} value={option.id}>
                      {option.name}
                    </option>
                  ))}
                </select>
              </label>
              <button
                type="button"
                onClick={() => removePartOfSpeech(wordInfoIndex)}
              >
                品詞を削除
              </button>
            </div>

            {wordInfo.japaneseMeans.map((mean, meanIndex) => (
              <div key={meanIndex}>
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
                <button
                  type="button"
                  onClick={() => removeJapaneseMean(wordInfoIndex, meanIndex)}
                >
                  日本語訳を削除
                </button>
              </div>
            ))}

            <button
              type="button"
              onClick={() => addJapaneseMean(wordInfoIndex)}
            >
              日本語訳を追加
            </button>
          </div>
        ))}

        <button type="button" onClick={addPartOfSpeech}>
          品詞を追加
        </button>

        <div>
          <button type="submit">単語を登録</button>
        </div>
      </form>
    </div>
  )
}

export default WordCreateForm
