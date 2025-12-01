**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## SLO & Results

**SLO（目標）**

- p95 latency < **200ms**
- Error rate < **1%**
- 対象: `GET /words?search=`（検索：DB 最適化前の性能測定）

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
- Executor: `ramping-arrival-rate`（到着率ベース）
- ステージ: `30s → 3 iter/s`, `1m → 8 iter/s`, `30s → 0`
- 最大 VU: **10**（事前割り当て: 3）
- ThinkTime: 0.2s
- 検索クエリ: ランダム選択（`able,test,go,ai,cat,run,play,have,make,good`）
- 検索パラメータ: `sortBy=name`
- 閾値（thresholds）:
  - `http_req_duration{endpoint:search}: p(95)<200ms`
  - `http_req_failed{endpoint:search}: rate<0.01`

**結果（要約）**

- SLO: **達成** ✅
- 平均 RPS: **4.11 req/s**
- チェック成功率: **100%** (493/493)
- エラー率: **0.00%** (0/493)
- 実行時間: **約 2 分**
- ドロップしたイテレーション: **1 件**（VU 不足によりスキップ）

**エンドポイント別メトリクス**

| Endpoint            | p50 (ms) | p95 (ms) | p99 (ms) | Error rate | 平均 RPS | SLO |
| ------------------- | -------: | -------: | -------: | ---------: | -------: | :-: |
| GET /words?search=  |    **72** |  **105** |    **-** |  **0.00%** | **4.11** | ✅  |

**HTTP メトリクス**

- http_req_duration (search):
  - avg: **78.16ms**
  - min: **50.17ms**
  - med: **71.53ms**
  - max: **465.28ms**
  - p(90): **96.81ms**
  - p(95): **104.95ms**
- http_req_failed (search): **0.00%** (0/493)
- http_reqs: **494** (4.11 req/s)
- イテレーション: **493** (4.10 iter/s)
- ドロップしたイテレーション: **1** (0.01 iter/s)

**チェック結果**

- ✓ 2xx: **100%** (493/493)

**ネットワーク**

- 受信データ: **1.5 MB** (12 kB/s)
- 送信データ: **146 kB** (1.2 kB/s)

**実行フロー**

1. `POST /users/auth/test-login` - テストログインでトークン取得（setup 関数）
2. `GET /words?search={ランダム}&sortBy=name` - 複数の検索クエリでランダムに検索

**解釈 / ボトルネック**

- **SLO 達成**: p95 latency が **104.95ms** で 200ms を大きく下回っており、SLO を達成しています
- **エラー率**: **0.00%** で安定しています
- **平均レイテンシ**: **78.16ms** で、検索処理としては良好な性能です
- **中央値（p50）**: **71.53ms** と低く、大半のリクエストが高速に処理されています
- **最大レイテンシ**: **465.28ms** は散発的に発生していますが、p95 が 200ms を下回っているため許容範囲内です
- **到着率ベースのテスト**: ramping-arrival-rate executor を使用しており、一定の到着率（RPS）を維持しながら負荷をかけています
- **ランダム検索クエリ**: 複数の検索クエリをランダムに使用することで、より現実的な負荷を再現しています
- **インデックス最適化前のベースライン**: この結果は、インデックス最適化前の性能ベースラインとして記録されます
- **VU 不足**: 1 件のイテレーションがドロップされました（VU 不足により）。これは maxVUs=10 では目標レートを完全に維持するには不十分だったことを示しています
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
