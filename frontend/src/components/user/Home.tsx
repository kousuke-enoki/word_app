import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useTheme } from '@/contexts/themeContext'

const Home: React.FC = () => {
  const [message, setMessage] = useState('')
  const { setTheme } = useTheme()

  useEffect(() => {
    const logoutMessage = localStorage.getItem('logoutMessage')
    setTheme('light')   // 初期テーマを設定
    if (logoutMessage) {
      setMessage(logoutMessage)
      localStorage.removeItem('logoutMessage')
    }
  }, [setTheme])

  return (
    <div>
      <p>トップページです。</p>
      <p>{message}</p>
      <p>
        <Link to="/sign_in">サインインはここから！</Link>
      </p>
    </div>
  )
}

export default Home
