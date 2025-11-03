import { check } from 'k6'
import { isWriteAllowed, login, req } from './helpers.js'

// endpoints.json を読み込む
const endpoints = JSON.parse(open('./endpoints.json'))
export const options = {
  vus: 1,
  iterations: 1,
  thresholds: {
    http_req_failed: ['rate==0'],
  },
}

export function setup() {
  const baseUrl = __ENV.BASE_URL
  const email = __ENV.TEST_EMAIL || ''
  const password = __ENV.TEST_PASSWORD || ''
  const env = __ENV.SMOKE_ENV || 'local'

  const token = login(baseUrl, email, password, env)
  return { token }
}

export default function (data) {
  const baseUrl = __ENV.BASE_URL

  // 公開エンドポイントを順次実行
  for (const ep of endpoints.open) {
    const res = req(baseUrl, ep, null)
    const isSuccess = check(res, {
      [`${ep.method} ${ep.path} 2xx`]: (r) => r.status >= 200 && r.status < 300,
    })

    if (!isSuccess) {
      console.error(
        `OPEN NG: ${ep.method} ${ep.path} status=${res.status} body=${String(res.body).slice(
          0,
          200,
        )}`,
      )
    }
  }

  // 保護エンドポイントを順次実行
  for (const ep of endpoints.protected) {
    // write:true かつ ALLOW_WRITE=false の場合はスキップ
    if (ep.write && !isWriteAllowed()) {
      console.log(`SKIP(write): ${ep.method} ${ep.path}`)
      continue
    }

    const headers = { Authorization: `Bearer ${data.token}` }
    const res = req(baseUrl, ep, headers)
    const isSuccess = check(res, {
      [`${ep.method} ${ep.path} 2xx`]: (r) => r.status >= 200 && r.status < 300,
    })

    if (!isSuccess) {
      console.error(
        `PROT NG: ${ep.method} ${ep.path} status=${res.status} body=${String(res.body).slice(
          0,
          200,
        )}`,
      )
    }
  }
}
