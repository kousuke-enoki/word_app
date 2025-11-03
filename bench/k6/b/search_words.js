import { check } from 'k6'
import http from 'k6/http'
import { getToken, think, withAuth } from './helpers_b.js'

const baseUrl = __ENV.BASE_URL
const q = __ENV.SEARCH_Q || 'test'
const sortBy = __ENV.SEARCH_SORT || 'name'

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  const res = http.get(
    `${baseUrl}/words?q=${encodeURIComponent(q)}&sortBy=${sortBy}`,
    withAuth(data.token),
  )
  const ok = check(res, { '2xx': (r) => r.status >= 200 && r.status < 300 })
  if (!ok) {
    console.error(`search_words NG: status=${res.status} body=${String(res.body).slice(0, 300)}`)
  }
  think()
}
