import React from 'react'
import { Navigate } from 'react-router-dom'

import { useAuth } from '../hooks/useAuth'

type UserRole = 'guest' | 'general' | 'admin' | 'root'

type PrivateRouteProps = {
  children: React.ReactNode
  requiredRole?: UserRole // 任意: 指定しなければ「ログインさえしていればOK」
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({
  children,
  requiredRole,
}) => {
  const { isLoggedIn, userRole, isLoading } = useAuth()

  if (isLoading) {
    // 認証状態を確認中はローディングを表示
    return <div>Loading...</div>
  }

  if (!isLoggedIn) {
    // 未ログインならリダイレクト
    return <Navigate to="/" />
  }

  if (requiredRole) {
    // もし requiredRole が "admin" なら 'admin' または 'root' は許可、'general' はNG にする等のロジック
    if (requiredRole === 'admin') {
      if (!(userRole === 'admin' || userRole === 'root')) {
        return <Navigate to="/" />
      }
    } else if (requiredRole === 'root') {
      // root 以外は通さない
      if (userRole !== 'root') {
        return <Navigate to="/" />
      }
    }
    // 他にも「一般ユーザー以上」的なロール判定を柔軟に追加できます
  }
  // ログイン済みの場合は子要素を表示
  return <>{children}</>
}

export default PrivateRoute
