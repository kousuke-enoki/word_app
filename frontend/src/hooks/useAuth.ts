import { useEffect, useState } from 'react'

import axiosInstance from '../axiosConfig'

type UserRole = 'guest' | 'general' | 'admin' | 'root'

export const useAuth = () => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)
  const [userRole, setUserRole] = useState<UserRole>('guest')
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const checkAuth = async () => {
      setIsLoading(true)
      const token = localStorage.getItem('token')
      if (!token) {
        // トークンがなければ未ログイン状態
        setIsLoggedIn(false)
        setUserRole('guest')
        setIsLoading(false)
        return
      }
      axiosInstance
        .get('/auth/check')
        .then((response) => {
          const user = response.data.user
          const isLogin = response.data.isLogin
          if (isLogin && user.id) {
            setIsLoggedIn(true)
            if (user.isAdmin) {
              setUserRole('admin')
            }
            if (user.isRoot) {
              setUserRole('root')
            }
          } else {
            setIsLoggedIn(false)
            setUserRole('guest')
          }
        })
        .catch(() => {
          setIsLoggedIn(false)
          setUserRole('guest')
        })
        .finally(() => {
          setIsLoading(false)
        })
    }

    checkAuth()
  }, [])

  return { isLoggedIn, userRole, isLoading }
}
