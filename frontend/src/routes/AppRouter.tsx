import React from 'react'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Home from '../components/user/Home'
import SignIn from '../components/user/SignIn'
import SignUp from '../components/user/SignUp'
import WordNew from '../components/word/WordNew'
import AllWordList from '../components/word/AllWordList'
import WordShow from '../components/word/WordShow'
import Header from '../components/Header'
// import Dashboard from '../components/Dashboard';
// import Footer from '../components/Footer';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/sign_in" element={<SignIn />} />
        <Route path="/sign_up" element={<SignUp />} />
        <Route path="/words/new" element={<WordNew />} />
        <Route path="/words" element={<AllWordList />} />
        <Route path="/words/:id" element={<WordShow />} />
      </Routes>
    </Router>
  )
}

export default AppRouter
