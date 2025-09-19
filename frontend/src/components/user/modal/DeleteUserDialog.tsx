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
}

// rootは消せないようにする。
// 自分は消せないようにする。
const DeleteUserDialog: React.FC<Props> = ({
  open,
  user,
  onClose,
  onSuccess,
  onError,
}) => {
  const submit = async () => {
    if (!user) return
    try {
      await axiosInstance.delete(`/users/${user.id}`)
      onSuccess('ユーザーを削除しました。')
    } catch {
      onError('ユーザー削除に失敗しました。')
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="元に戻せません。削除しますか？">
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onClose}>
          いいえ
        </Button>
        <Button onClick={submit}>はい</Button>
      </div>
    </Modal>
  )
}

export default DeleteUserDialog
