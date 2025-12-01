import { check, sleep } from "k6";
import http from "k6/http";
import {
  getToken,
  thresholds200p95,
  withAuth,
  isLambdaEnv,
  getLambdaThresholds,
} from "./helpers_c.js";

const baseUrl = __ENV.BASE_URL;
const label = __ENV.LABEL || "run";
const isLambda = isLambdaEnv(baseUrl);
const queries = (
  __ENV.SEARCH_SET || "able,test,go,ai,cat,run,play,have,make,good"
).split(",");

export const options = {
  setupTimeout: isLambda ? "120s" : "30s", // Lambda環境でのコールドスタートとDynamoDBタイムアウトを考慮
  scenarios: {
    set: {
      executor: "ramping-arrival-rate",
      startRate: 0,
      timeUnit: "1s",
      preAllocatedVUs: isLambda ? 1 : 3, // Lambda向け: ウォームアップ用に1 VUを確保
      maxVUs: 10,
      stages: isLambda
        ? [
            // Lambda向け: ウォームアップを追加
            { duration: "30s", target: 1 }, // コールドスタート対策: 1 req/sで30秒ウォームアップ
            { duration: "30s", target: 3 },
            { duration: "1m", target: 8 },
            { duration: "30s", target: 0 },
          ]
        : [
            // ローカル: 既存のステージ設定
            { duration: "30s", target: 3 },
            { duration: "1m", target: 8 },
            { duration: "30s", target: 0 },
          ],
      tags: {
        endpoint: "search",
        label,
        phase:
          label === "before_idx"
            ? "before_index"
            : label === "after_idx"
            ? "after_index"
            : "unknown",
      },
    },
  },
  thresholds: isLambda
    ? getLambdaThresholds("search") // Lambda向け: p95<1000ms
    : {
        ...thresholds200p95("search"), // ローカル: 共通SLO
      },
};

export function setup() {
  const token = getToken(baseUrl);
  return { token };
}

export default function (data) {
  // ランダムに検索クエリを選択
  const searchQuery = queries[Math.floor(Math.random() * queries.length)];

  // 注意: APIは ?search= パラメータを期待している（?q= ではない）
  const res = http.get(
    `${baseUrl}/words?search=${encodeURIComponent(searchQuery)}&sortBy=name`,
    withAuth(data.token, {
      endpoint: "search",
      label,
      phase:
        label === "before_idx"
          ? "before_index"
          : label === "after_idx"
          ? "after_index"
          : "unknown",
    })
  );
  const ok = check(res, { "2xx": (r) => r.status >= 200 && r.status < 300 });
  if (!ok) {
    console.error(
      `db_before_after NG: status=${res.status} body=${String(res.body).slice(
        0,
        200
      )}`
    );
  }
  sleep(0.2);
}
