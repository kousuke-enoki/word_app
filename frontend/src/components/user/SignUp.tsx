/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import { useTheme } from '@/contexts/themeContext'

type FieldError = { field: string; message: string }

const SignUp: React.FC = () => {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [errors, setErrors] = useState<FieldError[]>([])
  const navigate = useNavigate()
  const { setTheme } = useTheme()

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setMessage('')
    setErrors([])
    try {
      const res = await axiosInstance.post('/users/sign_up', {
        name,
        email,
        password,
      })
      const token = res.data?.token
      if (token) localStorage.setItem('token', token)
      localStorage.setItem('logoutMessage', 'サインアップしました。')

      // 任意：ユーザー設定からテーマ反映（失敗しても無視）
      try {
        const cfg = await axiosInstance.get('/setting/user_config')
        setTheme(cfg.data?.is_dark_mode ? 'dark' : 'light')
      } catch {
        /* empty */
      }

      navigate('/my_page')
    } catch (err: any) {
      const fieldErrors: FieldError[] = err?.response?.data?.errors || []
      setErrors(fieldErrors)
      setMessage(err?.response?.data?.message || 'サインアップに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  const getErrorMessages = (field: string) =>
    errors
      .filter((e) => e.field === field)
      .map((e, i) => (
        <p key={`${field}-${i}`} className="mt-1 text-sm text-red-600">
          {e.message}
        </p>
      ))

  return (
    <div className="mx-auto max-w-md">
      <div className="mb-8 text-center">
        <h1 className="text-2xl font-bold text-[var(--h1_fg)]">サインアップ</h1>
        <p className="mt-1 text-sm opacity-70">
          アカウントを作成して学習を始めましょう。
        </p>
      </div>

      <Card className="p-6">
        <form onSubmit={handleSignUp} className="space-y-5">
          <div>
            <label className="mb-1 block text-sm font-medium" htmlFor="name">
              Name
            </label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your Name"
              required
            />
            {getErrorMessages('name')}
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium" htmlFor="email">
              Email
            </label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
            {getErrorMessages('email')}
          </div>

          <div>
            <label
              className="mb-1 block text-sm font-medium"
              htmlFor="password"
            >
              Password
            </label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
            {getErrorMessages('password')}
          </div>

          {message && (
            <div
              className={`rounded-lg border-l-4 ${errors.length ? 'border-red-500 text-red-600' : 'border-[var(--success_pop_bc)]'} bg-[var(--container_bg)] px-3 py-2 text-sm`}
            >
              {message}
            </div>
          )}

          <Button disabled={loading} className="w-full">
            {loading ? 'サインアップ中…' : 'サインアップ'}
          </Button>
        </form>
      </Card>

      <p className="mt-4 text-center text-sm opacity-80">
        すでにアカウントをお持ちですか？{' '}
        <Link className="underline" to="/sign_in">
          サインイン
        </Link>
      </p>
    </div>
  )
}

export default SignUp
