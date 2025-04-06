import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useTheme } from '../../context/ThemeContext'

const Home: React.FC = () => {
  const [message, setMessage] = useState('')
  const { setTheme } = useTheme()

  useEffect(() => {
    const logoutMessage = localStorage.getItem('logoutMessage')
    setTheme('light') 
    if (logoutMessage) {
      setMessage(logoutMessage)
      localStorage.removeItem('logoutMessage')
    }
  }, [])

  return (
    <div>
      <p>トップページです。</p>
      <p>{message}</p>
      <p>
        <Link to="/sign_up">サインアップページ</Link>
      </p>
      <p>
        <Link to="/sign_in">サインインページ</Link>
      </p>
    </div>
  )
}

export default Home
