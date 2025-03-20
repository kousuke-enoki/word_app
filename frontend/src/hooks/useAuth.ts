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
          const data = response.data
          if (data.isLogin) {
            setIsLoggedIn(true)
            // setUserRole(data.userRole)
            if (data.isAdmin) {
              setUserRole('admin')
            }
            if (data.isRoot) {
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
    // ↓無限ループしてしまうので入れない
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { isLoggedIn, userRole, isLoading }
}
