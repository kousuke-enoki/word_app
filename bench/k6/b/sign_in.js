import { check } from "k6";
import http from "k6/http";
import { think } from "./helpers_b.js";

const profile = __ENV.PROFILE || "pr";
const baseUrl = __ENV.BASE_URL;
const email = __ENV.TEST_EMAIL || "k6-test@example.com";
const password = __ENV.TEST_PASSWORD || "k6-testPASS";

export const options = {
  scenarios: {
    sign_in: {
      executor: "ramping-vus",
      startVUs: 0,
      stages:
        profile === "nightly"
          ? [
              // ポートフォリオ用途: 最大10 VU
              { duration: "20s", target: 3 }, // ウォームアップ: 3 VU
              { duration: "1m", target: 10 }, // ランプアップ: 10 VU
              { duration: "2m", target: 10 }, // ピーク: 2分間 10 VU維持
              { duration: "20s", target: 0 }, // ランプダウン
            ]
          : [
              // PR: 最大5 VU
              { duration: "20s", target: 2 }, // ウォームアップ: 2 VU
              { duration: "1m", target: 5 }, // ランプアップ: 5 VU
              { duration: "1m30s", target: 5 }, // ピーク: 1分30秒 5 VU維持
              { duration: "20s", target: 0 }, // ランプダウン
            ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: {
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
  const r = http.post(`${baseUrl}/users/sign_up`, payload, {
    headers: { "Content-Type": "application/json" },
    tags: { endpoint: "sign_up" },
  });
  if (r.status >= 400 && r.status !== 409 && r.status !== 422) {
    throw new Error(
      `sign_up failed: ${r.status} body=${String(r.body).slice(0, 200)}`
    );
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
    "has token": (r) => !!r.json("token"),
  });
  if (!ok) {
    console.error(
      `sign_in NG: status=${res.status} body=${String(res.body).slice(0, 300)}`
    );
  }
  think();
}
