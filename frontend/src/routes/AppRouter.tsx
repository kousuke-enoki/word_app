import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from '../pages/Home';
import SignIn from '../pages/SignIn';
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
        {/* 他のルートをここに追加 */}
      </Routes>
      {/* <Footer /> */}
    </Router>
  );
};

export default AppRouter;