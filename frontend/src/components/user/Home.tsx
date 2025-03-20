import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

const Home: React.FC = () => {
  const [message, setMessage] = useState('')

  useEffect(() => {
    const logoutMessage = localStorage.getItem('logoutMessage')
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
