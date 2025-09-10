/* eslint-disable @typescript-eslint/no-explicit-any */
import { Input } from '@headlessui/react'
import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'

type SettingResponse = { is_line_auth: boolean }

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [lineAuthEnabled, setLineAuthEnabled] = useState(false)
  const [loadingSetting, setLoadingSetting] = useState(true)
  const navigate = useNavigate()

  // ▼ LINE 認証の有効/無効を取得
  useEffect(() => {
    let mounted = true
    ;(async () => {
      try {
        const { data } =
          await axiosInstance.get<SettingResponse>('/setting/auth')
        if (mounted) setLineAuthEnabled(!!data.is_line_auth)
      } catch {
        // 設定取得失敗時はボタン非表示のまま
      } finally {
        if (mounted) setLoadingSetting(false)
      }
    })()
    return () => {
      mounted = false
    }
  }, [])

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const res = await axiosInstance.post('/users/sign_in', {
        email,
        password,
      })
      const token = res.data?.token
      if (token) localStorage.setItem('token', token)
      navigate('/mypage') // ルートが /my_page の場合はここを変更
    } catch (err: any) {
      const fieldError: string =
        err.response?.data?.message || 'サインインに失敗しました'
      setError(fieldError)
    } finally {
      setLoading(false)
    }
  }

  // ▼ LINE ログイン
  const handleLineLogin = () => {
    const base = import.meta.env.VITE_API_URL
    window.location.href = `${base}/users/auth/line/login`
  }

  return (
    <div className="mx-auto max-w-md">
      <div className="mb-8 text-center">
        <h1 className="text-2xl font-bold text-[var(--h1_fg)]">サインイン</h1>
        <p className="mt-1 text-sm opacity-70">
          メールアドレスとパスワードを入力してください。
        </p>
      </div>

      <Card className="p-6 space-y-5">
        {/* メール/パスワード */}
        <form onSubmit={onSubmit} className="space-y-5">
          <div>
            <label className="mb-1 block text-sm font-medium">Email</label>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium">Password</label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>

          {error && (
            <div className="rounded-lg border-l-4 border-red-500 bg-[var(--container_bg)] px-3 py-2 text-sm text-red-600">
              {error}
            </div>
          )}

          <Button disabled={loading} className="w-full">
            {loading ? 'サインイン中…' : 'メールでサインイン'}
          </Button>
        </form>

        {/* 区切り */}
        {!loadingSetting && lineAuthEnabled && (
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t border-[var(--border)]" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-[var(--bg)] px-3 text-xs opacity-60">
                または
              </span>
            </div>
          </div>
        )}

        {/* LINE ログイン（設定が有効時のみ） */}
        {!loadingSetting && lineAuthEnabled && (
          <button
            type="button"
            onClick={handleLineLogin}
            className="
              inline-flex w-full items-center justify-center gap-2 rounded-xl px-4 py-2.5
              text-sm font-medium text-white
              ring-2 ring-transparent focus:outline-none focus:ring-2 focus:ring-[#06C755]
              transition
            "
            style={{
              backgroundColor: '#06C755', // LINE ブランドカラー
            }}
            aria-label="LINEでログイン"
          >
            {/* LINE ロゴ（SVG） */}
            <svg
              width="20"
              height="20"
              viewBox="0 0 36 36"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <rect width="36" height="36" rx="8" fill="white" />
              <path
                d="M18.03 7C12.02 7 7.14 10.99 7.14 15.93c0 3.27 2.2 6.11 5.46 7.57-.17.62-.61 2.2-.7 2.54-.11.43.16.43.33.31.14-.1 2.23-1.52 3.13-2.13.88.13 1.79.2 2.77.2 6.01 0 10.89-3.98 10.89-8.92S24.04 7 18.03 7Z"
                fill="#06C755"
              />
              <path
                d="M12.5 14.3h1.6v5.2h-1.6v-5.2Zm3.05 0h1.6v3.4h2.1v1.8h-3.7v-5.2Zm6.29 0h1.6v5.2h-1.6v-5.2Zm-3.15 0h1.6v5.2h-1.6v-5.2Z"
                fill="white"
              />
            </svg>
            LINEでログイン
          </button>
        )}
        {/* {!loadingSetting && !lineAuthEnabled && (
          <div className="text-center text-xs opacity-60">
            現在、LINEログインは無効です。
          </div>
        )} */}
      </Card>

      <p className="mt-4 text-center text-sm opacity-80">
        アカウント未作成ですか？{' '}
        <Link className="underline" to="/sign_up">
          サインアップ
        </Link>
      </p>
    </div>
  )
}

export default SignIn
