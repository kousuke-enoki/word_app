import React, { useEffect, useState } from 'react'
import axiosInstance from '@/axiosConfig'
import { EditingPermission } from '@/types/settingTypes'
import '@/styles/components/setting/RootSetting.css'

const RootSetting: React.FC = () => {
  /** 取得した初期値を保持しておく */
  const [initial, setInitial] = useState<{
    editing: EditingPermission
    test: boolean
    mail: boolean
    line: boolean
  }>({
    editing: 'admin',
    test: false,
    mail: false,
    line: false,
  })

  /** フォームの現在値 */
  const [editingPermission, setEditingPermission] =
    useState<EditingPermission>('admin')
  const [isTestUserMode, setIsTestUserMode] = useState(false)
  const [isEmailAuthCheck, setIsEmailAuthCheck] = useState(false)
  const [isLineAuth, setIsLineAuth] = useState(false)

  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)

  /* ─────────  初期値取得  ───────── */
  useEffect(() => {
    ;(async () => {
      try {
        const { data } = await axiosInstance.get('/setting/root_config')
        const config = data.Config
        setEditingPermission(config.editing_permission)
        setIsTestUserMode(config.is_test_user_mode)
        setIsEmailAuthCheck(config.is_email_authentication_check)
        setIsLineAuth(config.is_line_authentication)

        // 初期値を保持
        setInitial({
          editing: config.editing_permission,
          test: config.is_test_user_mode,
          mail: config.is_email_authentication_check,
          line: config.is_line_authentication,
        })
      } catch (e){
        alert('ルート設定の取得中にエラーが発生しました。')
        console.error(e)
      }
    })()
  }, [])

  /* ─────────  変更検知  ───────── */
  const isDirty =
    editingPermission !== initial.editing ||
    isTestUserMode !== initial.test ||
    isEmailAuthCheck !== initial.mail ||
    isLineAuth !== initial.line

  /* ─────────  保存処理  ───────── */
  const handleSave = async () => {
    setLoading(true)
    setMessage('')
    try {
      const res = await axiosInstance.post('/setting/root_config', {
        editing_permission: editingPermission,
        is_test_user_mode: isTestUserMode,
        is_email_authentication_check: isEmailAuthCheck,
        is_line_authentication: isLineAuth,
      })
      if (res.status === 200) {
        setMessage('設定が保存されました。')
        // 保存後を新しい初期値として採用
        setInitial({
          editing: editingPermission,
          test: isTestUserMode,
          mail: isEmailAuthCheck,
          line: isLineAuth,
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
          checked={isEmailAuthCheck}
          onChange={(e) => setIsEmailAuthCheck(e.target.checked)}
        />
      </div>

      <div className="row">
        <label>LINE認証：</label>
        <input
          type="checkbox"
          checked={isLineAuth}
          onChange={(e) => setIsLineAuth(e.target.checked)}
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
