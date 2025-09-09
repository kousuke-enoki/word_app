import '@/styles/components/word/WordList.css'

import { useMutation, useQuery } from '@tanstack/react-query'
import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import {
  getPartOfSpeech,
  PartOfSpeechOption,
} from '@/service/word/GetPartOfSpeech'

export type WordForUpdate = {
  id: number
  name: string
  wordInfos: WordInfoForUpdate[]
}
type WordInfoForUpdate = {
  id: number
  partOfSpeechId: number
  japaneseMeans: { id: number; name: string }[]
}

type ValidationErrorsU = {
  name?: string
  wordInfos?: Array<{ partOfSpeech?: string; japaneseMeans?: string[] }>
}

const WordEdit: React.FC = () => {
  const { id } = useParams()
  const navigate = useNavigate()

  const {
    data: fetchedWord,
    isLoading,
    isError,
    refetch,
  } = useQuery<WordForUpdate | null>({
    queryKey: ['word', id],
    queryFn: async () => (await axiosInstance.get(`/words/${id}`)).data,
    enabled: Boolean(id),
  })

  const [word, setWord] = useState<WordForUpdate | null>(null)
  const [validationErrors, setValidationErrors] = useState<ValidationErrorsU>(
    {},
  )
  const [errorMessage, setErrorMessage] = useState('')

  useEffect(() => {
    if (fetchedWord) setWord(fetchedWord)
  }, [fetchedWord])

  const wordNameRegex = /^[A-Za-z]+$/
  // eslint-disable-next-line no-control-regex
  const japaneseMeanRegex = /^[^-~｡-ﾟ~]+$/

  const validateWord = (target: WordForUpdate) => {
    const ne: ValidationErrorsU = {}
    if (!wordNameRegex.test(target.name))
      ne.name = '単語名は半角アルファベットのみ入力できます。'
    const info = target.wordInfos.map((w) => {
      const e: { partOfSpeech?: string; japaneseMeans?: string[] } = {}
      if (w.partOfSpeechId === 0) e.partOfSpeech = '品詞を選択してください。'
      const ms = w.japaneseMeans.map((m) =>
        japaneseMeanRegex.test(m.name)
          ? ''
          : '日本語訳はひらがな、カタカナ、漢字、または記号「~」のみ入力できます。',
      )
      if (ms.some((x) => x !== '')) e.japaneseMeans = ms
      return e
    })
    if (
      info.some(
        (x) => x.partOfSpeech || (x.japaneseMeans && x.japaneseMeans.length),
      )
    )
      ne.wordInfos = info
    return ne
  }

  const handleWordNameChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    word && setWord({ ...word, name: e.target.value })
  const handlePartOfSpeechChange = (i: number, v: string) => {
    if (!word) return
    const u = [...word.wordInfos]
    u[i].partOfSpeechId = parseInt(v, 10)
    setWord({ ...word, wordInfos: u })
  }
  const handleJapaneseMeanChange = (wi: number, mi: number, v: string) => {
    if (!word) return
    const u = [...word.wordInfos]
    u[wi].japaneseMeans[mi].name = v
    setWord({ ...word, wordInfos: u })
  }

  const updateWordMutation = useMutation<
    { name: string },
    unknown,
    WordForUpdate
  >({
    mutationFn: async (updated: WordForUpdate) =>
      (await axiosInstance.put(`/words/${id}`, updated)).data,
    onSuccess: (data) =>
      navigate(`/words/${id}` as string, {
        state: { successMessage: `${data.name}が正常に更新されました！` },
      }),
    onError: () => setErrorMessage('単語情報の更新中にエラーが発生しました。'),
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!word) return
    const errors = validateWord(word)
    setValidationErrors(errors)
    if (Object.keys(errors).length) return
    updateWordMutation.mutate(word)
  }

  const getAvailablePartOfSpeechOptions = (
    currentIndex: number,
  ): PartOfSpeechOption[] => {
    if (!word) return []
    const selected = word.wordInfos
      .filter((_, i) => i !== currentIndex)
      .map((i) => i.partOfSpeechId)
    return getPartOfSpeech.filter((o) => !selected.includes(o.id))
  }

  if (isLoading) return <p>読み込み中…</p>
  if (isError)
    return (
      <div className="rounded-xl border-l-4 border-red-500 bg-[var(--container_bg)] px-4 py-3 text-sm text-red-600">
        単語情報の取得中にエラーが発生しました。{' '}
        <Button variant="outline" className="ml-2" onClick={() => refetch()}>
          再取得
        </Button>
      </div>
    )
  if (!word) return <p>データが存在しません。</p>

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold text-[var(--h1_fg)]">単語更新</h1>
        <Badge>🛠️ Edit</Badge>
      </div>

      {errorMessage && (
        <div className="mb-4 rounded-xl border-l-4 border-red-500 bg-[var(--container_bg)] px-4 py-3 text-sm text-red-600">
          {errorMessage}
        </div>
      )}

      <Card className="p-6">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="mb-1 block text-sm font-medium">単語名</label>
            <Input value={word.name} onChange={handleWordNameChange} required />
            {validationErrors.name && (
              <p className="mt-1 text-sm text-red-600">
                {validationErrors.name}
              </p>
            )}
          </div>

          {word.wordInfos.map((wi, wiIndex) => {
            const infoErr = validationErrors.wordInfos?.[wiIndex]
            return (
              <Card key={wi.id} className="p-4">
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
                </div>

                <div className="mt-4 space-y-3">
                  {wi.japaneseMeans.map((m, mi) => (
                    <div key={m.id} className="grid gap-2 sm:grid-cols-1">
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
                          required
                        />
                        {infoErr?.japaneseMeans?.[mi] && (
                          <p className="mt-1 text-sm text-red-600">
                            {infoErr.japaneseMeans[mi]}
                          </p>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </Card>
            )
          })}

          <div className="flex flex-wrap items-center gap-3">
            <Button type="submit">単語を更新</Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate(`/words/${word.id}`)}
            >
              単語詳細に戻る
            </Button>
          </div>
        </form>
      </Card>
    </div>
  )
}

export default WordEdit
