import { check } from "k6";
import http from "k6/http";
import {
  think,
  isLambdaEnv,
  getLambdaStages,
  getLambdaThresholds,
  httpPostWithRetry,
} from "./helpers_b.js";

const profile = __ENV.PROFILE || "pr";
const baseUrl = __ENV.BASE_URL;
const email = __ENV.TEST_EMAIL || "k6-test@example.com";
const password = __ENV.TEST_PASSWORD || "k6-testPASS";

const isLambda = isLambdaEnv(baseUrl);

export const options = {
  scenarios: {
    sign_in: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: isLambda
        ? getLambdaStages(profile) // Lambda向け: ウォームアップ付き
        : profile === "nightly"
        ? [
            // ローカル: ポートフォリオ用途: 最大10 VU
            { duration: "20s", target: 3 }, // ウォームアップ: 3 VU
            { duration: "1m", target: 10 }, // ランプアップ: 10 VU
            { duration: "2m", target: 10 }, // ピーク: 2分間 10 VU維持
            { duration: "20s", target: 0 }, // ランプダウン
          ]
        : [
            // ローカル: PR: 最大5 VU
            { duration: "20s", target: 2 }, // ウォームアップ: 2 VU
            { duration: "1m", target: 5 }, // ランプアップ: 5 VU
            { duration: "1m30s", target: 5 }, // ピーク: 1分30秒 5 VU維持
            { duration: "20s", target: 0 }, // ランプダウン
          ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: isLambda
    ? getLambdaThresholds("sign_in") // Lambda向け: p95<1000ms
    : {
        // ローカル: 既存の閾値
        "http_req_failed{endpoint:sign_in}": ["rate<0.01"],
        "http_req_duration{endpoint:sign_in}": ["p(95)<200"],
      },
};

export function setup() {
  // サインインのためのユーザーを作っておく（既存なら409/422を許容）
  if (!email || !password) {
    throw new Error("TEST_EMAIL/TEST_PASSWORD are required for sign_in test");
  }
  const payload = JSON.stringify({ email, password, name: "k6-user" });

  const r = isLambda
    ? httpPostWithRetry(`${baseUrl}/users/sign_up`, payload, {
        headers: { "Content-Type": "application/json" },
        tags: { endpoint: "sign_up" },
      })
    : http.post(`${baseUrl}/users/sign_up`, payload, {
        headers: { "Content-Type": "application/json" },
        tags: { endpoint: "sign_up" },
        timeout: "10s", // ローカル環境でもタイムアウトを設定（10秒）
      });

  // ネットワークエラーの場合（status === 0）
  if (r.status === 0) {
    throw new Error(
      `sign_up failed: network error (server may not be running at ${baseUrl})`
    );
  }

  // レスポンスボディがnullの場合の処理
  const bodyStr = r.body ? String(r.body).slice(0, 200) : "(no body)";

  // エラーステータスの場合（409/422以外）
  if (r.status >= 400 && r.status !== 409 && r.status !== 422) {
    throw new Error(`sign_up failed: ${r.status} body=${bodyStr}`);
  }
}

export default function () {
  const res = http.post(
    `${baseUrl}/users/sign_in`,
    JSON.stringify({ email, password }),
    {
      headers: { "Content-Type": "application/json" },
      tags: { endpoint: "sign_in" }, // ← これだけを閾値判定
    }
  );
  const ok = check(res, {
    "2xx": (r) => r.status >= 200 && r.status < 300,
    "has token": (r) => {
      // レスポンスが成功した場合のみJSONパースを試みる
      if (r.status < 200 || r.status >= 300) {
        return false;
      }
      // レスポンスボディがnullの場合はfalseを返す
      if (!r.body) {
        return false;
      }
      try {
        return !!r.json("token");
      } catch (e) {
        // JSONパースエラーの場合はfalseを返す
        return false;
      }
    },
  });

  if (!ok) {
    const bodyStr = res.body ? String(res.body).slice(0, 300) : "(no body)";
    const errorMsg =
      res.status === 0
        ? "network error (server may not be running)"
        : `status=${res.status}`;
    console.error(`sign_in NG: ${errorMsg} body=${bodyStr}`);
  }
  think();
}
