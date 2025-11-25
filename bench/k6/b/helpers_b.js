import { fail, sleep } from 'k6'
import http from 'k6/http'

/**
 * test-login エンドポイントでトークンを取得
 * @param {string} baseUrl - ベースURL
 * @returns {string} JWTトークン
 */
export function getToken(baseUrl) {
  const res = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: { 'Content-Type': 'application/json' },
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
 * Authorizationヘッダー付きのリクエストオプションを返す
 * @param {string} token - JWTトークン
 * @returns {object} ヘッダーを含むオプションオブジェクト
 */
export function withAuth(token) {
  return {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
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
  // pr (default)
  return [
    { duration: '30s', target: 5 },
    { duration: '2m', target: 10 },
    { duration: '30s', target: 0 },
  ]
}

/**
 * 検索結果から先頭の単語IDを取得
 * @param {string} baseUrl - ベースURL
 * @param {string} token - JWTトークン
 * @param {string} q - 検索クエリ
 * @param {string} sortBy - ソート条件
 * @returns {number} 単語ID
 */
export function pickWordId(baseUrl, token, q, sortBy) {
  const res = http.get(
    `${baseUrl}/words?search=${encodeURIComponent(q)}&sortBy=${sortBy}&order=asc&page=1&limit=10`,
    withAuth(token),
  )
  if (res.status !== 200) {
    fail(`pickWordId list failed: ${res.status} ${String(res.body).slice(0, 200)}`)
  }
  const data = res.json()
  // 実装に応じて配列 or { items: [] } などに合わせて取得
  const item = data.words[0]
  if (!item || !item.id) {
    fail('no word id found in response' + JSON.stringify(data))
  }
  return item.id
}

/**
 * 共通の閾値設定を返す
 * @returns {object} 閾値設定オブジェクト
 */
export function commonThresholds() {
  return {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<200'],
  }
}

/**
 * ユーザーらしさのための think time（1秒）
 */
export function think() {
  sleep(1)
}
