// axiosConfig.ts
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// axiosのインスタンスを作成
const axiosInstance = axios.create({
  baseURL: '/api',      // 相対パスにする
  withCredentials: false,
  // baseURL: API_BASE_URL, // APIのベースURL
  timeout: 5000, // タイムアウト設定 (ミリ秒)
  headers: {
    'Content-Type': 'application/json', // リクエストのContent-TypeをJSONに設定
  },
})

// リクエストインターセプター
axiosInstance.interceptors.request.use(
  (config) => {
    // ローカルストレージからトークンを取得して、ヘッダーに追加
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    // リクエストエラーが発生した場合
    return Promise.reject(error)
  },
)

// レスポンスインターセプター
axiosInstance.interceptors.response.use(
  (response) => {
    // 成功時の処理
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      const isHomePage = window.location.pathname === '/' // 現在のページがトップページかを確認
      // トークン切れでトップページでなければリダイレクト
      if (!isHomePage) {
        window.location.href = '/' // トップページにリダイレクト
      }
    }
    return Promise.reject(error)
  },
)

export default axiosInstance
