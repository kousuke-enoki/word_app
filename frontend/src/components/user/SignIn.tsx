import { Input } from '@headlessui/react'
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Link } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import { Card, PageContainer } from '@/components/card'
import { PageShell } from '@/components/PageShell'
import { Button } from '@/components/ui'

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('root@example.com')
  const [password, setPassword] = useState('Password123$')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()

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
      if (token) {
        localStorage.setItem('token', token)
      }
      navigate('/my_page')
    } catch {
      // setError(err?.response?.data?.message || 'サインインに失敗しました')
    } finally {
      setLoading(false)
    }
  }

  return (
    <PageShell>
      <PageContainer>
        <div className="mx-auto max-w-md">
          <div className="mb-8 text-center">
            <h1 className="text-2xl font-bold text-[var(--h1_fg)]">
              サインイン
            </h1>
            <p className="mt-1 text-sm opacity-70">
              メールアドレスとパスワードを入力してください。
            </p>
          </div>

          <Card className="p-6">
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
                <label className="mb-1 block text-sm font-medium">
                  Password
                </label>
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
                {loading ? 'サインイン中…' : 'サインイン'}
              </Button>
            </form>
          </Card>

          <p className="mt-4 text-center text-sm opacity-80">
            アカウント未作成ですか？{' '}
            <Link className="underline" to="/sign_up">
              サインアップ
            </Link>
          </p>
        </div>
      </PageContainer>
    </PageShell>
  )
}

export default SignIn
