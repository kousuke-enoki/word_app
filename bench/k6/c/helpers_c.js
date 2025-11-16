import { fail } from 'k6'
import http from 'k6/http'

/**
 * test-login エンドポイントでトークンを取得
 * @param {string} baseUrl - ベースURL
 * @returns {string} JWTトークン
 */
export function getToken(baseUrl) {
  const res = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: { 'Content-Type': 'application/json' },
    tags: { step: 'test-login' },
  })
  if (res.status !== 200) {
    fail(`test-login failed: ${res.status} ${String(res.body).slice(0, 200)}`)
  }
  const token = res.json('token')
  if (!token) {
    fail('no token in test-login response')
  }
  return token
}

/**
 * Authorizationヘッダー付きのリクエストオプションを返す（タグ付き）
 * @param {string} token - JWTトークン
 * @param {object} extraTags - 追加タグ（例: { endpoint: 'search', phase: 'cold' }）
 * @returns {object} ヘッダーとタグを含むオプションオブジェクト
 */
export function withAuth(token, extraTags = {}) {
  return {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    tags: extraTags,
  }
}

/**
 * エンドポイントごとのタグ付き閾値設定を返す
 * @param {string} endpoint - エンドポイント名（例: 'search'）
 * @param {string} phase - フェーズ名（オプション、例: 'warm'）
 * @returns {object} 閾値設定オブジェクト
 */
export function thresholds200p95(endpoint, phase) {
  const tag = phase ? `{endpoint:${endpoint},phase:${phase}}` : `{endpoint:${endpoint}}`
  return {
    [`http_req_failed${tag}`]: ['rate<0.01'],
    [`http_req_duration${tag}`]: ['p(95)<200'],
  }
}

/**
 * プロファイルに応じたステージ設定を返す
 * @param {string} profile - 'pr' または 'nightly'
 * @returns {array} ステージ設定の配列
 */
export function profileStages(profile) {
  if (profile === 'nightly') {
    return [
      { duration: '45s', target: 10 },
      { duration: '4m', target: 40 },
      { duration: '45s', target: 0 },
    ]
  }
  return [
    { duration: '30s', target: 5 },
    { duration: '2m', target: 10 },
    { duration: '30s', target: 0 },
  ]
}
