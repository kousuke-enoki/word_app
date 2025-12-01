import { check, sleep } from 'k6'
import http from 'k6/http'
import { Counter, Trend } from 'k6/metrics'
import { isLambdaEnv } from './helpers_c.js'

const baseUrl = __ENV.BASE_URL
const isLambda = isLambdaEnv(baseUrl)

// --- カスタムメトリクス ---
export const rate_limit_429 = new Counter('rate_limit_429') // 429 件数
export const rate_limit_exceeded_header = new Counter('rate_limit_exceeded_header') // X-Rate-Limit-Exceeded ヘッダー件数
export const retry_after_seconds = new Trend('retry_after_seconds') // Retry-After 秒数

// User-Agentを統一して、レート制限のキー（IP + UAHash + Route）を同じにする
// これにより、複数のVUが同じレート制限キーで判定される
const UNIFIED_USER_AGENT = 'k6-rate-limit-test/1.0'
// IPアドレスも統一して、レート制限キーを同じにする
const UNIFIED_IP = '192.168.1.100' // テスト用の固定IPアドレス

export const options = {
  scenarios: {
    // スパイクフェーズ: キャッシュなしで連続test-login → 429を確実に検出
    spike_no_cache: {
      executor: 'constant-arrival-rate',
      rate: 10,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 20, // 増加: 10 → 20
      maxVUs: 50, // 増加: 10 → 50（余裕を持たせる）
      tags: { phase: 'spike', cache: 'no', endpoint: 'test-login' },
    },
    // 回復フェーズ: test-login → test-logout → test-login（実運用を再現）
    recovery_with_cache: {
      executor: 'constant-arrival-rate',
      startTime: '45s',
      rate: 5,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 10, // 増加: 5 → 10
      maxVUs: 20, // 増加: 10 → 20
      tags: { phase: 'recovery', cache: 'yes', endpoint: 'test-login' },
    },
  },
  thresholds: {
    // スパイクフェーズ: 429エラーが大量に発生すること（キャッシュなし）
    // User-Agentを統一した場合、レート制限は1分間に1リクエスト（RATE_LIMIT_MAX_REQUESTS=1）
    // 10 req/s × 30秒 = 約300リクエスト中、最初の1件のみ成功、残り299件は429
    // 429率は約99.7%になるはず（実際は最初の1件が成功するため、299/300 ≈ 99.7%）
    // ただし、複数のVUが同時にリクエストを送信する場合、最初の数件が成功する可能性がある
    // そのため、より現実的な期待値として、90%以上が429になることを期待する
    'http_req_failed{phase:spike,endpoint:test-login}': ['rate>0.9'], // 90%以上が429
    'rate_limit_429{phase:spike}': ['count>250'], // 250件以上の429（300件中250件以上）

    // 回復フェーズ: 正常に動作すること
    'http_req_failed{phase:recovery,endpoint:test-login}': ['rate<0.01'],
    // Lambda環境ではレスポンス時間の閾値を緩和
    [`http_req_duration{phase:recovery,endpoint:test-login}`]: [
      isLambda ? 'p(95)<1000' : 'p(95)<200',
    ],
    'rate_limit_429{phase:recovery}': ['count<=0'], // 回復時は429なし

    // X-Rate-Limit-Exceededヘッダーが付くこと（キャッシュありでレート制限超過時）
    // 回復フェーズでは、最初のtest-loginでキャッシュが保存され、
    // 2回目のtest-login（同じウィンドウ内）でキャッシュが返される
    // 5 req/s × 30秒 = 約150リクエスト（各リクエストでlogin1とlogin2の2回のtest-login）
    // 各iterationでlogin2がキャッシュを返すため、約150件のX-Rate-Limit-Exceededヘッダーが期待される
    // ただし、ウィンドウのタイミングによっては、一部のリクエストでキャッシュが失効している可能性がある
    // そのため、より現実的な期待値として、50件以上を期待する
    'rate_limit_exceeded_header{phase:recovery}': ['count>50'],
  },
}

// test-loginは認証不要なので、setup()は不要

// スパイクフェーズ: キャッシュなしで連続test-login → 429を確実に検出
function spikePhase() {
  // 最初の1件だけ成功し、残りは429になる（MaxRequests=1, WindowSeconds=60のため）
  // User-AgentとIPを統一して、レート制限のキーを同じにする
  const res = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': UNIFIED_USER_AGENT,
      'X-Forwarded-For': UNIFIED_IP, // IPアドレスを統一
      'X-Real-IP': UNIFIED_IP, // 追加: X-Real-IPも設定
    },
    tags: { endpoint: 'test-login' },
  })

  // デバッグ用: 最初の10件だけログ出力（サンプリング）
  if (!globalThis._spikeLogCount) {
    globalThis._spikeLogCount = 0
  }
  if (globalThis._spikeLogCount < 10) {
    console.log(
      `[SPIKE] status=${res.status}, Retry-After=${res.headers['Retry-After']}, X-Rate-Limit-Exceeded=${res.headers['X-Rate-Limit-Exceeded']}`,
    )
    globalThis._spikeLogCount++
  }

  const ok = check(res, {
    '2xx/429': (r) => (r.status >= 200 && r.status < 300) || r.status === 429,
  })
  if (!ok) {
    console.error(
      `rate_limit spike NG: status=${res.status} body=${String(res.body).slice(0, 200)}`,
    )
  }

  // 429 をカウント（フェーズタグはシナリオから引き継がれる）
  if (res.status === 429) {
    rate_limit_429.add(1)

    // Retry-After が数値ならTrendに入れる
    const ra = res.headers['Retry-After']
    if (ra && !Number.isNaN(Number(ra))) {
      retry_after_seconds.add(Number(ra))
    }
  }

  sleep(0.02) // executorが到着率を制御するため、最小限のスリープ
}

