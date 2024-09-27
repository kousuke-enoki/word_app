import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div>
      <h1>word app</h1>
      <p>トップページです。</p>
      <p>
        <Link to="/SignUp">サインアップページ</Link>
      </p>
      <p>
        <Link to="/SignIn">サインインページ</Link>
      </p>
      <p>
        <Link to="/mypage">マイページ</Link>
      </p>
    </div>
  );
};

export default Home;
