import React from 'react'
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Navigate,
} from 'react-router-dom'
import Home from '../components/user/Home'
import SignIn from '../components/user/SignIn'
import SignUp from '../components/user/SignUp'
import AllWordList from '../components/word/AllWordList'
import WordShow from '../components/word/WordShow'
import Header from '../components/Header'
import PrivateRoute from '../components/PrivateRoute' // 後述するPrivateRouteをインポート
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
        {/* ログイン必須ページはPrivateRouteで保護 */}
        <Route
          path="/words"
          element={
            <PrivateRoute>
              <AllWordList />
            </PrivateRoute>
          }
        />
        <Route
          path="/words/:id"
          element={
            <PrivateRoute>
              <WordShow />
            </PrivateRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  )
}

export default AppRouter
