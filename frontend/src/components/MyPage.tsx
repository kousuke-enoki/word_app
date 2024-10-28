import React from 'react'
import { Link } from 'react-router-dom'

// Propsの型定義
interface MyPageProps {
  user: { name: string }
  onSignOut: () => void // サインアウト関数を渡す
}

const MyPage: React.FC<MyPageProps> = ({ user, onSignOut }) => {
  // 今日の日付を取得
  const today = new Date().toLocaleDateString()

  return (
    <div>
      <h2>マイページ</h2>
      <p>ようこそ、{user.name}さん！</p>
      <p>今日の日付: {today}</p>
      <p>
        全単語リスト: <Link to="allwordlist">word app</Link>
      </p>
      <button onClick={onSignOut}>サインアウト</button>
    </div>
  )
}

export default MyPage
