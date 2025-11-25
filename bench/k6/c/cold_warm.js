import { check, sleep } from 'k6'
import http from 'k6/http'
import { getToken, profileStages, thresholds200p95, withAuth } from './helpers_c.js'

const baseUrl = __ENV.BASE_URL
const q = __ENV.SEARCH_Q || 'test'
const sortBy = __ENV.SEARCH_SORT || 'name'

export const options = {
  scenarios: {
    cold_once: {
      executor: 'per-vu-iterations',
      vus: 1,
      iterations: 1,
      startTime: '0s',
      maxDuration: '2m',
      tags: { phase: 'cold', endpoint: 'search' },
    },
    warm_ramp: {
      executor: 'ramping-vus',
      startTime: '60s', // 手順で事前にアイドルを作る。ここは実行内の待機。
      stages: profileStages('pr'), // pr: 0→5→10→0 など
      tags: { phase: 'warm', endpoint: 'search' },
    },
  },
  thresholds: {
    ...thresholds200p95('search', 'warm'), // warm のみ評価
    // cold は記録のみ
  },
}

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  // シナリオのタグは自動的にリクエストに適用される
  // cold_onceシナリオは phase: 'cold' タグ、warm_rampシナリオは phase: 'warm' タグが設定済み
  // リクエストタグで明示的に設定（シナリオタグとマージされる）
  const res = http.get(
    `${baseUrl}/words?q=${encodeURIComponent(q)}&sortBy=${sortBy}`,
    withAuth(data.token, { endpoint: 'search' }),
  )
  const ok = check(res, { '2xx': (r) => r.status >= 200 && r.status < 300 })
  if (!ok) {
    console.error(`search NG: status=${res.status} body=${String(res.body).slice(0, 300)}`)
  }
  sleep(1)
}
