import React from 'react';

// Propsの型定義
interface MyPageProps {
  user: { name: string };
  onSignOut: () => void;  // サインアウト関数を渡す
}

const MyPage: React.FC<MyPageProps> = ({ user, onSignOut }) => {
  console.log("mypage")
  // 今日の日付を取得
  const today = new Date().toLocaleDateString();

  return (
    <div>
      <h2>マイページ</h2>
      <p>ようこそ、{user.name}さん！</p>
      <p>今日の日付: {today}</p>
      <button onClick={onSignOut}>サインアウト</button>
    </div>
  );
};

export default MyPage;
