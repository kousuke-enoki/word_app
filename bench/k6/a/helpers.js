import { fail } from 'k6'
import http from 'k6/http'

// /**
//  * endpoints.json を読み込む
//  * @returns {object} エンドポイント定義オブジェクト
//  */
// export function loadEndpoints() {
//   const content = open('./endpoints.json') // smoke_all.js と同じフォルダ前提
//   return JSON.parse(content)
// }

/**
 * ログインしてトークンを取得
 * @param {string} baseUrl - ベースURL
 * @param {string} email - メールアドレス（sign_in用、test-loginでは未使用）
 * @param {string} password - パスワード（sign_in用、test-loginでは未使用）
 * @param {string} env - 環境（'aws' または 'local'）
 * @returns {string} JWTトークン
 */
export function login(baseUrl, email, password, env) {
  let url
  let payload

  // test-login を優先（デフォルト）
  if (!env || env === 'local' || env === 'aws') {
    url = `${baseUrl}/users/auth/test-login`
    payload = JSON.stringify({})
  } else {
    // フォールバック: sign_in
    url = `${baseUrl}/users/sign_in`
    payload = JSON.stringify({
      email: email,
      password: password,
    })
  }

  const res = http.post(url, payload, {
    headers: { 'Content-Type': 'application/json' },
  })

  if (res.status !== 200) {
    fail(`Login failed: ${res.status} - ${String(res.body).slice(0, 200)}`)
  }

  const body = JSON.parse(res.body)
  if (!body.token) {
    fail(`No token in response: ${String(res.body).slice(0, 200)}`)
  }

  return body.token
}

/**
 * エンドポイントに対してリクエストを送信
 * @param {string} baseUrl - ベースURL
 * @param {object} ep - エンドポイント定義 { method, path, body? }
 * @param {object} headers - 追加ヘッダー（Authorization など）
 * @returns {http.Response} HTTPレスポンス
 */
export function req(baseUrl, ep, headers) {
  const url = `${baseUrl}${ep.path}`
  const defaultHeaders = {
    'Content-Type': 'application/json',
  }

  const finalHeaders = { ...defaultHeaders, ...headers }

  let res
  const method = ep.method.toUpperCase()

  if (ep.body) {
    const payload = JSON.stringify(ep.body)
    switch (method) {
      case 'POST':
        res = http.post(url, payload, { headers: finalHeaders })
        break
      case 'PUT':
        res = http.put(url, payload, { headers: finalHeaders })
        break
      case 'PATCH':
        res = http.patch(url, payload, { headers: finalHeaders })
        break
      default:
        res = http.post(url, payload, { headers: finalHeaders })
    }
  } else {
    switch (method) {
      case 'GET':
        res = http.get(url, { headers: finalHeaders })
        break
      case 'DELETE':
        res = http.del(url, null, { headers: finalHeaders })
        break
      default:
        res = http.get(url, { headers: finalHeaders })
    }
  }

  return res
}

/**
 * 書き込み操作が許可されているかどうか
 * @returns {boolean}
 */
export function isWriteAllowed() {
  return __ENV.ALLOW_WRITE === 'true'
}
