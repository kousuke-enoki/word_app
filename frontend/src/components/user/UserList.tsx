// src/pages/user/UserList.tsx
// import '@/styles/components/user/UserList.css' // 必要ならCSSを用意（WordListと共通でもOK）

import React, { useEffect, useState } from 'react'
import { Link, useLocation } from 'react-router-dom'

import axiosInstance from '@/axiosConfig'
import PageBottomNav from '@/components/common/PageBottomNav'
import PageTitle from '@/components/common/PageTitle'
import Pagination from '@/components/common/Pagination'
import { Badge, Card, Input } from '@/components/ui/card'
import { Button } from '@/components/ui/ui'

// --- types（必要なら src/types/userTypes.ts に分離してください） ---
type User = {
  id: number
  name: string
  email?: string
  isAdmin: boolean
  isRoot: boolean
  isTest: boolean
  isSettedPassword?: boolean
  isLine?: boolean
}

type UserListResponse = {
  users: User[]
  totalPages: number
}

// --- component ---
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

  useEffect(() => {
    if (!isInitialized) return
    const fetchUsers = async () => {
      try {
        const { data } = await axiosInstance.get<UserListResponse>('/users', {
          params: { search, sortBy, order, page, limit },
        })
        setUsers(data.users)
        setTotalPages(data.totalPages)
      } catch (e) {
        console.error('Failed to fetch users:', e)
      }
    }
    fetchUsers()
  }, [search, sortBy, order, page, limit, isInitialized])

  const roleLabel = (u: User) => {
    if (u.isRoot) return 'Root'
    if (u.isAdmin) return 'Admin'
    if (u.isTest) return 'Test'
    return 'User'
  }

  const roleBadgeTone = (u: User) => {
    if (u.isRoot) return 'bg-[var(--badge_root_bg)] text-[var(--badge_root_fg)]'
    if (u.isAdmin)
      return 'bg-[var(--badge_admin_bg)] text-[var(--badge_admin_fg)]'
    if (u.isTest) return 'bg-[var(--badge_test_bg)] text-[var(--badge_test_fg)]'
    return 'bg-[var(--badge_user_bg)] text-[var(--badge_user_fg)]'
  }

  const Toolbar = (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <div className="flex-1 min-w-[220px]">
        <Input
          value={search}
          onChange={(e) => {
            setPage(1) // 新しい検索時は1ページ目へ
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
          // ソートキー変更時は先頭ページに戻す
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
  )

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <PageTitle title="ユーザー一覧" />
        {/* もしユーザー新規作成ページがあるなら */}
        {/* <Link to="/users/new"><Button>新規作成</Button></Link> */}
      </div>

      <Card className="p-4">
        {Toolbar}

        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="bg-[var(--thbc)] text-left">
                {['名前', 'メール', '役割', 'LINE連携', 'PW設定', '詳細'].map(
                  (th) => (
                    <th
                      key={th}
                      className="border-b border-[var(--thbd)] px-3 py-2 text-[var(--fg)]"
                    >
                      {th}
                    </th>
                  ),
                )}
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
                    <Link
                      to={`/users/${u.id}`}
                      state={{ search, sortBy, order, page, limit }}
                      className="underline"
                    >
                      詳細
                    </Link>
                  </td>
                </tr>
              ))}
              {users.length === 0 && (
                <tr>
                  <td
                    className="px-3 py-6 text-center text-[var(--muted_fg)]"
                    colSpan={6}
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
    </div>
  )
}

export default UserList
