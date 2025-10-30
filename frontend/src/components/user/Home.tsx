import React, { useEffect, useMemo, useRef, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { useTheme } from '@/contexts/themeContext'
import { testLogin } from '@/features/auth/testLogin'
import { useTestUserMode } from '@/features/setting/isTestUserMode'

type QuickLink = { title: string; desc: string; to: string; emoji: string }

const LINKS: QuickLink[] = [
  {
    title: '単語リスト',
    desc: '英単語のリストや詳細を表示',
    to: '/words',
    emoji: '📊',
  },
  {
    title: '一括登録',
    desc: '長文をコピペするだけで単語登録',
    to: '/Words/BulkRegister',
    emoji: '✍️',
  },
  { title: 'クイズ', desc: '10問からすぐに開始', to: '/quizs', emoji: '🧠' },
]

const Home: React.FC = () => {
  const [message, setMessage] = useState('')
  const [uiError, setUiError] = useState<string>('')
  const { setTheme } = useTheme()
  const [testing, setTesting] = useState(false) // API呼び出し中
  const [btnCooldown, setBtnCooldown] = useState(false) // 2〜5秒のUIクールダウン
  const testModeEnabled = useTestUserMode()
  const navigate = useNavigate()

  // 誤タップ連打対策（画面側クールダウン）
  const cooldownTimer = useRef<number | null>(null)
  const startCooldown = (ms = 2500) => {
    setBtnCooldown(true)
    if (cooldownTimer.current) window.clearTimeout(cooldownTimer.current)
    cooldownTimer.current = window.setTimeout(() => setBtnCooldown(false), ms)
  }

  useEffect(() => {
    const logoutMessage = localStorage.getItem('logoutMessage')
    setTheme('light') // 初期テーマを設定
    if (logoutMessage) {
      setMessage(logoutMessage)
      localStorage.removeItem('logoutMessage')
    }
    return () => {
      if (cooldownTimer.current) window.clearTimeout(cooldownTimer.current)
    }
  }, [setTheme])

  // 既にログイン済みか（トークンの有無だけの簡易判定）
  const isAuthed = useMemo(() => !!localStorage.getItem('token'), [])

  const doTestLoginThen = async (to: string) => {
    if (testing) return
    setUiError('')
    try {
      setTesting(true)
      startCooldown()
      // ログイン済みなら再発行しない
      if (!isAuthed) {
        await testLogin() // 1分以内の連打は同じ結果を返す
      }
      navigate(to)
    } catch (e: unknown) {
      if (
        typeof e === 'object' &&
        e !== null &&
        'remainingMs' in e &&
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        typeof (e as any).remainingMs === 'number'
      ) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const sec = Math.ceil((e as any).remainingMs / 1000)
        setUiError(`テストログインは ${sec} 秒後に再度お試しください。`)
      } else {
        setUiError(
          'テストログインに失敗しました。時間をおいて再度お試しください。',
        )
      }
    } finally {
      setTesting(false)
    }
  }

  const handleTestLoginClick = async () => {
    if (!testModeEnabled) return
    if (btnCooldown) {
      setUiError('操作が早すぎます。数秒後にもう一度お試しください。')
      return
    }
    await doTestLoginThen('/words')
  }

  const handleCardClick = (to: string) => async (e: React.MouseEvent) => {
    if (!testModeEnabled) {
      e.preventDefault()
      return
    }
    if (!isAuthed) {
      e.preventDefault()
      if (btnCooldown) {
        setUiError('操作が早すぎます。数秒後にもう一度お試しください。')
        return
      }
      await doTestLoginThen(to)
    }
  }

  return (
    <section className="relative">
      <div className="pointer-events-none absolute inset-0 -z-10 bg-gradient-to-b from-[var(--container_bg)]/60 to-transparent" />

      {/* 成功/失敗メッセージ */}
      {(message || uiError) && (
        <div
          className={`mb-6 rounded-xl border-l-4 px-4 py-3 text-sm ${
            uiError
              ? 'border-red-500 bg-[var(--container_bg)] text-red-600'
              : 'border-[var(--success_pop_bc)] bg-[var(--container_bg)]'
          }`}
          role="status"
          aria-live="polite"
        >
          {uiError || message}
        </div>
      )}

      <div className="text-center">
        <h1 className="mb-3 text-3xl font-bold tracking-tight text-[var(--h1_fg)]">
          英単語を、もっと覚えやすく。
        </h1>
        <p className="mx-auto mb-8 max-w-2xl text-[15px] opacity-80">
          単語登録・クイズ・成績可視化までを一つに。ポートフォリオとしても、普遍的に使える学習体験を目指しています。
        </p>
        <div className="flex items-center justify-center gap-3">
          <Link to="/sign_in">
            <Button>サインイン</Button>
          </Link>
          {/* テストモードON時のみ表示 */}
          {testModeEnabled && (
            <Button
              variant="primary"
              onClick={handleTestLoginClick}
              disabled={testing || btnCooldown}
              aria-disabled={testing || btnCooldown}
              title={testing ? 'テストログイン中…' : 'ワンクリックで試す'}
            >
              {testing ? 'テストログイン中…' : 'テストログイン'}
            </Button>
          )}
        </div>
      </div>

      {/* 下部カード：テストモードOFF時は非活性表示（クリック無効） */}
      <div className="mt-12 grid gap-4 sm:grid-cols-3">
        {LINKS.map((i) => {
          const disabled = !testModeEnabled && !isAuthed
          return (
            <Link
              key={i.title}
              to={disabled ? '#' : i.to}
              onClick={handleCardClick(i.to)}
              className={`group ${disabled ? 'pointer-events-none opacity-50' : ''}`}
              aria-disabled={disabled}
            >
              <Card className="h-full p-5 transition hover:shadow-md">
                <div className="mb-2 text-sm opacity-70">{i.emoji}</div>
                <div className="text-base font-semibold">{i.title}</div>
                <p className="mt-1 text-sm opacity-70">{i.desc}</p>
              </Card>
            </Link>
          )
        })}
      </div>
    </section>
  )
}

export default Home
