import React from 'react'
import { Link } from 'react-router-dom'
import { User } from '../../types/userTypes'

// Propsの型定義
interface MyPageProps {
  user: User
  onSignOut: () => void // サインアウト関数を渡す
}

const MyPage: React.FC<MyPageProps> = ({ user, onSignOut }) => {
  // 今日の日付を取得
  const today = new Date().toLocaleDateString()
  let isAdmin = <></>
  if (user.admin) {
    isAdmin = <p>管理ユーザーでログインしています。</p>
  }

  return (
    <div>
      <h2>マイページ</h2>
      {isAdmin}
      <p>ようこそ、{user.name}さん！</p>
      <p>今日の日付: {today}</p>
      <p>
        全単語リスト: <Link to="words">word app</Link>
      </p>
      <p>
        <Link to="words/new">単語登録画面</Link>
      </p>
      <p>
        <Link to="exams">test</Link>
      </p>
      <button onClick={onSignOut}>サインアウト</button>
    </div>
  )
}

export default MyPage
