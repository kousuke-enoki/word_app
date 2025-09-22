// src/components/user/DeleteUserDialog.tsx
import React from 'react'

import axiosInstance from '@/axiosConfig'
import Modal from '@/components/common/Modal'
import { Button } from '@/components/ui/ui'
import type { User } from '@/components/user/UserList'

type Props = {
  open: boolean
  user: User | null
  onClose: () => void
  onSuccess: (message: string) => void
  onError: (message: string) => void
  /** 追加：対象ユーザーが自分自身かどうか（親から渡す） */
  isSelf?: boolean
}

const DeleteUserDialog: React.FC<Props> = ({
  open,
  user,
  onClose,
  onSuccess,
  onError,
  isSelf = false,
}) => {
  const cannotDelete = !!user && (isSelf || user.isRoot)

  const submit = async () => {
    if (!user) return
    // ブロック対象は実行不可（念のため submit 側でも防止）
    if (cannotDelete) {
      onError('自分自身やrootユーザーを削除することはできません。')
      return
    }
    try {
      await axiosInstance.delete(`/users/${user.id}`)
      onSuccess('ユーザーを削除しました。')
    } catch {
      onError('ユーザー削除に失敗しました。')
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="元に戻せません。削除しますか？">
      <div className="space-y-3">
        {cannotDelete && (
          <p className="rounded-md bg-red-50 p-3 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-300">
            自分自身やrootユーザーを削除することはできません。
          </p>
        )}
        <div className="flex justify-center gap-2">
          <Button variant="outline" onClick={onClose}>
            いいえ
          </Button>
          <Button onClick={submit} disabled={cannotDelete}>
            はい
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default DeleteUserDialog
