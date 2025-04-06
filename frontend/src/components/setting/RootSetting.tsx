import React, { useEffect, useState } from 'react'
import axiosInstance from '../../axiosConfig'
import { EditingPermission } from '../../types/settingTypes'
import '../../styles/components/setting/RootSetting.css'

const RootSetting: React.FC = () => {
  /** 取得した初期値を保持しておく */
  const [initial, setInitial] = useState<{
    editing: EditingPermission
    test: boolean
    mail: boolean
  }>({ 
    editing: 'admin',
    test: false,
    mail: false,
  })

  /** フォームの現在値 */
  const [editingPermission, setEditingPermission] =
    useState<EditingPermission>('admin')
  const [isTestUserMode, setIsTestUserMode] = useState(false)
  const [isEmailAuth, setIsEmailAuth] = useState(false)

  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)

  /* ─────────  初期値取得  ───────── */
  useEffect(() => {
    ;(async () => {
      try {
        const res = await axiosInstance.get('/setting/root_config')
        setEditingPermission(res.data.editing_permission)
        setIsTestUserMode(res.data.is_test_user_mode)
        setIsEmailAuth(res.data.is_email_authentication)

        // 初期値を保持
        setInitial({
          editing: res.data.editing_permission,
          test: res.data.is_test_user_mode,
          mail: res.data.is_email_authentication,
        })
      } catch {
        alert('ルート設定の取得中にエラーが発生しました。')
      }
    })()
  }, [])

  /* ─────────  変更検知  ───────── */
  const isDirty =
    editingPermission !== initial.editing ||
    isTestUserMode !== initial.test ||
    isEmailAuth !== initial.mail

  /* ─────────  保存処理  ───────── */
  const handleSave = async () => {
    setLoading(true)
    setMessage('')
    try {
      const res = await axiosInstance.post('/setting/root_config', {
        editing_permission: editingPermission,
        is_test_user_mode: isTestUserMode,
        is_email_authentication: isEmailAuth,
      })
      if (res.status === 200) {
        setMessage('設定が保存されました。')
        // 保存後を新しい初期値として採用
        setInitial({
          editing: editingPermission,
          test: isTestUserMode,
          mail: isEmailAuth,
        })
      } else {
        setMessage('保存に失敗しました。')
      }
    } catch (e) {
      setMessage('エラーが発生しました。')
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="settings-wrapper">
      <h2 className="title">管理設定画面</h2>

      {/* 行 */}
      <div className="row">
        <label>編集権限ロール：</label>
        <select
          value={editingPermission}
          onChange={(e) => setEditingPermission(e.target.value as EditingPermission)}
        >
          <option value="user">一般ユーザー</option>
          <option value="admin">adminユーザー</option>
          <option value="root">ルートユーザー</option>
        </select>
      </div>

      <div className="row">
        <label>テストユーザーモード：</label>
        <input
          type="checkbox"
          checked={isTestUserMode}
          onChange={(e) => setIsTestUserMode(e.target.checked)}
        />
      </div>

      <div className="row">
        <label>メール認証：</label>
        <input
          type="checkbox"
          checked={isEmailAuth}
          onChange={(e) => setIsEmailAuth(e.target.checked)}
        />
      </div>

      <button onClick={handleSave} disabled={loading || !isDirty}>
        {loading ? '保存中...' : '保存'}
      </button>

      {message && <p className="msg">{message}</p>}

    </div>
  )
}

export default RootSetting
