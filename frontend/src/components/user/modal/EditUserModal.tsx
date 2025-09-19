// src/components/user/EditUserModal.tsx
import React, { useEffect, useMemo, useState } from 'react'

import axiosInstance from '@/axiosConfig'
import Modal from '@/components/common/Modal'
import { Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import type { User } from '@/components/user/UserList'

type RoleKey = 'root' | 'admin' | 'test' | 'user'

const EMAIL_RE = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
const validateName = (s: string) =>
  s.trim().length >= 3 && s.trim().length <= 20
const validateEmail = (s: string) => EMAIL_RE.test(s)
const validatePassword = (s: string) =>
  s === '' ? true : s.length >= 8 && s.length <= 72

function roleFromUser(u?: User | null): RoleKey {
  if (!u) return 'user'
  if (u.isRoot) return 'root'
  if (u.isAdmin) return 'admin'
  if (u.isTest) return 'test'
  return 'user'
}

type Props = {
  open: boolean
  user: User | null
  onClose: () => void
  onSuccess: (message: string) => void
  onError: (message: string) => void
}

const EditUserModal: React.FC<Props> = ({
  open,
  user,
  onClose,
  onSuccess,
  onError,
}) => {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState<RoleKey>('user')

  const [confirmOpen, setConfirmOpen] = useState(false)
  const isValid = useMemo(
    () =>
      validateName(name) && validateEmail(email) && validatePassword(password),
    [name, email, password],
  )

  useEffect(() => {
    if (!open || !user) return
    setName(user.name || '')
    setEmail(user.email || '')
    setPassword('')
    setRole(roleFromUser(user))
  }, [open, user])

  const submit = async () => {
    if (!user) return
    const payload: {
      name: string
      email: string
      role: RoleKey
      password?: string
    } = { name: name.trim(), email: email.trim(), role }
    if (password.trim() !== '') payload.password = password
    try {
      await axiosInstance.put(`/users/${user.id}`, payload)
      onSuccess('ユーザーを更新しました。')
    } catch {
      onError('ユーザー更新に失敗しました。')
    }
  }

  return (
    <>
      <Modal open={open} onClose={onClose} title="ユーザー編集">
        <div className="space-y-3">
          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              名前
            </label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="3〜20文字"
            />
            {!validateName(name) && (
              <p className="mt-1 text-xs text-red-600">
                3〜20文字で入力してください
              </p>
            )}
          </div>
          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              メール
            </label>
            <Input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@domain.com"
            />
            {!validateEmail(email) && (
              <p className="mt-1 text-xs text-red-600">メール形式が不正です</p>
            )}
          </div>
          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              パスワード（空欄＝変更しない）
            </label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="8〜72文字"
            />
            {!validatePassword(password) && (
              <p className="mt-1 text-xs text-red-600">
                8〜72文字で入力してください（空欄は未変更）
              </p>
            )}
          </div>
          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              役割
            </label>
            <select
              className="w-full rounded-xl border px-3 py-2"
              value={role}
              onChange={(e) => setRole(e.target.value as RoleKey)}
            >
              <option value="user">User</option>
              <option value="test">Test</option>
              <option value="admin">Admin</option>
              <option value="root">Root</option>
            </select>
          </div>

          <div className="mt-4 flex justify-end gap-2">
            <Button variant="outline" onClick={onClose}>
              戻る
            </Button>
            <Button disabled={!isValid} onClick={() => setConfirmOpen(true)}>
              編集確認
            </Button>
          </div>
        </div>
      </Modal>

      <Modal
        open={confirmOpen}
        onClose={() => setConfirmOpen(false)}
        title="更新しますか？"
      >
        <div className="flex justify-end gap-2">
          <Button variant="outline" onClick={() => setConfirmOpen(false)}>
            いいえ
          </Button>
          <Button
            onClick={() => {
              setConfirmOpen(false)
              submit()
            }}
          >
            はい
          </Button>
        </div>
      </Modal>
    </>
  )
}

export default EditUserModal
