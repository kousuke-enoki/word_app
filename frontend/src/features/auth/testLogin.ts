// src/features/auth/testLogin.ts
import axiosInstance from '@/axiosConfig'

// テストログインレスポンスの型定義
export interface TestLoginResponse {
  token: string
  user_id: number
  user_name: string
  jump: string
}

// LocalStorageに保存するデータの型
interface CachedTestLoginData {
  ts: number // タイムスタンプ（ミリ秒）
  payload: TestLoginResponse // レスポンス全体
}

export class TestLoginCooldownError extends Error {
  remainingMs: number
  cachedResponse?: TestLoginResponse // キャッシュされたレスポンス（再表示用）
  constructor(remainingMs: number, cachedResponse?: TestLoginResponse) {
    super('TEST_LOGIN_COOLDOWN')
    this.name = 'TestLoginCooldownError'
    this.remainingMs = remainingMs
    this.cachedResponse = cachedResponse
  }
}

// LocalStorageのキー
const STORAGE_KEY = 'test_login_cache'
const COOLDOWN_MS = 60_000 // 60秒

// キャッシュデータを取得
function getCachedData(): CachedTestLoginData | null {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (!stored) return null
    const data: CachedTestLoginData = JSON.parse(stored)
    return data
  } catch {
    return null
  }
}

