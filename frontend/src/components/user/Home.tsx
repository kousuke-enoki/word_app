import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { Card,PageContainer } from '@/components/card'
import { PageShell } from '@/components/PageShell'
import { Button } from '@/components/ui'
import { useTheme } from '@/contexts/themeContext'

const Home: React.FC = () => {
  // const [message, setMessage] = useState('')
  const [message, setMessage] = useState('')
  // const logoutMessage = localStorage.getItem('logoutMessage') || ''
  const { setTheme } = useTheme()

  useEffect(() => {
    const logoutMessage = localStorage.getItem('logoutMessage')
    setTheme('light')   // 初期テーマを設定
    if (logoutMessage) {
      setMessage(logoutMessage)
      // setMessage(logoutMessage)
      localStorage.removeItem('logoutMessage')
    }
  }, [setTheme])

  return (
    <PageShell>
      <section className="relative">
      <div className="pointer-events-none absolute inset-0 -z-10 bg-gradient-to-b from-[var(--container_bg)]/60 to-transparent" />
      <PageContainer>
        {message ? (
          <div className="mb-6 rounded-xl border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-4 py-3 text-sm">
            {message}
          </div>
        ) : null}


      <div className="text-center">
        <h1 className="mb-3 text-3xl font-bold tracking-tight text-[var(--h1_fg)]">英単語を、もっと覚えやすく。</h1>
        <p className="mx-auto mb-8 max-w-2xl text-[15px] opacity-80">
          単語登録・クイズ・成績可視化までを一つに。ポートフォリオとしても、普遍的に使える学習体験を目指しています。
        </p>
        <div className="flex items-center justify-center gap-3">
          <Link to="/sign_in">
            <Button>サインイン</Button>
          </Link>
        </div>
      </div>


      <div className="mt-12 grid gap-4 sm:grid-cols-3">
        {[
          { title: '単語リスト', desc: '英単語のリストや詳細を表示', to: '/words' },
          { title: '登録', desc: '英単語と意味、品詞などを登録', to: '/words/new' },
          { title: 'クイズ', desc: '10問からすぐに開始', to: '/quizs' },
        ].map((i) => (
        <Link key={i.title} to={i.to} className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-2 text-sm opacity-70">{i.title === '登録' ? '✍️' : i.title === 'クイズ' ? '🧠' : '📊'}</div>
            <div className="text-base font-semibold">{i.title}</div>
            <p className="mt-1 text-sm opacity-70">{i.desc}</p>
          </Card>
        </Link>
        ))}
      </div>
    </PageContainer>
    </section>
  </PageShell>
  )
}

export default Home
