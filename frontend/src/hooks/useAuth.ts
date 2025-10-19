import { useEffect, useState } from 'react'

import axiosInstance from '../axiosConfig'

type UserRole = 'guest' | 'general' | 'admin' | 'root' | 'test'

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

          if (user.id) {
            setIsLoggedIn(true)
            if (user.isAdmin) {
              setUserRole('admin')
            }
            if (user.isRoot) {
              setUserRole('root')
            }
            if (user.isTest) {
              setUserRole('test')
            }
          } else {
            setIsLoggedIn(false)
            setUserRole('guest')
            localStorage.removeItem('token')
          }
        })
        .catch(() => {
          setIsLoggedIn(false)
          setUserRole('guest')
          localStorage.removeItem('token')
        })
        .finally(() => {
          setIsLoading(false)
        })
    }

    checkAuth()
  }, [])

  return { isLoggedIn, userRole, isLoading }
}
