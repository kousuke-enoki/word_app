// src/features/auth/testLogin.ts
import axiosInstance from '@/axiosConfig'

export class TestLoginCooldownError extends Error {
  remainingMs: number
  constructor(remainingMs: number) {
    super('TEST_LOGIN_COOLDOWN')
    this.name = 'TestLoginCooldownError'
    this.remainingMs = remainingMs
  }
}

// 1分以内の再クリックは同じ結果を返す（/auth/test-login のサーバ側仕様に合わせる）
// クリック連打防止のため
let inflight: Promise<string> | null = null
let lastOkAt = 0
let lastToken: string | null = null
const COOLDOWN_MS = 60_000

export async function testLogin(): Promise<string> {
  const now = Date.now()
  const hasLocalToken = !!localStorage.getItem('token')

  // 直近成功から1分未満 & いまは未ログイン（＝再テストログインを試みている）
  if (lastToken && !hasLocalToken && now - lastOkAt < COOLDOWN_MS) {
    const remaining = COOLDOWN_MS - (now - lastOkAt)
    throw new TestLoginCooldownError(remaining)
  }
  if (inflight) return inflight

  // ログイン済みならAPIを叩かず即返す
  if (hasLocalToken) {
    return localStorage.getItem('token') as string
  }

  inflight = (async () => {
    try {
      const res = await axiosInstance.post('/users/auth/test-login')
      const token = res.data?.token as string
      if (token) {
        localStorage.setItem('token', token)
        lastOkAt = Date.now()
        lastToken = token
      }
      return token
    } finally {
      inflight = null
    }
  })()

  return inflight
}
