import React from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

type PrivateRouteProps = {
  children: React.ReactNode
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const { isLoggedIn, isLoading } = useAuth()

  if (isLoading) {
    // 認証状態を確認中はローディングを表示
    return <div>Loading...</div>
  }

  if (!isLoggedIn) {
    // 未ログインならリダイレクト
    return <Navigate to="/" />
  }

  // ログイン済みの場合は子要素を表示
  return <>{children}</>
}

export default PrivateRoute
