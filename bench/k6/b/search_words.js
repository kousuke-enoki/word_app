import { check } from "k6";
import http from "k6/http";
import {
  getToken,
  randomSearchQuery,
  randomSortBy,
  think,
  withAuth,
  isLambdaEnv,
  getLambdaThresholds,
} from "./helpers_b.js";

const baseUrl = __ENV.BASE_URL;
const profile = __ENV.PROFILE || "pr";
const isLambda = isLambdaEnv(baseUrl);

export const options = {
  setupTimeout: isLambda ? "120s" : "30s", // Lambda環境でのコールドスタートとDynamoDBタイムアウトを考慮して120秒に設定
  scenarios: {
    search_words: {
      executor: "constant-arrival-rate",
      // ポートフォリオ用途: 同時5-10ユーザー想定
      rate: profile === "nightly" ? 10 : 5, // PR: 5 req/s, Nightly: 10 req/s
      timeUnit: "1s",
      duration: profile === "nightly" ? "3m" : "2m", // PR: 2分, Nightly: 3分
      preAllocatedVUs: isLambda ? 1 : 3, // Lambda向け: ウォームアップ用に1 VUを確保
      maxVUs: profile === "nightly" ? 10 : 5,
      gracefulStop: "30s",
    },
  },
  thresholds: isLambda
    ? getLambdaThresholds("search") // Lambda向け: p95<1000ms
    : {
        // ローカル: 既存の閾値
        "http_req_failed{endpoint:search}": ["rate<0.01"],
        "http_req_duration{endpoint:search}": ["p(95)<200"], // ポートフォリオ用途に適した閾値
      },
};

export function setup() {
  const token = getToken(baseUrl);
  return { token };
}

export default function (data) {
  // ランダムな検索クエリとソート条件を使用
  const q = randomSearchQuery();
  const sortBy = randomSortBy();

  const res = http.get(
    `${baseUrl}/words?search=${encodeURIComponent(q)}&sortBy=${sortBy}`,
    {
      ...withAuth(data.token),
      tags: { endpoint: "search" },
    }
  );

  const ok = check(res, {
    "2xx": (r) => r.status >= 200 && r.status < 300,
  });

  if (!ok) {
    console.error(
      `search_words NG: status=${res.status} body=${String(res.body).slice(
        0,
        300
      )}`
    );
  }

  think();
}
