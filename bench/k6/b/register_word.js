import { check } from "k6";
import http from "k6/http";
import {
  getToken,
  pickWordId,
  randomSearchQuery,
  randomSortBy,
  resetRegisteredWords,
  think,
  withAuth,
  isLambdaEnv,
  getLambdaStages,
  getLambdaThresholds,
} from "./helpers_b.js";

const baseUrl = __ENV.BASE_URL;
const profile = __ENV.PROFILE || "pr";
const isLambda = isLambdaEnv(baseUrl);

export const options = {
  setupTimeout: isLambda ? "120s" : "30s", // Lambda環境でのコールドスタートとDynamoDBタイムアウトを考慮して120秒に設定
  scenarios: {
    register_word: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: isLambda
        ? getLambdaStages(profile) // Lambda向け: ウォームアップ付き
        : profile === "nightly"
        ? [
            // ローカル: ポートフォリオ用途: 最大10 VU
            { duration: "20s", target: 3 },
            { duration: "1m", target: 10 },
            { duration: "2m", target: 10 },
            { duration: "20s", target: 0 },
          ]
        : [
            // ローカル: PR: 最大5 VU
            { duration: "20s", target: 2 },
            { duration: "1m", target: 5 },
            { duration: "1m30s", target: 5 },
            { duration: "20s", target: 0 },
          ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: isLambda
    ? getLambdaThresholds("register_word", { exclude429: true }) // Lambda向け: p95<1000ms, 429除外
    : {
        // ローカル: 既存の閾値
        // 429エラーを許容するため、429以外のエラーのみをチェック
        "http_req_failed{endpoint:register_word,status:!429}": ["rate<0.01"],
        "http_req_duration{endpoint:register_word}": ["p(95)<200"],
      },
};

export function setup() {
  const token = getToken(baseUrl);

  // テスト実行前に登録済み単語を全件リセット（上限200件まで）
  resetRegisteredWords(baseUrl, token, 200);

  return { token };
}

export default function (data) {
  // ランダムな検索クエリとソート条件で単語IDを取得
  const q = randomSearchQuery();
  const sortBy = randomSortBy();
  const id = __ENV.REGISTER_WORD_ID
    ? Number(__ENV.REGISTER_WORD_ID)
    : pickWordId(baseUrl, data.token, q, sortBy);

  const res = http.post(
    `${baseUrl}/words/register`,
    JSON.stringify({ wordId: id, isRegistered: true }),
    {
      ...withAuth(data.token),
      tags: { endpoint: "register_word" },
    }
  );

  const ok = check(res, {
    "2xx or 409 or acceptable errors": (r) => {
      // 2xx: 成功
      if (r.status >= 200 && r.status < 300) return true;
      // 409: ユニーク制約違反（許容）
      if (r.status === 409) return true;
      // 400: 登録状態に変更なし（許容）
      if (r.status === 400) {
        const bodyStr = String(r.body);
        if (bodyStr.includes("no change in registration state")) {
          return true;
        }
      }
      // 429: 登録単語数上限超過（許容）
      if (r.status === 429) {
        const bodyStr = String(r.body);
        if (bodyStr.includes("registered words limit exceeded")) {
          return true;
        }
      }
      return false;
    },
  });

  if (!ok) {
    console.error(
      `register_word NG: status=${res.status} body=${String(res.body).slice(
        0,
        300
      )} id=${id}`
    );
  }

  think();
}
