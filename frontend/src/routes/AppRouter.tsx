import React from 'react'
import {
  BrowserRouter as Router,
  Navigate,
  Route,
  Routes,
} from 'react-router-dom'

import { PageContainer } from '@/components/ui/card'
import { PageShell } from '@/components/ui/PageShell'

import QuizMenu from '../components/quiz/QuizMenu'
import ResultIndex from '../components/result/ResultIndex'
import ResultShow from '../components/result/ResultShow/ResultShow'
import RootSetting from '../components/setting/RootSetting'
import UserSetting from '../components/setting/UserSetting'
import Home from '../components/user/Home'
import LineCallback from '../components/user/LineCallback'
import MyPage from '../components/user/MyPage'
import SignIn from '../components/user/SignIn'
import SignUp from '../components/user/SignUp'
import WordBulkRegister from '../components/word/WordBulkRegister'
import WordEdit from '../components/word/WordEdit'
import WordList from '../components/word/WordList'
import WordNew from '../components/word/WordNew'
import WordShow from '../components/word/WordShow'
import PrivateRoute from './PrivateRoute'
import PublicRoute from './PublicRoute'

const routerFutureFlags = {
  v7_startTransition: true,
  v7_relativeSplatPath: true,
}

const AppRouter: React.FC = () => {
  return (
    <Router future={routerFutureFlags}>
      <PageShell>
        <PageContainer>
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
            <Route
              path="/line/callback"
              element={
                <PublicRoute>
                  <LineCallback />
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
                <PrivateRoute requiredRole={'general'}>
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
            <Route
              path="/Words/BulkRegister"
              element={
                <PrivateRoute>
                  <WordBulkRegister />
                </PrivateRoute>
              }
            />
            <Route
              path="/quizs"
              element={
                <PrivateRoute>
                  <QuizMenu />
                </PrivateRoute>
              }
            />
            <Route
              path="/results"
              element={
                <PrivateRoute>
                  <ResultIndex />
                </PrivateRoute>
              }
            />
            <Route
              path="/results/:quizNo"
              element={
                <PrivateRoute>
                  <ResultShow />
                </PrivateRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </PageContainer>
      </PageShell>
    </Router>
  )
}

export default AppRouter
