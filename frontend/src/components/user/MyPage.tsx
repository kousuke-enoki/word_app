import React, { useEffect, useState } from 'react'
import axiosInstance from '../../axiosConfig'
import { Link, useNavigate } from 'react-router-dom'
import { User } from '../../types/userTypes'

const MyPage: React.FC = () => {
  const [message] = useState(() => localStorage.getItem('logoutMessage') || '')
  const [user, setUser] = useState<User | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await axiosInstance.get('/users/my_page')
        setUser(response.data.user)

        // ログアウトメッセージがあれば表示し、一度表示したら削除
        if (message) {
          localStorage.removeItem('logoutMessage')
        }
      } catch (error) {
        localStorage.removeItem('token')
        localStorage.setItem('logoutMessage', 'ログインしてください')
        setTimeout(() => {
          navigate('/')
        }, 2000)
      }
    }

    fetchUserData()
  }, [message, navigate])

  // 今日の日付を取得
  const today = new Date().toLocaleDateString()

  let showRole = <></>
  if (user?.root) {
    showRole = <p>ルートユーザーでログインしています。</p>
  } else if (user?.admin) {
    showRole = <p>管理ユーザーでログインしています。</p>
  }

  const handleSignOut = () => {
    localStorage.removeItem('token')
    localStorage.setItem('logoutMessage', 'ログアウトしました')
    setUser(null)

    setTimeout(() => {
      navigate('/')
    })
  }

  return (
    <div>
      <h2>マイページ</h2>
      {message && <p>{message}</p>}
      {showRole}
      {user ? (
        <p>ようこそ、{user.name}さん！</p>
      ) : (
        <p>ユーザー情報がありません。</p>
      )}
      <p>今日の日付: {today}</p>
      <p>
        全単語リスト: <Link to="/words">word list</Link>
      </p>
      {user?.admin ? (
        <p>
          <Link to="/words/new">単語登録画面</Link>
        </p>
      ) : null}
      <button onClick={handleSignOut}>サインアウト</button>
    </div>
  )
}

export default MyPage
