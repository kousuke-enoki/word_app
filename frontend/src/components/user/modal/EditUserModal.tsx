// src/components/user/modal/EditUserModal.tsx
import React, { useEffect, useMemo, useState } from 'react'

import axiosInstance from '@/axiosConfig'
import Modal from '@/components/common/Modal'
import { Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import type { User } from '@/components/user/UserList'

type RoleKey = 'admin' | 'user' | 'root' | 'test'
const EMAIL_RE = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
const validateName = (s: string) => {
  const t = s.trim()
  return t.length >= 3 && t.length <= 20
}
const validateEmail = (s: string) => EMAIL_RE.test(s.trim())
const validatePasswordNew = (s: string) =>
  s === '' || (s.length >= 8 && s.length <= 72)

type Props = {
  open: boolean
  user: User | null
  isSelf: boolean // 自分編集か？（current PW 必須判定に使用）
  canEditRole: boolean // 呼び出し元で制御（root かつ 対象が root/test 以外の時のみ true）
  onClose: () => void
  onSuccess: (message: string) => void
  onError: (message: string) => void
}

const EditUserModal: React.FC<Props> = ({
  open,
  user,
  isSelf,
  canEditRole,
  onClose,
  onSuccess,
  onError,
}) => {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [role, setRole] = useState<RoleKey>('user')
  const [pwNew, setPwNew] = useState('')
  const [pwCurrent, setPwCurrent] = useState('')
  const [confirmOpen, setConfirmOpen] = useState(false)

  const isTestTarget = !!user?.isTest
  const hasPw = !!user?.isSettedPassword

  // role 初期値（表示のみ。root/test はそもそも編集不可なので何でも良い）
  useEffect(() => {
    if (!open || !user) return
    setName(user.name || '')
    setEmail(user.email || '')
    setPwNew('')
    setPwCurrent('')
    let role: RoleKey = 'user'
    if (user.isAdmin) role = 'admin'
    if (user.isRoot) role = 'root'
    if (user.isTest) role = 'test'
    setRole(role)
  }, [open, user])

  // current パスワードが「必須になる」か
  const needCurrentPw = isSelf && hasPw && pwNew.trim() !== ''

  // 入力バリデーション
  const isValid = useMemo(() => {
    if (isTestTarget) return false // test は編集不可
    if (
      !validateName(name) ||
      !validateEmail(email) ||
      !validatePasswordNew(pwNew)
    )
      return false
    if (needCurrentPw && pwCurrent.trim() === '') return false
    return true
  }, [name, email, pwNew, pwCurrent, needCurrentPw, isTestTarget])

  // 変更のないフィールドは送らない（部分更新）
  const buildPayload = () => {
    if (!user) return null
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {}
    const tName = name.trim()
    const tEmail = email.trim().toLowerCase()

    if (tName !== (user.name || '')) payload.name = tName
    if (tEmail !== (user.email || '').toLowerCase()) payload.email = tEmail
    if (pwNew.trim() !== '') {
      payload.password = { new: pwNew.trim() }
      if (needCurrentPw) payload.password.current = pwCurrent.trim()
    }
    // role は admin / user のみ。変更があり、かつ編集可の時だけ送る
    const targetRole: RoleKey = user.isAdmin ? 'admin' : 'user'
    if (canEditRole && role !== targetRole) {
      payload.role = role
    }
    return payload
  }

  const submit = async () => {
    if (!user) return
    const payload = buildPayload()
    if (!payload || Object.keys(payload).length === 0) {
      onClose()
      return
    }
    try {
      await axiosInstance.put(`/users/${user.id}`, payload)
      onSuccess('ユーザーを更新しました。')
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (e: any) {
      // 代表的なAPIエラーのマッピング
      const status = e?.response?.status
      const data = e?.response?.data
      if (status === 400 && data?.errors?.length) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        onError(data.errors.map((x: any) => x.message).join(' / '))
      } else if (status === 403) {
        onError('権限がありません。')
      } else if (status === 404) {
        onError('ユーザーが見つかりません。')
      } else if (status === 409) {
        onError('このメールは既に使われています。')
      } else {
        onError('ユーザー更新に失敗しました。')
      }
    }
  }

  return (
    <>
      <Modal open={open} onClose={onClose} title="ユーザー編集">
        <div className="space-y-3">
          {isTestTarget && (
            <div className="rounded-md bg-[var(--warn_bg)] px-3 py-2 text-[var(--warn_fg)]">
              テストユーザーは編集できません。
            </div>
          )}

          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              名前
            </label>
            <Input
              value={name}
              disabled={isTestTarget}
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
              disabled={isTestTarget}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="example@domain.com"
            />
            {!validateEmail(email) && (
              <p className="mt-1 text-xs text-red-600">メール形式が不正です</p>
            )}
          </div>

          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              新しいパスワード（空欄＝変更しない）
            </label>
            <Input
              type="password"
              value={pwNew}
              disabled={isTestTarget}
              onChange={(e) => setPwNew(e.target.value)}
              placeholder="8〜72文字"
            />
            {!validatePasswordNew(pwNew) && (
              <p className="mt-1 text-xs text-red-600">
                8〜72文字で入力してください（空欄は未変更）
              </p>
            )}
          </div>

          {needCurrentPw && (
            <div>
              <label className="mb-1 block text-xs text-[var(--muted_fg)]">
                現在のパスワード
              </label>
              <Input
                type="password"
                value={pwCurrent}
                onChange={(e) => setPwCurrent(e.target.value)}
                placeholder="現在のパスワードを入力"
              />
              {pwCurrent.trim() === '' && (
                <p className="mt-1 text-xs text-red-600">
                  現在のパスワードが必要です
                </p>
              )}
            </div>
          )}

          <div>
            <label className="mb-1 block text-xs text-[var(--muted_fg)]">
              役割
            </label>
            <select
              className="w-full rounded-xl border px-3 py-2"
              value={role}
              disabled={isTestTarget || !canEditRole}
              onChange={(e) => setRole(e.target.value as RoleKey)}
            >
              {/* admin / user 以外は出さない（API仕様に準拠） */}
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
            {!canEditRole && (
              <p className="mt-1 text-xs text-[var(--muted_fg)]">
                役割は変更できません（Root/Testは不可、または権限なし）。
              </p>
            )}
          </div>

          <div className="mt-4 flex justify-center gap-2">
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
        <div className="flex justify-center gap-2">
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
