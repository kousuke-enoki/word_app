import { check, sleep } from "k6";
import http from "k6/http";
import { getToken, randomQuizParams, withAuth } from "./helpers_b.js";

const baseUrl = __ENV.BASE_URL;
const profile = __ENV.PROFILE || "pr";

export const options = {
  scenarios: {
    quiz_new: {
      executor: "ramping-vus",
      startVUs: 0,
      stages:
        profile === "nightly"
          ? [
              // ポートフォリオ用途: 最大10 VU（クイズ生成は重い処理のため控えめに）
              { duration: "20s", target: 3 },
              { duration: "1m", target: 10 },
              { duration: "2m", target: 10 },
              { duration: "20s", target: 0 },
            ]
          : [
              // PR: 最大5 VU
              { duration: "20s", target: 2 },
              { duration: "1m", target: 5 },
              { duration: "1m30s", target: 5 },
              { duration: "20s", target: 0 },
            ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: {
    // 429エラーはクォータ制限のため許容（テストユーザーの制約）
    "http_req_failed{endpoint:quiz_new,status:!429}": ["rate<0.01"],
    "http_req_duration{endpoint:quiz_new}": ["p(95)<200"],
  },
};

// setup()を削除 - 各VUごとにトークンを取得するため

export default function () {
  // 各iterationでトークンを取得（各VUが異なるテストユーザーを使用）
  const token = getToken(baseUrl);

  // ランダムなクイズパラメータを使用
  const body = randomQuizParams();

  const res = http.post(`${baseUrl}/quizzes/new`, JSON.stringify(body), {
    ...withAuth(token),
    tags: { endpoint: "quiz_new" },
  });

  const ok = check(res, {
    "2xx": (r) => r.status >= 200 && r.status < 300,
    // 429エラーはクォータ制限のため許容
    "2xx or 429": (r) =>
      (r.status >= 200 && r.status < 300) || r.status === 429,
  });

  if (!ok && res.status !== 429) {
    // 429以外のエラーのみログ出力
    console.error(
      `quiz_new NG: status=${res.status} body=${String(res.body).slice(0, 300)}`
    );
  }

  // クイズ生成は重い処理なので、少し長めのThinkTime
  sleep(5 + Math.random() * 5); // 5-10秒
}
