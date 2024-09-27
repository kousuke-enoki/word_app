import React from 'react';
import { Link } from 'react-router-dom';

const Header: React.FC = () => {
  return (
    <div>
      <p>
        <Link to="/">word app</Link>
      </p>
      <p>
      <Link to="/mypage">マイページ</Link>
      </p>
    </div>
  );
};

export default Header;
