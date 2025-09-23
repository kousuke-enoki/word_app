/* eslint-disable @typescript-eslint/no-explicit-any */
import { Input } from '@headlessui/react'
import React, { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import axios from '@/axiosConfig'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'

const LineComplete: React.FC = () => {
  const nav = useNavigate()
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState('')

  const tempToken = useMemo(() => sessionStorage.getItem('line_temp_token'), [])
  const suggestedMail = useMemo(
    () => sessionStorage.getItem('line_suggested_mail') || '',
    [],
  )

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!tempToken) {
      setError(
        'セッションが切れました。もう一度LINEログインからやり直してください。',
      )
      return
    }
    if (password.length < 8) {
      setError('パスワードは8文字以上にしてください。')
      return
    }
    if (password !== confirm) {
      setError('確認用パスワードが一致しません。')
      return
    }

    try {
      const { data } = await axios.post('/users/auth/line/complete', {
        temp_token: tempToken,
        password,
      })
      // サーバは { token: string } を返す前提
      localStorage.setItem('token', data.token)
      // 使い終わったら破棄
      sessionStorage.removeItem('line_temp_token')
      sessionStorage.removeItem('line_suggested_mail')

      nav('/mypage', { replace: true })
    } catch (err: any) {
      setError(err?.response?.data?.error || '本登録に失敗しました。')
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <div className="mb-6">
        <h1 className="text-xl font-bold">LINE連携の最終ステップ</h1>
        <p className="text-sm opacity-70 mt-1">
          パスワードを設定して、アカウント作成を完了します。
        </p>
      </div>

      <Card className="p-6 space-y-4">
        {/* 参考表示（メールが来ていれば表示だけ） */}
        {suggestedMail && (
          <div className="text-sm opacity-80">
            連携されたメール: <b>{suggestedMail}</b>
          </div>
        )}

        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium">パスワード</label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.currentTarget.value)}
              placeholder="8文字以上"
              required
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium">
              パスワード（確認）
            </label>
            <Input
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.currentTarget.value)}
              required
            />
          </div>

          {error && (
            <div className="rounded-lg border-l-4 border-red-500 bg-[var(--container_bg)] px-3 py-2 text-sm text-red-600">
              {error}
            </div>
          )}

          <Button className="w-full">本登録してはじめる</Button>
        </form>
      </Card>
    </div>
  )
}

export default LineComplete
