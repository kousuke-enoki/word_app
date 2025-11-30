import { fail, sleep } from "k6";
import http from "k6/http";

/**
 * test-login エンドポイントでトークンを取得
 * @param {string} baseUrl - ベースURL
 * @returns {string} JWTトークン
 */
export function getToken(baseUrl) {
  const res = http.post(
    `${baseUrl}/users/auth/test-login`,
    JSON.stringify({}),
    {
      headers: { "Content-Type": "application/json" },
    }
  );
  if (res.status !== 200) {
    fail(`test-login failed: ${res.status} ${String(res.body).slice(0, 200)}`);
  }
  const token = res.json("token");
  if (!token) {
    fail("no token in test-login response");
  }
  return token;
}

/**
 * Authorizationヘッダー付きのリクエストオプションを返す
 * @param {string} token - JWTトークン
 * @returns {object} ヘッダーを含むオプションオブジェクト
 */
export function withAuth(token) {
  return {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  };
}

/**
 * プロファイルに応じたステージ設定を返す
 * ポートフォリオ用途: 同時5-10ユーザー想定、DB接続プール5を考慮
 * @param {string} profile - 'pr' または 'nightly'
 * @returns {array} ステージ設定の配列
 */
export function profileStages(profile) {
  if (profile === "nightly") {
    // Nightly: 最大10 VU
    return [
      { duration: "20s", target: 3 },
      { duration: "1m", target: 10 },
      { duration: "2m", target: 10 },
      { duration: "20s", target: 0 },
    ];
  }
  // PR: 最大5 VU
  return [
    { duration: "20s", target: 2 },
    { duration: "1m", target: 5 },
    { duration: "1m30s", target: 5 },
    { duration: "20s", target: 0 },
  ];
}

/**
 * 検索結果から先頭の単語IDを取得
 * @param {string} baseUrl - ベースURL
 * @param {string} token - JWTトークン
 * @param {string} q - 検索クエリ
 * @param {string} sortBy - ソート条件
 * @param {number} maxRetries - 最大リトライ回数
 * @returns {number} 単語ID
 */
export function pickWordId(baseUrl, token, q, sortBy, maxRetries = 5) {
  let currentQuery = q;
  for (let i = 0; i < maxRetries; i++) {
    // ランダムなページを使用して多様な単語を選ぶ
    const page = Math.floor(Math.random() * 10) + 1; // 1-10ページ
    const res = http.get(
      `${baseUrl}/words?search=${encodeURIComponent(
        currentQuery
      )}&sortBy=${sortBy}&order=asc&page=${page}&limit=10`,
      withAuth(token)
    );
    if (res.status !== 200) {
      if (i === maxRetries - 1) {
        fail(
          `pickWordId list failed: ${res.status} ${String(res.body).slice(
            0,
            200
          )}`
        );
      }
      sleep(0.5); // リトライ前に少し待つ
      continue;
    }
    const data = res.json();
    if (!data.words || data.words.length === 0) {
      // 検索結果が空の場合、別のクエリを試す
      if (i < maxRetries - 1) {
        currentQuery = randomSearchQuery(); // 新しいクエリを生成
        sleep(0.1);
        continue;
      }
      fail("no words found in response");
    }
    // ランダムに単語を選ぶ（配列内から）
    const randomIndex = Math.floor(Math.random() * data.words.length);
    const item = data.words[randomIndex];
    if (item && item.id) {
      return item.id;
    }
    // 検索結果が空の場合、別のクエリを試す
    if (i < maxRetries - 1) {
      currentQuery = randomSearchQuery(); // 新しいクエリを生成
      sleep(0.1);
    }
  }
  fail("no word id found after retries");
}

/**
 * 共通の閾値設定を返す
 * @returns {object} 閾値設定オブジェクト
 */
export function commonThresholds() {
  return {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<200"],
  };
}

/**
 * ユーザーらしさのための think time（1秒）
 */
