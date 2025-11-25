**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `GET /words?q=`（検索：DB 最適化前の性能測定）

**計測メタデータ**

- **Date (JST)**: 2025-11-25 16:59:11 JST
- **Git SHA**: `e672e06`
- **Go Version**: 1.25.4
- **k6 Version**: v1.3.0
- **DB Version**: PostgreSQL（ローカル）
- **DB Size**: words=要確認件数 / users=固定シード
- **Seed Type**: 固定シード
- **ラベル**: **before_idx**（インデックス最適化前）

**テスト環境**

- Frontend: ローカル（未使用）
- Backend: ローカル（Go 1.25.4）
- DB: PostgreSQL（ローカル）
- リージョン: ローカル
- キャッシュ: なし（DB 直）

**ワークロード（k6）**

- シナリオ: DB before/after 比較（単語検索）
- エグゼキューター: `ramping-arrival-rate`（到着率ベース）
- ステージ: `30s → 10 iter/s`, `1m30s → 30 iter/s`, `30s → 0`
- 最大 VU: 100（事前割り当て: 20）
- ThinkTime: 0.2s
- 検索クエリ: ランダム選択（`able,test,go,ai,cat,run,play,have,make,good`）
- 検索パラメータ: `sortBy=name`
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search}: p(95)<200ms`
  - `http_req_failed{endpoint:search}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- 平均 RPS: **16.0 req/s**
- チェック成功率: **100%** (2399/2399)
- エラー率: **0.00%** (0/2400)
- 実行時間: **約 2 分 30 秒**

**エンドポイント別メトリクス**

| Endpoint         | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 平均 RPS | SLO |
| ---------------- | -------: | -------: | -------: | ---------: | -------: | :-: |
| GET /words?q=... |   **12** |   **17** |   **74** |  **0.00%** | **16.0** | ✅  |

**HTTP メトリクス**

- http_req_duration (search):
  - avg: **11.49ms**
  - min: **3.91ms**
  - med: **11.53ms**
  - max: **74.02ms**
  - p(90): **16.61ms**
  - p(95): **17.42ms**
- http_req_failed (search): **0.00%** (0/2399)
- http_reqs: **2400** (16.0 req/s)
- イテレーション: **2399** (16.0 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (2399/2399)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得
2. `GET /words?q={ランダム}&sortBy=name` - 複数の検索クエリでランダムに検索

**解釈 / ボトルネック**

- **SLO 達成**: p95 latency が **17.42ms** で 200ms を大きく下回っており、SLO を達成しています
- **エラー率**: **0.00%** で安定しています
- **平均レイテンシ**: **11.49ms** で、検索処理としては非常に良好な性能です
- **最大レイテンシ**: **74.02ms** で、200ms を大きく下回っており、インデックスが効いていることを示しています
- **到着率ベースのテスト**: ramping-arrival-rate エグゼキューターを使用しており、一定の到着率（RPS）を維持しながら負荷をかけています
- **ランダム検索クエリ**: 複数の検索クエリをランダムに使用することで、より現実的な負荷を再現しています
- **インデックス最適化前のベースライン**: この結果は、インデックス最適化前の性能ベースラインとして記録されます
- **after テストとの比較**: [db_after_local.md](db_after_local.md) の結果と比較することで、インデックス最適化の効果を測定できます
- 詳細な比較グラフは [db_index_before_after.png](../db_index_before_after.png) を参照してください

**再現手順**

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）
# インデックス最適化前の状態で実行

export BASE_URL="http://localhost:8080"
export LABEL="before_idx"
npm --prefix bench run k6:c:db:before

# 結果は bench/k6/out/db_before.json に出力されます
```

**注意事項**

- このテストはインデックス最適化**前**の性能を測定します
- インデックス最適化**後**のテスト（`k6:c:db:after`）と比較することで、最適化の効果を確認できます
- ラベル（`LABEL=before_idx`）により、結果を区別できます
- 検索クエリはランダムに選択されるため、実行ごとに異なるクエリパターンが使用されます
