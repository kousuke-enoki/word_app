import React from 'react'
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Navigate,
} from 'react-router-dom'
import PrivateRoute from './PrivateRoute'
import PublicRoute from './PublicRoute'
import Home from '../components/user/Home'
import MyPage from '../components/user/MyPage'
import UserSetting from '../components/setting/UserSetting'
import RootSetting from '../components/setting/RootSetting'
import SignIn from '../components/user/SignIn'
import SignUp from '../components/user/SignUp'
import WordNew from '../components/word/WordNew'
import WordEdit from '../components/word/WordEdit'
import WordList from '../components/word/WordList'
import WordShow from '../components/word/WordShow'
import Header from '../components/Header'
// import Dashboard from '../components/Dashboard';
// import Footer from '../components/Footer';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Header />
      <Routes>
        {/* 未ログインのみアクセス可 */}

        <Route
          path="/"
          element={
            <PublicRoute>
              <Home />
            </PublicRoute>
          }
        />
        <Route
          path="/sign_in"
          element={
            <PublicRoute>
              <SignIn />
            </PublicRoute>
          }
        />
        <Route
          path="/sign_up"
          element={
            <PublicRoute>
              <SignUp />
            </PublicRoute>
          }
        />

        {/* ログイン済みのみアクセス可 */}
        <Route
          path="/mypage"
          element={
            <PrivateRoute>
              <MyPage />
            </PrivateRoute>
          }
        />
        <Route
          path="/user/userSetting"
          element={
            <PrivateRoute requiredRole={'root'}>
              <UserSetting />
            </PrivateRoute>
          }
        />
        <Route
          path="/user/rootSetting"
          element={
            <PrivateRoute requiredRole={'root'}>
              <RootSetting />
            </PrivateRoute>
          }
        />
        <Route
          path="/words/new"
          element={
            <PrivateRoute requiredRole={'admin'}>
              <WordNew />
            </PrivateRoute>
          }
        />
        <Route
          path="/words/edit/:id"
          element={
            <PrivateRoute requiredRole={'admin'}>
              <WordEdit />
            </PrivateRoute>
          }
        />
        <Route
          path="/words"
          element={
            <PrivateRoute>
              <WordList />
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
