import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Badge, Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { clearTestLoginCache } from '@/features/auth/testLogin'

import { User } from '../../types/userTypes'
import PageTitle from '../common/PageTitle'

const MyPage: React.FC = () => {
  const [message] = useState(() => localStorage.getItem('logoutMessage') || '')
  const [user, setUser] = useState<User | null>(null)
  const [signingOut, setSigningOut] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    const run = async () => {
      try {
        const res = await axiosInstance.get('/users/my_page')
        setUser(res.data.user)
        if (message) localStorage.removeItem('logoutMessage')
      } catch {
        localStorage.removeItem('token')
        localStorage.setItem('logoutMessage', 'ログインしてください')
        setTimeout(() => navigate('/'), 1500)
      }
    }
    run()
  }, [message, navigate])

  const today = new Date().toLocaleDateString()

  const signOutLocally = (msg: string) => {
    localStorage.removeItem('token')
    clearTestLoginCache() // テストログインキャッシュもクリア
    localStorage.setItem('logoutMessage', msg)
    setUser(null)
    navigate('/')
  }

  // 通常サインアウト
  const onSignOut = () => {
    if (signingOut) return
    signOutLocally('ログアウトしました')
  }

  // テストユーザー専用：確認→削除API→トークン破棄
  const onTestLogout = async () => {
    if (signingOut) return
    const ok = window.confirm(
      'テストユーザーはアカウントごと削除されます。よろしいですか？',
    )
    if (!ok) return
    setSigningOut(true)
    try {
      await axiosInstance.post('/users/auth/test-logout')
      // 成功でも失敗でもローカルは破棄&遷移
      signOutLocally(
        'テストユーザーを削除しました。ご利用ありがとうございました。',
      )
    } catch {
      // 失敗しても冪等：トークン破棄してトップへ
      signOutLocally('削除に失敗しましたが、サインアウトしました')
    } finally {
      setSigningOut(false)
    }
  }

  const renderSignOutButton = () => {
    if (!user) return null
    if (user.isTest) {
      return (
        <Button onClick={onTestLogout} disabled={signingOut}>
          {signingOut ? '処理中…' : 'テストユーザー削除（サインアウト）'}
        </Button>
      )
    }
    return (
      <Button onClick={onSignOut} disabled={signingOut}>
        {signingOut ? 'サインアウト中…' : 'サインアウト'}
      </Button>
    )
  }

  return (
    <div>
      {message && (
        <div className="mb-4 rounded-xl border-l-4 border-[var(--success_pop_bc)] bg-[var(--container_bg)] px-4 py-3 text-sm">
          {message}
        </div>
      )}

      <div className="mb-6 flex items-center justify-between">
        <div>
          <PageTitle title="マイページ" />
          <p className="mt-1 text-sm opacity-80">今日の日付: {today}</p>
        </div>
        <div>
          {user?.isRoot ? (
            <Badge>⭐ Root</Badge>
          ) : user?.isAdmin ? (
            <Badge>🔧 Admin</Badge>
          ) : user?.isTest ? (
            <Badge>👾 Test</Badge>
          ) : (
            <Badge>👤 User</Badge>
          )}
        </div>
      </div>

      <Card className="mb-6 p-5">
        {user ? (
          <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm opacity-70">ようこそ</p>
              <p className="text-lg font-semibold">{user.name} さん</p>
            </div>
            <div className="mt-3 sm:mt-0">{renderSignOutButton()}</div>
          </div>
        ) : (
          <p>ユーザー情報がありません。</p>
        )}
      </Card>

      <div className="grid gap-4 sm:grid-cols-2">
        <Link to="/me" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">🙋</div>
            <div className="text-base font-semibold">ユーザー情報詳細</div>
            <p className="mt-1 text-sm opacity-70">
              登録情報の確認・編集・削除
            </p>
          </Card>
        </Link>

        {user?.isRoot && (
          <Link to="/users" className="group">
            <Card className="h-full p-5 transition hover:shadow-md">
              <div className="mb-1 text-sm opacity-70">🤖</div>
              <div className="text-base font-semibold">ユーザーリスト</div>
              <p className="mt-1 text-sm opacity-70">
                検索・ソート・ページネーションに対応
              </p>
            </Card>
          </Link>
        )}

        <Link to="/words" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">📚</div>
            <div className="text-base font-semibold">全単語リスト</div>
            <p className="mt-1 text-sm opacity-70">
              検索・ソート・ページネーションに対応
            </p>
          </Card>
        </Link>

        <Link to="/Words/BulkRegister" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">📥</div>
            <div className="text-base font-semibold">まとめて登録</div>
            <p className="mt-1 text-sm opacity-70">
              英文のコピペで楽に登録可能
            </p>
          </Card>
        </Link>

        {user?.isAdmin && (
          <Link to="/words/new" className="group">
            <Card className="h-full p-5 transition hover:shadow-md">
              <div className="mb-1 text-sm opacity-70">✍️</div>
              <div className="text-base font-semibold">単語登録</div>
              <p className="mt-1 text-sm opacity-70">新しい単語を追加</p>
            </Card>
          </Link>
        )}

        <Link to="/quizs" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">🧠</div>
            <div className="text-base font-semibold">単語クイズ</div>
            <p className="mt-1 text-sm opacity-70">10問から手軽に開始</p>
          </Card>
        </Link>

        <Link to="/results" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">📊</div>
            <div className="text-base font-semibold">クイズ成績一覧</div>
            <p className="mt-1 text-sm opacity-70">
              進捗を確認して学習を最適化
            </p>
          </Card>
        </Link>

        <Link to="/user/userSetting" className="group">
          <Card className="h-full p-5 transition hover:shadow-md">
            <div className="mb-1 text-sm opacity-70">⚙️</div>
            <div className="text-base font-semibold">ユーザー設定</div>
            <p className="mt-1 text-sm opacity-70">テーマ設定など</p>
          </Card>
        </Link>

        {user?.isRoot && (
          <Link to="/user/rootSetting" className="group">
            <Card className="h-full p-5 transition hover:shadow-md">
              <div className="mb-1 text-sm opacity-70">🛡️</div>
              <div className="text-base font-semibold">管理設定</div>
              <p className="mt-1 text-sm opacity-70">ルート設定にアクセス</p>
            </Card>
          </Link>
        )}
      </div>
    </div>
  )
}

export default MyPage