// キャッシュデータを保存
function saveCachedData(response: TestLoginResponse): void {
  try {
    const data: CachedTestLoginData = {
      ts: Date.now(),
      payload: response,
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  } catch (error) {
    console.warn('Failed to save test login cache:', error)
  }
}

// キャッシュをクリア（ログアウト時など）
export function clearTestLoginCache(): void {
  try {
    localStorage.removeItem(STORAGE_KEY)
  } catch (error) {
    console.warn('Failed to clear test login cache:', error)
  }
}

// キャッシュが有効かチェック（60秒以内か）
function isCacheValid(cached: CachedTestLoginData): boolean {
  const now = Date.now()
  const elapsed = now - cached.ts
  return elapsed < COOLDOWN_MS
}

// ヘッダーからRetry-After値を取得（大文字小文字を考慮）
function getRetryAfter(
  headers: Record<string, string | undefined>,
): string | undefined {
  const headerKeys = Object.keys(headers)
  const retryAfterKey = headerKeys.find(
    (key) => key.toLowerCase() === 'retry-after',
  )
  if (retryAfterKey) {
    return headers[retryAfterKey]
  }
  // 直接アクセスも試す（axiosが正規化している場合）
  return headers['retry-after'] || headers['Retry-After']
}

// ヘッダーからX-Rate-Limit-Exceededフラグを取得
function isRateLimitExceeded(
  headers: Record<string, string | undefined>,
): boolean {
  const headerKeys = Object.keys(headers)
  const rateLimitExceededKey = headerKeys.find(
    (key) => key.toLowerCase() === 'x-rate-limit-exceeded',
  )
  return (
    rateLimitExceededKey !== undefined &&
    headers[rateLimitExceededKey] === 'true'
  )
}

// 進行中のリクエスト（重複呼び出し防止）
let inflight: Promise<TestLoginResponse> | null = null

/**
 * テストログインを実行
 *
 * 仕様:
 * - 60秒以内に成功レスポンスがあれば、APIを叩かずにキャッシュを返す
 * - キャッシュがある場合は、成功レスポンス全体を返す（再表示用）
 * - ログイン済み（tokenあり）の場合は、APIを叩かずに既存tokenを返す
 *
 * @returns 成功レスポンス全体（token, user_id, user_name, jump）
 * @throws TestLoginCooldownError 60秒以内でキャッシュがある場合（cachedResponseを含む）
 */
export async function testLogin(): Promise<TestLoginResponse> {
  // ログイン済みなら既存tokenからレスポンスを構築して返す
  const existingToken = localStorage.getItem('token')
  if (existingToken) {
    // 既存tokenがある場合は、そのtokenを使ってレスポンスを構築
    // ただし、完全な情報がないのでキャッシュがあればそれを使う
    const cached = getCachedData()
    if (
      cached &&
      isCacheValid(cached) &&
      cached.payload.token === existingToken
    ) {
      return cached.payload
    }
    // キャッシュがない場合は、最小限のレスポンスを返す
    return {
      token: existingToken,
      user_id: 0, // 不明
      user_name: '', // 不明
      jump: 'quiz', // デフォルト
    }
  }

  // キャッシュをチェック
  const cached = getCachedData()
  if (cached && isCacheValid(cached)) {
    // 60秒以内でキャッシュがある → APIを叩かずにキャッシュを返す
    // ただし、エラーとして扱う（cachedResponseを含む）ので、呼び出し側で処理できる
    const remaining = COOLDOWN_MS - (Date.now() - cached.ts)
    throw new TestLoginCooldownError(remaining, cached.payload)
  }

  // 進行中のリクエストがあれば、それを返す
  if (inflight) {
    return inflight
  }

  // API呼び出し
  inflight = (async () => {
    try {
      const res = await axiosInstance.post<TestLoginResponse>(
        '/users/auth/test-login',
      )

      // レート制限超過の判定（200 OKでもRetry-Afterヘッダーがあれば超過）
      const headers = (res.headers || {}) as Record<string, string | undefined>
      const retryAfter = getRetryAfter(headers)
      const rateLimitExceeded = isRateLimitExceeded(headers)
      const isRateLimited = !!retryAfter || rateLimitExceeded

      if (isRateLimited) {
        // レート制限超過として処理
        const retryAfterMs = retryAfter
          ? parseInt(retryAfter, 10) * 1000
          : COOLDOWN_MS

        // レスポンスデータはあるので、キャッシュとして保存
        const response = res.data
        if (response && response.token) {
          localStorage.setItem('token', response.token)
          saveCachedData(response)
        }

        // エラーとして扱う（キャッシュされたレスポンスを含む）
        const cached = getCachedData()
        throw new TestLoginCooldownError(retryAfterMs, cached?.payload)
      }

      const response = res.data
      if (!response || !response.token) {
        throw new Error('Invalid response: missing token')
      }

      // tokenをLocalStorageに保存（既存の動作を維持）
      localStorage.setItem('token', response.token)

      // レスポンス全体をキャッシュに保存
      saveCachedData(response)

      return response
    } catch (error: unknown) {
      // 429エラー（レート制限超過）の処理
      if (
        error &&
        typeof error === 'object' &&
        'response' in error &&
        error.response &&
        typeof error.response === 'object' &&
        'status' in error.response &&
        error.response.status === 429
      ) {
        const axiosError = error as {
          response: {
            status: number
            headers?: Record<string, string | undefined> | Headers
          }
        }
        const headers = (axiosError.response.headers || {}) as Record<
          string,
          string | undefined
        >
        const retryAfter = getRetryAfter(headers)
        const retryAfterMs = retryAfter
          ? parseInt(retryAfter, 10) * 1000
          : COOLDOWN_MS

        // キャッシュがあればそれを使う
        const cached = getCachedData()
        if (cached && isCacheValid(cached)) {
          throw new TestLoginCooldownError(retryAfterMs, cached.payload)
        }

        // キャッシュがない場合は、Retry-Afterの時間だけ待つ必要がある
        throw new TestLoginCooldownError(retryAfterMs)
      }

      // その他のエラーはそのまま再スロー
      throw error
    } finally {
      inflight = null
    }
  })()

  return inflight
}

/**
 * キャッシュからレスポンスを取得（APIを叩かない）
 * 60秒以内のキャッシュがあれば返す。なければnull
 */
export function getCachedTestLoginResponse(): TestLoginResponse | null {
  const cached = getCachedData()
  if (cached && isCacheValid(cached)) {
    return cached.payload
  }
  return null
}
