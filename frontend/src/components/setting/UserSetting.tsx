import '@/styles/components/setting/UserSetting.css'

import React, { useEffect, useState } from 'react'

import axiosInstance from '@/axiosConfig'
import { useTheme } from '@/contexts/themeContext'

const UserSetting: React.FC = () => {
  const [initial, setInitial] = useState<{
    isDarkMode: boolean
  }>({
    isDarkMode: false
  })

  const [isDarkMode, SetIsDarkMode] = useState(false)
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)
  const { setTheme } = useTheme()

  useEffect(() => {
    const fetchUserSettingData = async () => {
      try {
        const response = await axiosInstance.get('/setting/user_config')
        const config = response.data.Config
        SetIsDarkMode(config.is_dark_mode)
        // 初期値を保持
        setInitial({
          isDarkMode: config.is_dark_mode,
        })
      } catch (error) {
        console.error(error)
        alert('ユーザー設定の取得中にエラーが発生しました。')
      }
    }
    fetchUserSettingData()
  }, [])

  const isDirty =
  isDarkMode !== initial.isDarkMode

  const handleSave = async () => {
    setLoading(true)
    setMessage('')
    try {
      const response = await axiosInstance.post('/setting/user_config', {
        is_dark_mode: isDarkMode,
      })
      if (response.status === 200) {
        setTheme(isDarkMode ? 'dark' : 'light') 
        setMessage('設定を保存しました。')
        setInitial({
          isDarkMode: isDarkMode,
        })
      } else {
        setMessage('保存に失敗しました。')
      }
    } catch (error) {
      console.error(error)
      setMessage('保存中にエラーが発生しました。')
    } finally {
      setLoading(false)
    }
  }

  return (
    // <div style={{ padding: '2rem', maxWidth: '500px', margin: '0 auto' }}>
    <div className="settings-wrapper">
      <h2 className="title">ユーザー設定画面</h2>

      {/* <div style={{ marginBottom: '1rem' }}> */}
      <div className="row">
        <label>ダークモード：</label>
        <input
          type="checkbox"
          checked={isDarkMode}
          onChange={(e) => SetIsDarkMode(e.target.checked)}
        />
        <span style={{ marginLeft: '0.5rem' }}>{isDarkMode ? 'ON' : 'OFF'}</span>
      </div>

      <button onClick={handleSave} disabled={loading || !isDirty}>
        {loading ? '保存中...' : '保存'}
      </button>

      {message && <p style={{ marginTop: '1rem' }}>{message}</p>}
    </div>
  )
}

export default UserSetting
