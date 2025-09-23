import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import axios from '@/axiosConfig'

type LineCallbackResponse = {
  token?: string
  need_password?: boolean
  temp_token?: string
  suggested_mail?: string // サーバが出すなら利用
}

const LineCallback: React.FC = () => {
  const nav = useNavigate()
  const [err, setErr] = useState<string>('')

  useEffect(() => {
    ;(async () => {
      const qs = new URLSearchParams(window.location.search)
      const code = qs.get('code')
      const state = qs.get('state')

      if (!code || !state) {
        setErr('LINEの認証パラメータが不足しています。')
        return
      }

      try {
        const { data } = await axios.get<LineCallbackResponse>(
          '/users/auth/line/callback',
          {
            params: { code: qs.get('code'), state: qs.get('state') },
            withCredentials: true,
          },
        )

        // 既存ユーザー：即JWT
        if (data.token) {
          localStorage.setItem('token', data.token)
          nav('/mypage', { replace: true })
          return
        }

        // 未紐付け：仮登録（TempJWT）→ 完了ページへ
        if (data.need_password && data.temp_token) {
          sessionStorage.setItem('line_temp_token', data.temp_token)
          if (data.suggested_mail) {
            sessionStorage.setItem('line_suggested_mail', data.suggested_mail)
          }
          nav('/line/complete', { replace: true })
          return
        }

        // 予期しないレスポンス
        setErr('予期しない応答です。もう一度お試しください。')
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (e: any) {
        // 403(state mismatch) / 502(token_exchange_or_verify_failed) / 500 など
        setErr(
          e?.response?.data?.error || 'LINEログインでエラーが発生しました。',
        )
      }
    })()
  }, [nav])

  if (err) return <p className="text-red-600">{err}</p>
  return <p>LINE 認証処理中...</p>
}

export default LineCallback
