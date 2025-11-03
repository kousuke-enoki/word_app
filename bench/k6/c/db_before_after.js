import { check, sleep } from 'k6'
import http from 'k6/http'
import { getToken, thresholds200p95, withAuth } from './helpers_c.js'

const baseUrl = __ENV.BASE_URL
const label = __ENV.LABEL || 'run'
const queries = (__ENV.SEARCH_SET || 'able,test,go,ai,cat,run,play,have,make,good').split(',')

export const options = {
  scenarios: {
    set: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 20,
      maxVUs: 100,
      stages: [
        { duration: '30s', target: 10 },
        { duration: '1m30s', target: 30 },
        { duration: '30s', target: 0 },
      ],
      tags: { endpoint: 'search', label },
    },
  },
  thresholds: {
    ...thresholds200p95('search'), // 共通SLO
  },
}

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  const q = queries[Math.floor(Math.random() * queries.length)]
  const res = http.get(
    `${baseUrl}/words?q=${encodeURIComponent(q)}&sortBy=name`,
    withAuth(data.token, { endpoint: 'search', label }),
  )
  const ok = check(res, { '2xx': (r) => r.status >= 200 && r.status < 300 })
  if (!ok) {
    console.error(`db_before_after NG: status=${res.status} body=${String(res.body).slice(0, 200)}`)
  }
  sleep(0.2)
}
