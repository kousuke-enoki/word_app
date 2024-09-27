import React, { useEffect, useState } from 'react';
import axiosInstance from '../api/axiosConfig';

const MyPage: React.FC = () => {
  const [name, setName] = useState('');
  const [date, setDate] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    const fetchUserData = async () => {
    const token = localStorage.getItem('token');
      try {
        const response = await axiosInstance.get('/users/mypage', {
          // 'headers': {
          //   'Authorization': `Bearer ${token}`,
          // },
        });
        setName(response.data.name);
        setDate(response.data.date);
      } catch (error) {
        setMessage('Failed to fetch user data');
      }
    };

    fetchUserData();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setMessage('Logged out successfully');
  };

  return (
    <div>
      <h1>マイページ</h1>
      <p>今日の日付: {date}</p>
      <p>ユーザー名: {name}</p>
      <button onClick={handleLogout}>ログアウト</button>
      {message && <p>{message}</p>}
    </div>
  );
};

export default MyPage;
