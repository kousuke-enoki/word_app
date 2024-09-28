import React, { useEffect, useState } from 'react';
import axios from 'axios';

const MyPage: React.FC = () => {
  const [user, setUser] = useState<{ name: string } | null>(null);
  const [message, setMessage] = useState('');

  useEffect(() => {
    // ローカルストレージからJWTトークンを取得
    const token = localStorage.getItem('token');
    // トークンがない場合はログインしてくださいと表示
    if (!token) {
      setMessage('ログインしてくださいトークンがないです');
      return;
    }

    // ユーザー情報を取得するためのリクエストを送信
    axios.get('http://localhost:8080/users/my_page', {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((response) => {
        setUser(response.data);
        setMessage('');
      })
      .catch((error) => {
        console.error(error);
        setMessage('ログインしてください');
      });
  }, []);

  // サインアウト処理
  const handleSignOut = () => {
    localStorage.removeItem('token');
    setUser(null);
    setMessage('ログアウトしました');
  };

  // 今日の日付を取得
  const today = new Date().toLocaleDateString();

  return (
    <div>
      {user ? (
        <>
          <h1>マイページ</h1>
          <p>ようこそ、{user.name}さん！</p>
          <p>今日の日付: {today}</p>
          <button onClick={handleSignOut}>サインアウト</button>
        </>
      ) : (
        <p>{message}</p>
      )}
    </div>
  );
};

export default MyPage;
