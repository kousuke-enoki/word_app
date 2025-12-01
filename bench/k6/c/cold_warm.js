import { check, sleep } from "k6";
import http from "k6/http";
import {
  getToken,
  profileStages,
  thresholds200p95,
  withAuth,
  isLambdaEnv,
  getLambdaStages,
  getLambdaThresholds,
} from "./helpers_c.js";

const baseUrl = __ENV.BASE_URL;
const profile = __ENV.PROFILE || "pr";
const searchQuery = __ENV.SEARCH_Q || "test";
const sortBy = __ENV.SEARCH_SORT || "name";
const isLambda = isLambdaEnv(baseUrl);

export const options = {
  setupTimeout: isLambda ? "120s" : "30s", // Lambda環境でのコールドスタートとDynamoDBタイムアウトを考慮
  scenarios: {
    cold_once: {
      executor: "per-vu-iterations",
      vus: 1,
      iterations: 1,
      startTime: "0s",
      maxDuration: isLambda ? "3m" : "2m", // Lambda環境では長めに設定
      tags: { phase: "cold", endpoint: "search" },
    },
    warm_ramp: {
      executor: "ramping-vus",
      startTime: "60s", // 手順で事前にアイドルを作る。ここは実行内の待機。
      stages: isLambda
        ? getLambdaStages(profile) // Lambda向け: ウォームアップ付き
        : profileStages(profile), // ローカル: 既存のステージ設定
      tags: { phase: "warm", endpoint: "search" },
    },
  },
  thresholds: isLambda
    ? getLambdaThresholds("search", "warm") // Lambda向け: p95<1000ms
    : {
        ...thresholds200p95("search", "warm"), // ローカル: warm のみ評価
        // cold は記録のみ
      },
};

export function setup() {
  const token = getToken(baseUrl);
  return { token };
}

export default function (data) {
  // シナリオのタグは自動的にリクエストに適用される
  // cold_onceシナリオは phase: 'cold' タグ、warm_rampシナリオは phase: 'warm' タグが設定済み
  // リクエストタグで明示的に設定（シナリオタグとマージされる）
  // 注意: APIは ?search= パラメータを期待している（?q= ではない）
  const res = http.get(
    `${baseUrl}/words?search=${encodeURIComponent(
      searchQuery
    )}&sortBy=${sortBy}`,
    withAuth(data.token, { endpoint: "search" })
  );
  const ok = check(res, { "2xx": (r) => r.status >= 200 && r.status < 300 });
  if (!ok) {
    console.error(
      `search NG: status=${res.status} body=${String(res.body).slice(0, 300)}`
    );
  }
  sleep(1);
}