export function think() {
  sleep(1);
}

/**
 * ランダムな検索クエリを生成
 * @returns {string} 検索クエリ
 */
export function randomSearchQuery() {
  // より確実に結果が返るクエリに変更
  const queries = [
    "a", // 最も一般的な文字
    "the", // 最も一般的な単語
    "and", // 一般的な単語
    "or", // 一般的な単語
    "but", // 一般的な単語
    "test", // テスト用の単語
    "word", // 一般的な単語
    "able", // 一般的な単語
    "example", // 一般的な単語
    "search", // 一般的な単語
  ];
  return queries[Math.floor(Math.random() * queries.length)];
}

/**
 * ランダムなソート条件を返す
 * @returns {string} 'name' または 'registrationCount'
 */
export function randomSortBy() {
  return Math.random() < 0.5 ? "name" : "registrationCount";
}

/**
 * ランダムなクイズパラメータを生成
 * テストユーザーのクォータ制限（デフォルト20）を考慮して、最大20に制限
 * テストユーザーは登録単語が0件のため、isRegisteredWords: 0（全単語）に固定
 * @returns {object} クイズリクエストボディ
 */
export function randomQuizParams() {
  // テストユーザーのクォータ制限を考慮して、10, 20のみに制限（50, 100を削除）
  const questionCounts = [10, 20];
  const questionCount =
    questionCounts[Math.floor(Math.random() * questionCounts.length)];

  return {
    questionCount,
    isSaveResult: false,
    // テストユーザーは登録単語が0件のため、0（全単語）に固定
    isRegisteredWords: 0,
    correctRate: 1,
    attentionLevelList: [],
    partsOfSpeeches: [1],
    isIdioms: Math.random() < 0.5 ? 0 : 1,
    isSpecialCharacters: Math.random() < 0.5 ? 0 : 1,
  };
}

/**
 * ランダムな単語登録用の単語名を生成（ユニーク性を保つため）
 * @param {number} vuId - VU ID（ユニーク性のため）
 * @returns {string} 単語名
 */
export function randomWordName(vuId) {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  return `k6-test-word-${vuId}-${timestamp}-${random}`;
}

/**
 * 登録済み単語をリセット（テスト実行前に呼び出す）
 * @param {string} baseUrl - ベースURL
 * @param {string} token - JWTトークン
 * @param {number} maxUnregister - 最大解除数（デフォルト200、上限まで）
 */
export function resetRegisteredWords(baseUrl, token, maxUnregister = 200) {
  let totalUnregistered = 0;
  let page = 1;
  const limit = 50; // 1ページあたりの取得数

  // ページネーションで全件取得して解除
  while (totalUnregistered < maxUnregister) {
    const res = http.get(
      `${baseUrl}/words?sortBy=register&order=asc&page=${page}&limit=${limit}`,
      withAuth(token)
    );

    if (res.status !== 200) {
      console.warn(
        `Failed to get registered words (page ${page}): ${res.status} ${String(
          res.body
        ).slice(0, 200)}`
      );
      break;
    }

    const data = res.json();
    if (!data.words || data.words.length === 0) {
      // これ以上登録済み単語がない
      break;
    }

    // 登録済み単語を解除
    for (const word of data.words) {
      if (totalUnregistered >= maxUnregister) {
        break;
      }
      if (word.isRegistered && word.id) {
        const unregisterRes = http.post(
          `${baseUrl}/words/register`,
          JSON.stringify({ wordId: word.id, isRegistered: false }),
          withAuth(token)
        );
        if (unregisterRes.status >= 200 && unregisterRes.status < 300) {
          totalUnregistered++;
        }
      }
    }

    // 取得した件数がlimit未満なら、これが最後のページ
    if (data.words.length < limit) {
      break;
    }

    page++;
  }

  if (totalUnregistered > 0) {
    console.log(`Reset ${totalUnregistered} registered words`);
  }
}
