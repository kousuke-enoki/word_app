import React, { useMemo, useState } from 'react'

import axiosInstance from '@/axiosConfig'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'

import PageBottomNav from '../common/PageBottomNav'
import PageTitle from '../common/PageTitle'

type Token = { word: string; checked: boolean }

// 最大文字数
const MAX_LEN = 5000

const WordBulkRegister: React.FC = () => {
  // ① textarea 入力値
  const [text, setText] = useState('')

  // ② サーバ応答
  const [tokens, setTokens] = useState<Token[]>([])
  const [notExistWords, setNotExistWords] = useState<string[]>([])
  const [registeredWords, setRegisteredWords] = useState<string[]>([])

  // ③ フィルタ
  const [filter, setFilter] = useState('')

  // ④ ロード/メッセージ
  const [loading, setLoading] = useState(false)
  const [msg, setMsg] = useState('')
  const [registedMsg, setRegistedMsg] = useState('')

  // ---------- 抽出 ----------
  const handleExtract = async () => {
    if (!text.trim()) return
    setLoading(true)
    setMsg('')
    setTokens([])
    setNotExistWords([])
    setRegisteredWords([])
    try {
      const { data } = await axiosInstance.post('/words/bulk_tokenize', {
        text,
      })
      console.log(data)
      const cands = (data.candidates || []) as string[]
      const notExists = (data.not_exists || []) as string[]
      const regs = (data.registered || []) as string[]

      if (cands.length > 0) {
        setTokens(cands.map((w) => ({ word: w, checked: false })))
        setMsg(`抽出に成功しました（${cands.length} 語）`)
      } else {
        setTokens([])
        setMsg('登録できる単語がありませんでした。')
      }
      setNotExistWords(notExists)
      setRegisteredWords(regs)
    } catch (e: any) {
      console.log(e)
      const errorMsg = e?.response?.data?.error || '抽出に失敗しました'
      setMsg(errorMsg)
      if (e?.response?.status === 429) {
        setMsg('1日のリクエスト上限に達しました')
      }
    } finally {
      setLoading(false)
    }
  }

  // ---------- 一括登録 ----------
  const handleRegister = async () => {
    const selected = tokens.filter((t) => t.checked).map((t) => t.word)
    if (selected.length === 0 || selected.length > 200) return
    setLoading(true)
    setRegistedMsg('')
    try {
      const { data } = await axiosInstance.post('/words/bulk_register', {
        words: selected,
      })
      console.log(data)
      let resMsg = ''
      if (data.success && data.failed) {
        resMsg = `結果： ${data.success.length} 件登録 / 失敗 ${data.failed.length} 件`
      } else if (data.success) {
        resMsg = `結果： ${data.success.length} 件登録`
      } else if (data.failed) {
        resMsg = `結果： ${data.failed.length} 件失敗`
      }
      setRegistedMsg(resMsg)

      // 登録成功した単語をチェック解除
      if (data.success && data.success.length > 0) {
        const successSet = new Set(data.success)
        setTokens((prev) =>
          prev.map((t) =>
            successSet.has(t.word) ? { ...t, checked: false } : t,
          ),
        )
      }
    } catch (e: any) {
      console.log(e)
      const errorMsg = e?.response?.data?.error || '登録に失敗しました'
      setRegistedMsg(errorMsg)
      if (e?.response?.status === 429) {
        setRegistedMsg('1日のリクエスト上限に達しました')
      }
    } finally {
      setLoading(false)
    }
  }

  // ---------- 初期化 ----------
  const handleReset = () => {
    setText('')
    setTokens([])
    setNotExistWords([])
    setRegisteredWords([])
    setFilter('')
    setMsg('')
    setRegistedMsg('')
  }

  // ---------- ソート/フィルタ ----------
  const sorted = useMemo(() => {
    const f = filter.trim().toLowerCase()
    return tokens
      .filter((t) => (f ? t.word.toLowerCase().includes(f) : true))
      .sort((a, b) => a.word.localeCompare(b.word))
  }, [tokens, filter])

  // ---------- チェック数 ----------
  const checkedCount = useMemo(
    () => tokens.filter((t) => t.checked).length,
    [tokens],
  )

  // 全選択/全解除
  const toggleAll = (next: boolean) =>
    setTokens((prev) => prev.map((p) => ({ ...p, checked: next })))

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <PageTitle title="単語一括登録" />
        <Badge>β</Badge>
      </div>
      {/* 入力カード */}
      <Card className="mb-6 p-6">
        <div className="mb-3 flex items-center justify-between text-sm opacity-80">
          <span>英語の長文を貼り付けてください（{MAX_LEN} 文字まで）</span>
          <span>
            {text.length}/{MAX_LEN}
          </span>
        </div>

        <textarea
          className="w-full rounded-xl border border-[var(--textarea_bd,var(--input_bd))] bg-[var(--textarea_bg)] p-3 text-[var(--textarea_c)] outline-none focus:ring-2 ring-[var(--button_bg)]"
          rows={8}
          maxLength={MAX_LEN}
          placeholder="Paste your English paragraph here..."
          value={text}
          onChange={(e) => setText(e.target.value)}
        />

        <div className="mt-4 flex flex-wrap items-center gap-3">
          <Button onClick={handleExtract} disabled={loading || !text.trim()}>
            {loading ? '抽出中…' : '抽出'}
          </Button>
          <Button
            variant="outline"
            onClick={handleReset}
            disabled={loading && !tokens.length}
          >
            初期化
          </Button>

          {msg && (
            <span className="ml-auto rounded-lg border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-3 py-1.5 text-sm">
              {msg}
            </span>
          )}
        </div>
      </Card>
      {/* 候補カード */}
      {tokens.length > 0 && (
        <Card className="mb-6 p-6">
          {/* ツールバー */}
          <div className="mb-4 flex flex-wrap items-center gap-3">
            <div className="min-w-[220px] flex-1">
              <Input
                placeholder="検索..."
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
              />
            </div>
            <Badge>候補: {sorted.length}</Badge>
            <Badge>選択: {checkedCount} / 200</Badge>
            <Button variant="outline" onClick={() => toggleAll(true)}>
              全選択
            </Button>
            <Button variant="outline" onClick={() => toggleAll(false)}>
              全解除
            </Button>
            <Button
              className="ml-auto"
              onClick={handleRegister}
              disabled={loading || checkedCount === 0}
            >
              {loading ? '登録中…' : 'まとめて登録'}
            </Button>
          </div>

          {/* 候補一覧 */}
          <div className="max-h-[360px] overflow-y-auto rounded-xl border border-[var(--border)] p-3">
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
              {sorted.map((t) => (
                <label
                  key={t.word}
                  className={`flex cursor-pointer items-center gap-2 rounded-lg border px-3 py-2 text-sm ${
                    t.checked
                      ? 'border-[var(--btn-subtle-bd)] bg-[var(--btn-subtle-bg)]'
                      : 'border-[var(--border)] bg-[var(--container_bg)]'
                  }`}
                >
                  <input
                    type="checkbox"
                    checked={t.checked}
                    onChange={() =>
                      setTokens((prev) =>
                        prev.map((p) =>
                          p.word === t.word ? { ...p, checked: !p.checked } : p,
                        ),
                      )
                    }
                  />
                  <span className="truncate">{t.word}</span>
                </label>
              ))}
            </div>
          </div>

          {registedMsg && (
            <div className="mt-4 rounded-lg border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-3 py-2 text-sm">
              {registedMsg}
            </div>
          )}
        </Card>
      )}
      {/* 補助情報 */}
      {registeredWords.length > 0 && (
        <Card className="mb-4 p-5">
          <div className="mb-2 text-sm font-semibold">すでに登録済みの単語</div>
          <div className="max-h-40 overflow-y-auto">
            <div className="flex flex-wrap gap-2">
              {registeredWords.map((w) => (
                <Badge key={`reg-${w}`}>{w}</Badge>
              ))}
            </div>
          </div>
        </Card>
      )}
      {notExistWords.length > 0 && (
        <Card className="mb-4 p-5">
          <div className="mb-2 text-sm font-semibold">
            データが存在しないため登録できない単語
          </div>
          <div className="max-h-40 overflow-y-auto">
            <div className="flex flex-wrap gap-2">
              {notExistWords.map((w) => (
                <Badge key={`no-${w}`}>{w}</Badge>
              ))}
            </div>
          </div>
        </Card>
      )}{' '}
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

export default WordBulkRegister
