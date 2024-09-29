import React, {useEffect, useState} from 'react';
import { Link } from 'react-router-dom';
import axiosInstance from '../axiosConfig';
import MyPage from '../components/MyPage';

const Home: React.FC = () => {
  const [user, setUser] = useState<{ name: string } | null>(null);
  const [message, setMessage] = useState('');

  useEffect(() => {
    // ローカルストレージからJWTトークンを取得
    const token = localStorage.getItem('token');

    // トークンがない場合はメッセージを表示
    // if (!token) {
    //   setMessage('ログインしてください。トークンがありません。');
    //   return;
    // }

    // ユーザー情報を取得するためのリクエストを送信
    axiosInstance.get('/users/my_page', {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((response) => {
        console.log(response)
        setUser(response.data.user); // ユーザー情報を保存
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
  );
};

export default Home;
