import { check } from 'k6'
import http from 'k6/http'
import { getToken, think, withAuth } from './helpers_b.js'

const baseUrl = __ENV.BASE_URL

export function setup() {
  const token = getToken(baseUrl)
  return { token }
}

export default function (data) {
  const body = {
    questionCount: 10,
    isSaveResult: false,
    isRegisteredWords: 0,
    correctRate: 1,
    attentionLevelList: [],
    partsOfSpeeches: [1],
    isIdioms: 0,
    isSpecialCharacters: 0,
  }
  const res = http.post(`${baseUrl}/quizzes/new`, JSON.stringify(body), withAuth(data.token))
  const ok = check(res, { '2xx': (r) => r.status >= 200 && r.status < 300 })
  if (!ok) {
    console.error(`quiz_new NG: status=${res.status} body=${String(res.body).slice(0, 300)}`)
  }
  think()
}
