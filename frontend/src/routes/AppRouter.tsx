import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from '../pages/user/Home';
import SignIn from '../pages/user/SignIn';
import SignUp from '../pages/user/SignUp';
import AllWordList from '../pages/word/AllWordList';
import Header from '../components/Header';
// import Dashboard from '../pages/Dashboard';
// import Footer from '../components/Footer';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/allwordlist" element={<AllWordList />} />
        {/* 他のルートをここに追加 */}
      </Routes>
      {/* <Footer /> */}
    </Router>
  );
};

export default AppRouter;