import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import axiosInstance from '../../axiosConfig'
import MyPage from '../../components/MyPage'
import { User } from '../../types/userTypes'

const Home: React.FC = () => {
  const [user, setUser] = useState<User | null>(null)
  const [message, setMessage] = useState('')

  useEffect(() => {
    // ローカルストレージからJWTトークンを取得
    const token = localStorage.getItem('token')
    if (!token) {
      return
    }
    // ユーザー情報を取得するためのリクエストを送信
    axiosInstance
      .get('/users/my_page', {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((response) => {
        setUser(response.data.user) // ユーザー情報を保存
        setMessage('')
      })
      .catch((error) => {
        console.error(error)
        localStorage.removeItem('token')
        setMessage('ログインしてください')
      })
  }, [])

  // サインアウト処理
  const handleSignOut = () => {
    localStorage.removeItem('token')
    setUser(null)
    setMessage('ログアウトしました')
  }

  return (
    <div>
      {/* ログイン状態の時はマイページを表示 */}
      {user ? (
        <MyPage user={user} onSignOut={handleSignOut} />
      ) : (
        <>
          <p>トップページです。</p>
          <p>{message}</p>
          <p>
            <Link to="/SignUp">サインアップページ</Link>
          </p>
          <p>
            <Link to="/SignIn">サインインページ</Link>
          </p>
        </>
      )}
    </div>
  )
}

export default Home
