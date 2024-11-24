import React from 'react'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Home from '../pages/user/Home'
import SignIn from '../pages/user/SignIn'
import SignUp from '../pages/user/SignUp'
import AllWordList from '../pages/word/AllWordList'
import WordShow from '../pages/word/WordShow'
import Header from '../components/Header'
// import Dashboard from '../pages/Dashboard';
// import Footer from '../components/Footer';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/sign_in" element={<SignIn />} />
        <Route path="/sign_up" element={<SignUp />} />
        <Route path="/words" element={<AllWordList />} />
        <Route path="/words/:id" element={<WordShow />} />
      </Routes>
    </Router>
  )
}

export default AppRouter