// 回復フェーズ: test-login → test-login（実運用を再現、キャッシュありでX-Rate-Limit-Exceededを検証）
function recoveryPhase() {
  // 1. test-login（キャッシュ保存）
  // User-AgentとIPを統一して、スパイクフェーズと同じレート制限キーを使用する
  const res1 = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': UNIFIED_USER_AGENT,
      'X-Forwarded-For': UNIFIED_IP, // IPアドレスを統一
      'X-Real-IP': UNIFIED_IP, // 追加: X-Real-IPも設定
    },
    tags: { endpoint: 'test-login', step: 'login1' },
  })

  if (res1.status !== 200) {
    console.error(
      `rate_limit recovery login1 NG: status=${res1.status} body=${String(res1.body).slice(
        0,
        200,
      )}`,
    )
    sleep(4) // 5 req/s に合わせて調整
    return
  }

  const token = res1.json('token')
  if (!token) {
    console.error('rate_limit recovery: no token in response')
    sleep(4)
    return
  }

  // 少し待機（レート制限ウィンドウ内で2回目のtest-loginを試みる）
  sleep(0.5)

  // 2. 2回目のtest-login（1分以内なのでレート制限がかかる）
  // キャッシュが残っている場合、200 OK + X-Rate-Limit-Exceeded が返る
  // キャッシュがない場合、429が返る
  // User-AgentとIPを統一して、login1と同じレート制限キーを使用する
  const res2 = http.post(`${baseUrl}/users/auth/test-login`, JSON.stringify({}), {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': UNIFIED_USER_AGENT,
      'X-Forwarded-For': UNIFIED_IP, // IPアドレスを統一
      'X-Real-IP': UNIFIED_IP, // 追加: X-Real-IPも設定
    },
    tags: { endpoint: 'test-login', step: 'login2' },
  })

  const login2Ok = check(res2, {
    '2xx/429': (r) => (r.status >= 200 && r.status < 300) || r.status === 429,
  })
  if (!login2Ok) {
    console.error(
      `rate_limit recovery login2 NG: status=${res2.status} body=${String(res2.body).slice(
        0,
        200,
      )}`,
    )
  }

  // 429 をカウント
  if (res2.status === 429) {
    rate_limit_429.add(1)
    const ra = res2.headers['Retry-After']
    if (ra && !Number.isNaN(Number(ra))) {
      retry_after_seconds.add(Number(ra))
    }
  }

  // X-Rate-Limit-Exceeded ヘッダーをチェック（キャッシュありでレート制限超過時）
  // キャッシュがある場合は200 OK + X-Rate-Limit-Exceededが返る
  if (res2.status === 200 && res2.headers['X-Rate-Limit-Exceeded'] === 'true') {
    rate_limit_exceeded_header.add(1)

    // Retry-After もチェック
    const ra = res2.headers['Retry-After']
    if (ra && !Number.isNaN(Number(ra))) {
      retry_after_seconds.add(Number(ra))
    }
  }

  // デバッグ用: 最初の10件だけログ出力（サンプリング）
  if (!globalThis._recoveryLogCount) {
    globalThis._recoveryLogCount = 0
  }
  if (globalThis._recoveryLogCount < 10) {
    console.log(
      `[RECOVERY] login2: status=${res2.status}, Retry-After=${res2.headers['Retry-After']}, X-Rate-Limit-Exceeded=${res2.headers['X-Rate-Limit-Exceeded']}`,
    )
    globalThis._recoveryLogCount++
  }

  // test-logoutでキャッシュ削除（実運用を再現）
  const res3 = http.post(`${baseUrl}/users/auth/test-logout`, null, {
    headers: { Authorization: `Bearer ${token}` },
    tags: { endpoint: 'test-logout', step: 'logout' },
  })

  const logoutOk = check(res3, {
    204: (r) => r.status === 204,
  })
  if (!logoutOk) {
    console.error(
      `rate_limit recovery logout NG: status=${res3.status} body=${String(res3.body).slice(
        0,
        200,
      )}`,
    )
  }

  sleep(3.5) // 合計4秒（5 req/s に合わせて調整、実際の到着率はexecutorが制御）
}

export default function () {
  // k6では各シナリオが独立したVUプールで実行される
  // spike_no_cache: startTimeなし（0秒から開始、30秒間）
  // recovery_with_cache: startTime: '45s'（45秒から開始、30秒間）

  // 実行開始からの経過時間を取得
  // k6の実行開始時刻は環境変数やグローバル変数では取得できないため、
  // 最初の実行時に記録する
  if (!globalThis._k6StartTime) {
    globalThis._k6StartTime = Date.now() / 1000
  }

  const elapsed = Date.now() / 1000 - globalThis._k6StartTime

  // recovery_with_cacheは45秒後に開始されるため、
  // 45秒以降の実行は回復フェーズ、それ以前はスパイクフェーズ
  if (elapsed >= 45) {
    // 回復フェーズ: test-login → test-logout → test-login
    recoveryPhase()
  } else {
    // スパイクフェーズ: キャッシュなしで連続test-login
    spikePhase()
  }
}
