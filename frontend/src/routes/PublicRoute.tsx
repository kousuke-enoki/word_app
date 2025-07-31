// PublicRoute.tsx
import React from 'react'
import { Navigate } from 'react-router-dom'

import { useAuth } from '../hooks/useAuth'

type PublicRouteProps = {
  children: React.ReactNode
}

// 未ログイン時だけアクセス可
// ログイン済みなら /mypage へリダイレクト
const PublicRoute: React.FC<PublicRouteProps> = ({ children }) => {
  const { isLoggedIn, isLoading } = useAuth()

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (isLoggedIn) {
    return <Navigate to="/mypage" />
  }

  return <>{children}</>
}

export default PublicRoute
