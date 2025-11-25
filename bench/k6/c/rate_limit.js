import { check, sleep } from 'k6'
import http from 'k6/http'
import { Counter, Trend } from 'k6/metrics'
import { getToken, withAuth } from './helpers_c.js'

const baseUrl = __ENV.BASE_URL
const q = __ENV.SEARCH_Q || 'test'

// --- カスタムメトリクス ---
export const rate_limit_429 = new Counter('rate_limit_429') // 429 件数
export const retry_after_seconds = new Trend('retry_after_seconds') // Retry-After 秒数（任意）

export const options = {
  scenarios: {
    spike: {
      executor: 'constant-arrival-rate',
      rate: 50,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 50,
      tags: { phase: 'spike', endpoint: 'search' },
    },
    recovery: {
      executor: 'constant-arrival-rate',
      startTime: '45s',
      rate: 5,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 5,
      tags: { phase: 'recovery', endpoint: 'search' },
    },
  },
  thresholds: {
    // 回復フェーズは SLO を満たすこと
    'http_req_failed{endpoint:search,phase:recovery}': ['rate<0.01'],
    'http_req_duration{endpoint:search,phase:recovery}': ['p(95)<200'],

    // スパイクで 429 が1件以上、回復では 0 件
    'rate_limit_429{phase:spike}': ['count>0'],
    'rate_limit_429{phase:recovery}': ['count<=0'],
  },
}

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  // シナリオのタグ（phase）が自動付与される。ここでも endpoint を付与しておく。
  const res = http.get(
    `${baseUrl}/words?q=${encodeURIComponent(q)}&sortBy=name`,
    withAuth(data.token, { endpoint: 'search' }),
  )

  const ok = check(res, {
    '2xx/429': (r) => (r.status >= 200 && r.status < 300) || r.status === 429,
  })
  if (!ok) {
    console.error(`rate_limit NG: status=${res.status} body=${String(res.body).slice(0, 200)}`)
  }

  // 429 をカウント（フェーズタグはシナリオから引き継がれる）
  if (res.status === 429) {
    rate_limit_429.add(1) // ← phase:spike|recovery タグが自動で紐づく

    // Retry-After が数値ならTrendに入れる（秒）。数値でなければスキップ。
    const ra = res.headers['Retry-After']
    if (ra && !Number.isNaN(Number(ra))) {
      retry_after_seconds.add(Number(ra))
    }
  }

  sleep(0.2)
}
