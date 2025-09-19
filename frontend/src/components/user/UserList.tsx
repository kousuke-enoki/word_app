// src/pages/user/UserList.tsx
import React, { useCallback, useEffect, useMemo, useState } from 'react'
import { useLocation } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import PageBottomNav from '@/components/common/PageBottomNav'
import PageTitle from '@/components/common/PageTitle'
import Pagination from '@/components/common/Pagination'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'
import DeleteUserDialog from '@/components/user/modal/DeleteUserDialog'
import EditUserModal from '@/components/user/modal/EditUserModal'

// ※必要であれば src/types/userTypes.ts に分離可
export type User = {
  id: number
  name: string
  email?: string
  isAdmin: boolean
  isRoot: boolean
  isTest: boolean
  isSettedPassword?: boolean
  isLine?: boolean
}
type UserListResponse = { users: User[]; totalPages: number }

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([])
  const [search, setSearch] = useState('')
  const [sortBy, setSortBy] = useState<'name' | 'email' | 'role'>('name')
  const [order, setOrder] = useState<'asc' | 'desc'>('asc')
  const location = useLocation()
  const [page, setPage] = useState<number>(location.state?.page || 1)
  const [totalPages, setTotalPages] = useState(1)
  const [limit, setLimit] = useState(10)
  const [isInitialized, setIsInitialized] = useState(false)

  // フラッシュメッセージ
  const [flash, setFlash] = useState<{
    type: 'success' | 'error'
    text: string
  } | null>(null)

  // モーダル制御
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [target, setTarget] = useState<User | null>(null)

  useEffect(() => {
    if (location.state) {
      setSearch(location.state.search || '')
      setSortBy((location.state.sortBy as typeof sortBy) || 'name')
      setOrder((location.state.order as 'asc' | 'desc') || 'asc')
      setPage(location.state.page || 1)
      setLimit(location.state.limit || 10)
    }
    setIsInitialized(true)
  }, [location.state])

  const fetchUsers = useCallback(async () => {
    const { data } = await axiosInstance.get<UserListResponse>('/users', {
      params: { search, sortBy, order, page, limit },
    })
    setUsers(data.users)
    setTotalPages(data.totalPages)
  }, [search, sortBy, order, page, limit])

  useEffect(() => {
    if (!isInitialized) return
    fetchUsers().catch(() =>
      setFlash({ type: 'error', text: 'ユーザー取得に失敗しました。' }),
    )
  }, [isInitialized, fetchUsers])

  const roleLabel = (u: User) =>
    u.isRoot ? 'Root' : u.isAdmin ? 'Admin' : u.isTest ? 'Test' : 'User'
  const roleBadgeTone = (u: User) =>
    u.isRoot
      ? 'bg-[var(--badge_root_bg)] text-[var(--badge_root_fg)]'
      : u.isAdmin
        ? 'bg-[var(--badge_admin_bg)] text-[var(--badge_admin_fg)]'
        : u.isTest
          ? 'bg-[var(--badge_test_bg)] text-[var(--badge_test_fg)]'
          : 'bg-[var(--badge_user_bg)] text-[var(--badge_user_fg)]'

  const Toolbar = useMemo(
    () => (
      <div className="mb-4 flex flex-wrap items-center gap-3">
        <div className="flex-1 min-w-[220px]">
          <Input
            value={search}
            onChange={(e) => {
              setPage(1)
              setSearch(e.target.value)
            }}
            placeholder="ユーザー名・メール検索"
          />
        </div>

        <select
          className="rounded-xl border border-[var(--input_bd)] bg-[var(--select)] px-3 py-2 text-[var(--select_c)]"
          value={sortBy}
          onChange={(e) => {
            const v = e.target.value as 'name' | 'email' | 'role'
            if (v !== sortBy) setPage(1)
            setSortBy(v)
          }}
        >
          <option value="name">名前</option>
          <option value="email">メール</option>
          <option value="role">役割</option>
        </select>

        <Button
          variant="outline"
          onClick={() => {
            setPage(1)
            setOrder(order === 'asc' ? 'desc' : 'asc')
          }}
        >
          {order === 'asc' ? '昇順' : '降順'}
        </Button>

        <Badge>総ページ: {totalPages}</Badge>
      </div>
    ),
    [search, sortBy, order, totalPages],
  )

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <PageTitle title="ユーザー一覧" />
      </div>

      {flash && (
        <div
          className={`mb-4 rounded-xl border-l-4 px-4 py-3 text-sm ${
            flash.type === 'success'
              ? 'border-green-500 bg-green-50 text-green-800'
              : 'border-red-500 bg-red-50 text-red-800'
          }`}
        >
          {flash.text}
        </div>
      )}

      <Card className="p-4">
        {Toolbar}

        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="bg-[var(--thbc)] text-left">
                {[
                  '名前',
                  'メール',
                  '役割',
                  'LINE連携',
                  'PW設定',
                  '編集',
                  '削除',
                ].map((th) => (
                  <th
                    key={th}
                    className="border-b border-[var(--thbd)] px-3 py-2 text-[var(--fg)]"
                  >
                    {th}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id} className="even:bg-[var(--table_tr_e)]">
                  <td className="px-3 py-2">{u.name}</td>
                  <td className="px-3 py-2">{u.email ?? '-'}</td>
                  <td className="px-3 py-2">
                    <span
                      className={`rounded px-2 py-1 text-xs ${roleBadgeTone(u)}`}
                    >
                      {roleLabel(u)}
                    </span>
                  </td>
                  <td className="px-3 py-2">
                    {u.isLine ? (
                      <span className="rounded bg-[var(--ok_bg)] px-2 py-1 text-xs text-[var(--ok_fg)]">
                        連携済
                      </span>
                    ) : (
                      <span className="rounded bg-[var(--muted_bg)] px-2 py-1 text-xs text-[var(--muted_fg)]">
                        未連携
                      </span>
                    )}
                  </td>
                  <td className="px-3 py-2">
                    {u.isSettedPassword ? (
                      <span className="rounded bg-[var(--ok_bg)] px-2 py-1 text-xs text-[var(--ok_fg)]">
                        設定済
                      </span>
                    ) : (
                      <span className="rounded bg-[var(--warn_bg)] px-2 py-1 text-xs text-[var(--warn_fg)]">
                        未設定
                      </span>
                    )}
                  </td>
                  <td className="px-3 py-2">
                    <Button
                      onClick={() => {
                        setTarget(u)
                        setEditOpen(true)
                      }}
                    >
                      編集
                    </Button>
                  </td>
                  <td className="px-3 py-2">
                    <Button
                      variant="outline"
                      onClick={() => {
                        setTarget(u)
                        setDeleteOpen(true)
                      }}
                    >
                      削除
                    </Button>
                  </td>
                </tr>
              ))}
              {users.length === 0 && (
                <tr>
                  <td
                    className="px-3 py-6 text-center text-[var(--muted_fg)]"
                    colSpan={7}
                  >
                    該当するユーザーが見つかりませんでした
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <Pagination
          className="mt-4"
          page={page}
          totalPages={totalPages}
          onPageChange={(p) => setPage(p)}
          pageSize={limit}
          onPageSizeChange={(n) => {
            setPage(1)
            setLimit(n)
          }}
          pageSizeOptions={[10, 20, 30, 50]}
        />
      </Card>

      <Card className="mt1 p-2">
        <PageBottomNav className="mt-1" showHome inline compact />
      </Card>

      {/* 編集モーダル（別コンポーネント） */}
      <EditUserModal
        open={editOpen}
        user={target}
        onClose={() => setEditOpen(false)}
        onSuccess={async (msg) => {
          setEditOpen(false)
          setFlash({ type: 'success', text: msg })
          await fetchUsers()
        }}
        onError={(msg) => setFlash({ type: 'error', text: msg })}
      />

      {/* 削除モーダル（別コンポーネント） */}
      <DeleteUserDialog
        open={deleteOpen}
        user={target}
        onClose={() => setDeleteOpen(false)}
        onSuccess={async (msg) => {
          setDeleteOpen(false)
          setFlash({ type: 'success', text: msg })
          await fetchUsers()
        }}
        onError={(msg) => setFlash({ type: 'error', text: msg })}
      />
    </div>
  )
}

export default UserList
