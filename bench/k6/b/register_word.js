import { check } from 'k6'
import http from 'k6/http'
import { getToken, pickWordId, think, withAuth } from './helpers_b.js'

const baseUrl = __ENV.BASE_URL
const q = __ENV.SEARCH_Q || 'able'
const sortBy = __ENV.SEARCH_SORT || 'name'

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  const id = __ENV.REGISTER_WORD_ID
    ? Number(__ENV.REGISTER_WORD_ID)
    : pickWordId(baseUrl, data.token, q, sortBy)

  const res = http.post(
    `${baseUrl}/words/register`,
    JSON.stringify({ wordId: id, isRegistered: true }),
    withAuth(data.token),
  )
  const ok = check(res, {
    '2xx or 409': (r) => (r.status >= 200 && r.status < 300) || r.status === 409,
  })
  if (!ok) {
    console.error(
      `register_word NG: status=${res.status} body=${String(res.body).slice(0, 300)} id=${id}`,
    )
  }
  think()
}
