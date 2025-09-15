import '@/styles/components/word/WordNew.css'

import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import {
  getPartOfSpeech,
  PartOfSpeechOption,
} from '@/service/word/GetPartOfSpeech'

import PageBottomNav from '../common/PageBottomNav'
import PageTitle from '../common/PageTitle'

export type WordForNew = { name: string; wordInfos: WordInfoForNew[] }
export type WordInfoForNew = {
  partOfSpeechId: number
  japaneseMeans: JapaneseMeansForNew[]
}
export type JapaneseMeansForNew = { name: string }
export type ValidationErrors = {
  name?: string
  wordInfos?: Array<{ partOfSpeech?: string; japaneseMeans?: string[] }>
}

const WordNew: React.FC = () => {
  const [word, setWord] = useState<WordForNew>({
    name: '',
    wordInfos: [{ partOfSpeechId: 0, japaneseMeans: [{ name: '' }] }],
  })
  const [validationErrors, setValidationErrors] = useState<ValidationErrors>({})
  const [errorMessage, setErrorMessage] = useState('')
  const navigate = useNavigate()

  const MAX_PART_OF_SPEECH = 10
  const MAX_JAPANESE_MEANS = 10
  const wordNameRegex = /^[A-Za-z]+$/
  // eslint-disable-next-line no-control-regex
  const japaneseMeanRegex = /^[^-~｡-ﾟ~]+$/

  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    if (value === '' || wordNameRegex.test(value))
      setWord({ ...word, name: value })
    else alert('単語名は半角アルファベットのみ入力できます。')
  }

  const handlePartOfSpeechChange = (index: number, value: string) => {
    const updated = [...word.wordInfos]
    updated[index].partOfSpeechId = parseInt(value, 10)
    setWord({ ...word, wordInfos: updated })
  }

  const handleJapaneseMeanChange = (wi: number, mi: number, value: string) => {
    if (value === '' || japaneseMeanRegex.test(value)) {
      const updated = [...word.wordInfos]
      updated[wi].japaneseMeans[mi].name = value
      setWord({ ...word, wordInfos: updated })
    } else
      alert(
        '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      )
  }

  const addPartOfSpeech = () =>
    setWord({
      ...word,
      wordInfos: [
        ...word.wordInfos,
        { partOfSpeechId: 0, japaneseMeans: [{ name: '' }] },
      ],
    })
  const removePartOfSpeech = (index: number) =>
    setWord({
      ...word,
      wordInfos: word.wordInfos.filter((_, i) => i !== index),
    })
  const addJapaneseMean = (index: number) => {
    const updated = [...word.wordInfos]
    updated[index].japaneseMeans.push({ name: '' })
    setWord({ ...word, wordInfos: updated })
  }
  const removeJapaneseMean = (wi: number, mi: number) => {
    const updated = [...word.wordInfos]
    updated[wi].japaneseMeans = updated[wi].japaneseMeans.filter(
      (_, i) => i !== mi,
    )
    setWord({ ...word, wordInfos: updated })
  }

  const getAvailablePartOfSpeechOptions = (
    currentIndex: number,
  ): PartOfSpeechOption[] => {
    const selected = word.wordInfos
      .filter((_, i) => i !== currentIndex)
      .map((i) => i.partOfSpeechId)
    return getPartOfSpeech.filter((o) => !selected.includes(o.id))
  }

  const validateWord = (target: WordForNew) => {
    const newErrors: ValidationErrors = {}
    if (!wordNameRegex.test(target.name))
      newErrors.name = '単語名は半角アルファベットのみ入力できます。'
    const infoErrors = target.wordInfos.map((info) => {
      const e: { partOfSpeech?: string; japaneseMeans?: string[] } = {}
      if (info.partOfSpeechId === 0) e.partOfSpeech = '品詞を選択してください。'
      const means = info.japaneseMeans.map((m) =>
        japaneseMeanRegex.test(m.name)
          ? ''
          : '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      )
      if (means.some((x) => x !== '')) e.japaneseMeans = means
      return e
    })
    if (
      infoErrors.some(
        (x) => x.partOfSpeech || (x.japaneseMeans && x.japaneseMeans.length),
      )
    )
      newErrors.wordInfos = infoErrors
    return newErrors
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const errors = validateWord(word)
    setValidationErrors(errors)
    if (Object.keys(errors).length) return
    try {
      const { data } = await axiosInstance.post('/words/new', word)
      navigate(`/words/${data.id}`, {
        state: { successMessage: `${data.name}が正常に登録されました！` },
      })
    } catch {
      setErrorMessage('単語の登録中にエラーが発生しました。')
    }
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <PageTitle title="単語登録" />
        <Badge>✍️ New</Badge>
      </div>

      {errorMessage && (
        <div className="mb-4 rounded-xl border-l-4 border-red-500 bg-[var(--container_bg)] px-4 py-3 text-sm text-red-600">
          {errorMessage}
        </div>
      )}

      <Card className="p-6">
        <form
          onSubmit={handleSubmit}
          className="space-y-6"
          aria-label="word-create-form"
        >
          <div>
            <label className="mb-1 block text-sm font-medium">単語名</label>
            <Input
              value={word.name}
              onChange={handleWordNameChange}
              placeholder="example"
              required
            />
            {validationErrors.name && (
              <p className="mt-1 text-sm text-red-600">
                {validationErrors.name}
              </p>
            )}
          </div>
          {word.wordInfos.map((wi, wiIndex) => {
            const infoErr = validationErrors.wordInfos?.[wiIndex]
            return (
              <Card key={wiIndex} className="p-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label className="mb-1 block text-sm font-medium">
                      品詞
                    </label>
                    <select
                      className="w-full rounded-xl border border-[var(--input_bd)] bg-[var(--select)] px-3 py-2 text-[var(--select_c)]"
                      value={wi.partOfSpeechId}
                      onChange={(e) =>
                        handlePartOfSpeechChange(wiIndex, e.target.value)
                      }
                      required
                    >
                      <option value={0}>選択してください</option>
                      {getAvailablePartOfSpeechOptions(wiIndex).map((o) => (
                        <option key={o.id} value={o.id}>
                          {o.name}
                        </option>
                      ))}
                    </select>
                    {infoErr?.partOfSpeech && (
                      <p className="mt-1 text-sm text-red-600">
                        {infoErr.partOfSpeech}
                      </p>
                    )}
                  </div>

                  <div className="flex items-start justify-end gap-2">
                    {word.wordInfos.length > 1 && (
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => removePartOfSpeech(wiIndex)}
                      >
                        品詞を削除
                      </Button>
                    )}
                  </div>
                </div>

                <div className="mt-4 space-y-3">
                  {wi.japaneseMeans.map((m, mi) => (
                    <div
                      key={mi}
                      className="grid gap-2 sm:grid-cols-[1fr_auto] sm:items-center"
                    >
                      <div>
                        <label className="mb-1 block text-sm font-medium">
                          日本語訳
                        </label>
                        <Input
                          value={m.name}
                          onChange={(e) =>
                            handleJapaneseMeanChange(
                              wiIndex,
                              mi,
                              e.target.value,
                            )
                          }
                          placeholder="意味"
                          required
                        />
                        {infoErr?.japaneseMeans?.[mi] && (
                          <p className="mt-1 text-sm text-red-600">
                            {infoErr.japaneseMeans[mi]}
                          </p>
                        )}
                      </div>
                      {wi.japaneseMeans.length > 1 && (
                        <Button
                          type="button"
                          variant="outline"
                          onClick={() => removeJapaneseMean(wiIndex, mi)}
                        >
                          削除
                        </Button>
                      )}
                    </div>
                  ))}

                  {wi.japaneseMeans.length < MAX_JAPANESE_MEANS && (
                    <Button
                      type="button"
                      variant="ghost"
                      onClick={() => addJapaneseMean(wiIndex)}
                    >
                      日本語訳を追加
                    </Button>
                  )}
                </div>
              </Card>
            )
          })}
          {word.wordInfos.length < MAX_PART_OF_SPEECH && (
            <Button type="button" variant="ghost" onClick={addPartOfSpeech}>
              品詞を追加
            </Button>
          )}
        </form>
      </Card>
      <Card className="mt1 p-2">
        <PageBottomNav
          className="mt-1"
          actions={[{ label: '単語一覧', to: '/words' }]}
          showHome
          inline
          compact
        />
      </Card>
    </div>
  )
}

export default WordNew
