// src/pages/user/UserDetail.tsx
import React, { useCallback, useEffect, useMemo, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import axios from '@/axiosConfig'
import PageBottomNav from '@/components/common/PageBottomNav'
import PageTitle from '@/components/common/PageTitle'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import DeleteUserDialog from '@/components/user/modal/DeleteUserDialog'
import EditUserModal from '@/components/user/modal/EditUserModal'
import type { UserDetail } from '@/components/user/UserList'

const UserDetailPage: React.FC = () => {
  const params = useParams<{ id?: string }>()
  const isMe = useLocation().pathname === '/me'
  const navigate = useNavigate()

  const [user, setUser] = useState<UserDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [flash, setFlash] = useState<{
    type: 'success' | 'error'
    text: string
  } | null>(null)

  const fetchDetail = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const url = isMe ? '/users/me' : `/users/${params.id}`
      const { data } = await axios.get<UserDetail>(url)
      setUser(data)
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (e: any) {
      if (e?.response?.status === 401 || e?.response?.status === 403) {
        setError('権限がありません')
      } else if (e?.response?.status === 404) {
        setError('ユーザーが存在しません')
      } else {
        setError('取得に失敗しました')
      }
    } finally {
      setLoading(false)
    }
  }, [isMe, params.id])

  useEffect(() => {
    fetchDetail()
  }, [fetchDetail])

  const roleLabel = useMemo(() => {
    if (!user) return '-'
    return user.isRoot
      ? 'Root'
      : user.isAdmin
        ? 'Admin'
        : user.isTest
          ? 'Test'
          : 'User'
  }, [user])

  const canEditRole =
    !!user &&
    /* 呼び出し側の認可はバックで担保済を前提にUI制御だけ */ !user.isRoot &&
    !user.isTest
  const isSelf = isMe // ID比較にする場合は auth.user.id === user?.id
  const needCurrentPasswordToUpdate = isSelf && !!user?.isSettedPassword

  const fmt = (s?: string) => (s ? new Date(s).toLocaleString() : '-')

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <PageTitle title={isMe ? 'プロフィール' : `ユーザー詳細`} />
        {!isMe && (
          <Button variant="ghost" onClick={() => navigate('/users')}>
            ← 一覧に戻る
          </Button>
        )}
      </div>

      {flash && (
        <div
          className={`rounded-xl border-l-4 px-4 py-3 text-sm ${
            flash.type === 'success'
              ? 'border-green-500 bg-green-50 text-green-800'
              : 'border-red-500 bg-red-50 text-red-800'
          }`}
        >
          {flash.text}
        </div>
      )}

      <Card className="p-4">
        {loading && (
          <div className="animate-pulse space-y-3">
            <div className="h-6 w-40 rounded bg-[var(--skeleton)]" />
            <div className="h-4 w-72 rounded bg-[var(--skeleton)]" />
            <div className="h-4 w-52 rounded bg-[var(--skeleton)]" />
          </div>
        )}

        {!loading && error && (
          <div className="text-[var(--warn_fg)]">{error}</div>
        )}

        {!loading && !error && user && (
          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">名前</div>
              <div className="text-lg">{user.name}</div>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">メール</div>
              <div>{user.email ?? '未設定'}</div>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">種別</div>
              <span className="rounded bg-[var(--badge_bg)] px-2 py-1 text-xs">
                {roleLabel}
              </span>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">
                LINE連携
              </div>
              <span
                className="rounded px-2 py-1 text-xs
                ${user.isLine ? 'bg-[var(--ok_bg)] text-[var(--ok_fg)]' : 'bg-[var(--muted_bg)] text-[var(--muted_fg)]'}"
              >
                {user.isLine ? '連携済' : '未連携'}
              </span>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">PW設定</div>
              <span
                className="rounded px-2 py-1 text-xs
                ${user.isSettedPassword ? 'bg-[var(--ok_bg)] text-[var(--ok_fg)]' : 'bg-[var(--warn_bg)] text-[var(--warn_fg)]'}"
              >
                {user.isSettedPassword ? '設定済' : '未設定'}
              </span>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">
                作成日時
              </div>
              <div>{fmt(user.createdAt)}</div>
            </div>
            <div>
              <div className="text-sm text-[var(--muted_fg)] mb-1">
                更新日時
              </div>
              <div>{fmt(user.updatedAt)}</div>
            </div>

            <div className="mt-4 flex justify-center gap-2">
              <Button onClick={() => setEditOpen(true)}>編集</Button>
              <Button
                variant="outline"
                disabled={user.isRoot}
                onClick={() => setDeleteOpen(true)}
                title={user.isRoot ? 'rootは削除できません' : '削除'}
              >
                削除
              </Button>
            </div>
          </div>
        )}
      </Card>

      <Card className="mt1 p-2">
        <PageBottomNav className="mt-1" showHome inline compact />
      </Card>

      {/* 既存モーダル流用 */}
      <EditUserModal
        open={editOpen}
        user={user || null}
        isSelf={isSelf}
        canEditRole={canEditRole}
        needCurrentPasswordToUpdate={needCurrentPasswordToUpdate}
        onClose={() => setEditOpen(false)}
        onSuccess={async (msg) => {
          setEditOpen(false)
          setFlash({ type: 'success', text: msg })
          await fetchDetail()
        }}
        onError={(msg) => setFlash({ type: 'error', text: msg })}
        operatorIsRoot={user ? user.isRoot : false}
      />

      <DeleteUserDialog
        open={deleteOpen}
        user={user || null}
        onClose={() => setDeleteOpen(false)}
        onSuccess={async (msg) => {
          setDeleteOpen(false)
          setFlash({ type: 'success', text: msg })
          if (isSelf) {
            // 自分削除 → ここでサインアウト処理/遷移
            // await auth.signOut()
            navigate('/')
          } else {
            navigate('/users') // 一覧へ戻る
          }
        }}
        onError={(msg) => setFlash({ type: 'error', text: msg })}
      />
    </div>
  )
}

export default UserDetailPage
