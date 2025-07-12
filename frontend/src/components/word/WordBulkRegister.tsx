/* WordBulkRegister.tsx */
import React, { useState, useMemo } from 'react'
import axiosInstance from '@/axiosConfig'


type Token = { word: string; checked: boolean }

const MAX_LEN = 3000

const WordBulkRegister: React.FC = () => {
  /* ① textarea 入力値 */
  const [text, setText] = useState('')

  /* ② 抽出後サーバから返ってくる配列 */
  const [tokens, setTokens] = useState<Token[]>([])
  const [notExistWords, setNotExistWords] = useState<string[]>([])
  const [registeredWords, setRegisteredWords] = useState<string[]>([])

  /* ③ フィルタ文字列 */
  const [filter, setFilter] = useState('')

  /* ④ ロード／メッセージ */
  const [loading, setLoading] = useState(false)
  const [msg, setMsg] = useState('')
  const [registedMsg, setRegistedMsg] = useState('')

  /* ---------- 抽出処理 ---------- */
  const handleExtract = async () => {
    if (!text.trim()) return
    setLoading(true)
    setMsg('')
    try {
      const { data } = await axiosInstance.post('/words/bulk_tokenize', { text })
      const cands = data.candidates as string[]
      const notExists = data.not_exists as string[]
      const regs = data.registered as string[]
      
      if(cands && cands?.length > 0){
        setTokens(cands.map(w => ({ word: w, checked: false })))
        setMsg(
          cands.length > 0
            ? `抽出に成功しました（${cands.length} 語）`
            : '登録できる単語がありませんでした。',
        )
      } else {
        setMsg('登録できる単語がありませんでした。')
      }
      setNotExistWords(notExists)
      setRegisteredWords(regs)

    } catch {
      setMsg('抽出に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  /* ───────── 登録処理 ───────── */
  const handleRegister = async () => {
    const selected = tokens.filter(t => t.checked).map(t => t.word)
    if (selected.length === 0 || selected.length > 200) return
    setLoading(true)
    setRegistedMsg('')
    try {
      const {data} = await axiosInstance.post('/words/bulk_register', { words: selected })
      let resMsg = ''
      if (data.success && data.failed) {
        resMsg = `結果： ${data.success.length} 件登録 / 失敗 ${data.failed.length} 件`
      } else if (data.success) {
        resMsg = `結果： ${data.success.length} 件登録`
      } else if (data.failed) {
        resMsg = `結果： ${data.failed.length} 件失敗`
      }
      setRegistedMsg(resMsg)
    } catch {
      setRegistedMsg('登録に失敗しました')
    } finally {
      setLoading(false)
    }
  }

  /* ---------- フォーム初期化 ---------- */
  const handleReset = () => {
    setText('')
    setTokens([])
    setNotExistWords([])
    setRegisteredWords([])
    setFilter('')
    setMsg('')
  }

  /* ───────── ソート / フィルタ ───────── */
  const sorted = useMemo(() => {
    return tokens
      .filter(t => t.word.includes(filter.toLowerCase()))
      .sort((a, b) => a.word.localeCompare(b.word))
  }, [tokens, filter])

  /* ---------- チェック時 ---------- */
  const checkedCount = useMemo(
    () => tokens.filter(t => t.checked).length,
    [tokens],
  )

   /* ---------- UI ---------- */
   return (
    <div style={{ maxWidth: 700, margin: '0 auto', padding: 20 }}>
      <h2>単語一括登録</h2>
      <label>5000文字まで送信できます。</label>

      <textarea
        rows={6}
        maxLength={MAX_LEN}
        style={{ width: '100%', marginBottom: 4 }}
        placeholder="英語の長文を貼り付けてください"
        value={text}
        onChange={e => setText(e.target.value)}
      />

      <div style={{ fontSize: 12, textAlign: 'right', marginBottom: 8 }}>
        {text.length}/{MAX_LEN}
      </div>

      <button onClick={handleExtract} disabled={loading}>
        抽出
      </button>
      <button onClick={handleReset} style={{ marginLeft: 8 }} disabled={loading}>
        初期化
      </button>

      {msg && <p style={{ marginTop: 16 }}>{msg}</p>}

      {/* ソート・検索 */}
      {tokens && tokens.length > 0 && (
        <>
          <div style={{ margin: '16px 0' }}>
            <input
              placeholder="検索..."
              value={filter}
              onChange={e => setFilter(e.target.value)}
            />
            <span style={{ marginLeft: 8 }}>件数: {sorted.length}</span>
          </div>

          {/* 候補一覧（チェック可） */}
          <div
            style={{
              maxHeight: 300,
              overflowY: 'auto',
              border: '1px solid #ccc',
              padding: 8,
            }}
          >
            {sorted.map((t) => (
              <label key={t.word} style={{ display: 'block' }}>
                <input
                  type="checkbox"
                  checked={t.checked}
                  onChange={() =>{
                    // handleCheck
                    setTokens(prev =>
                      prev.map(p =>
                        p.word === t.word ? { ...p, checked: !p.checked } : p,
                      ),
                    )
                  }}
                />
                {t.word}
              </label>
            ))}
          </div>
          <label>同時に登録できるのは200個までです。チェック：{checkedCount}個</label>
          <button
            style={{ marginTop: 12 }}
            onClick={handleRegister}
            disabled={loading || tokens.every(t => !t.checked)}
          >
            まとめて登録
          </button>
          {registedMsg && <p style={{ marginTop: 16 }}>{registedMsg}</p>}
        </>
      )}

      {/* 登録済み単語 */}
      {registeredWords && registeredWords.length > 0 ? (
        <div style={{ marginTop: 24 }}>
          <h4>すでに登録済みの単語</h4>
          <div>{registeredWords.join(', ')}</div>
        </div>
      ):null}

      {/* 登録できない単語 */}
      {notExistWords && notExistWords.length > 0 ? (
        <div style={{ marginTop: 16 }}>
          <h4>データが存在しないため登録できない単語</h4>
          <div>{notExistWords.join(', ')}</div>
        </div>
      ):null}
    </div>
  )
}

export default WordBulkRegister