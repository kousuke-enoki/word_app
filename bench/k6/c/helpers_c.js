import { fail } from 'k6'
import http from 'k6/http'

/**
 * test-login エンドポイントでトークンを取得
 * @param {string} baseUrl - ベースURL
 * @returns {string} JWTトークン
 */
export function getToken(baseUrl) {
  // 環境に応じてタイムアウトを設定
  // Lambda環境の判定（isLambdaEnv()の前に定義されているため、直接判定ロジックを使用）
  const isLambda =
    __ENV.IS_LAMBDA === 'true' ||
    (baseUrl &&
      (baseUrl.includes('amazonaws.com') ||
        baseUrl.includes('lambda-url') ||
        baseUrl.includes('execute-api')))
  const timeout = isLambda ? '30s' : '10s' // Lambda: コールドスタート考慮、ローカル: 早めにエラー検出

  const res = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: { 'Content-Type': 'application/json' },
    tags: { step: 'test-login' },
    timeout,
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
 * プロファイルに応じたステージ設定を返す（ポートフォリオ用途に最適化）
 * @param {string} profile - 'pr' または 'nightly'
 * @returns {array} ステージ設定の配列
 */
export function profileStages(profile) {
  if (profile === 'nightly') {
    // ポートフォリオ用途: 最大10 VU, 10-15 req/s
    return [
      { duration: '30s', target: 3 },
      { duration: '2m', target: 10 },
      { duration: '30s', target: 0 },
    ]
  }
  // PR用: 最大5 VU, 5-10 req/s
  return [
    { duration: '20s', target: 2 },
    { duration: '1m30s', target: 5 },
    { duration: '20s', target: 0 },
  ]
}

/**
 * Lambda環境かどうかを判定
 * @param {string} baseUrl - ベースURL
 * @returns {boolean} Lambda環境の場合true
 */
export function isLambdaEnv(baseUrl) {
  // 環境変数で明示的に指定されている場合
  if (__ENV.IS_LAMBDA === 'true') return true
  if (__ENV.IS_LAMBDA === 'false') return false

  // BASE_URLから自動判定
  return (
    baseUrl &&
    (baseUrl.includes('amazonaws.com') ||
      baseUrl.includes('lambda-url') ||
      baseUrl.includes('execute-api'))
  )
}

/**
 * Lambda環境向けのステージ設定を返す（ウォームアップフェーズ付き）
 * @param {string} profile - 'pr' または 'nightly'
 * @returns {array} ステージ設定の配列
 */
export function getLambdaStages(profile) {
  if (profile === 'nightly') {
    // Lambda向け: ウォームアップを追加
    return [
      { duration: '30s', target: 1 }, // コールドスタート対策: 1 VUで30秒ウォームアップ
      { duration: '30s', target: 3 },
      { duration: '2m', target: 10 },
      { duration: '30s', target: 0 },
    ]
  }
  // PR プロファイル: Lambda向け: ウォームアップを追加
  return [
    { duration: '30s', target: 1 }, // コールドスタート対策: 1 VUで30秒ウォームアップ
    { duration: '20s', target: 2 },
    { duration: '1m30s', target: 5 },
    { duration: '20s', target: 0 },
  ]
}

/**
 * Lambda環境向けの閾値設定を返す
 * @param {string} endpoint - エンドポイント名（例: 'search'）
 * @param {string} phase - フェーズ名（オプション、例: 'warm'）
 * @param {object} options - 追加オプション
 * @param {boolean} options.exclude429 - 429エラーを除外するか（デフォルト: false）
 * @returns {object} 閾値設定オブジェクト
 */
export function getLambdaThresholds(endpoint, phase, options = {}) {
  const { exclude429 = false } = options

  // Lambda環境では閾値を緩和（コールドスタートとネットワークレイテンシを考慮）
  const p95Threshold = 1000 // Lambda: 1秒

  // phaseが指定されている場合はタグに含める
  const tag = phase ? `{endpoint:${endpoint},phase:${phase}}` : `{endpoint:${endpoint}}`

  const thresholds = {
    [`http_req_duration${tag}`]: [`p(95)<${p95Threshold}`],
  }

  // 429エラーを除外する場合
  if (exclude429) {
    thresholds[`http_req_failed${tag},status:!429`] = ['rate<0.01']
  } else {
    thresholds[`http_req_failed${tag}`] = ['rate<0.01']
  }

  return thresholds
}
