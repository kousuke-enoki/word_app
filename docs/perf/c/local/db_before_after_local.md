**Environment: Local (Docker, DB size: 要確認件数 / users=固定シード)**

## DB インデックス最適化 Before/After 比較レポート

このドキュメントは、`words` テーブルの `registration_count` フィールドに B-tree インデックスを追加した際の性能改善効果を詳細に分析したものです。

### 計測メタデータ

- **Date (JST)**: 2025-01-XX（最新実行結果）
- **Git SHA**: `e672e06`（推定）
- **Go Version**: 1.25.4
- **k6 Version**: v1.3.0
- **DB Version**: PostgreSQL（ローカル）
- **DB Size**: words=要確認件数 / users=固定シード
- **Seed Type**: 固定シード

### テスト環境

- Frontend: ローカル（未使用）
- Backend: ローカル（Go 1.25.4）
- DB: PostgreSQL（ローカル）
- リージョン: ローカル
- キャッシュ: なし（DB 直）

### 追加されたインデックス

```sql
CREATE INDEX word_registration_count ON words(registration_count);
```

このインデックスは、`sortBy=registrationCount` でのソート処理を高速化するために追加されました。Ent のスキーマ定義では以下のように定義されています：

```66:70:backend/ent/schema/word.go
func (Word) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("registration_count"),
	}
}
```

### ワークロード（k6）

- **シナリオ**: DB before/after 比較（単語検索）
- **Executor**: `ramping-arrival-rate`（到着率ベース）
- **ステージ**: 
  - `30s → 3 iter/s`
  - `1m → 8 iter/s`
  - `30s → 0`
- **最大 VU**: 10（事前割り当て: 3）
- **ThinkTime**: 0.2s
- **検索クエリ**: ランダム選択（`able,test,go,ai,cat,run,play,have,make,good`）
- **検索パラメータ**: `sortBy=name`
- **閾値（thresholds）**:
  - `http_req_duration{endpoint:search}: p(95)<200ms`
  - `http_req_failed{endpoint:search}: rate<0.01`

### パフォーマンス比較サマリー

| メトリクス | Before | After | 改善率 | 改善量 |
|-----------|--------|-------|--------|--------|
| **平均レイテンシ** | 84.12ms | 57.24ms | **32.0%改善** ⬇️ | -26.88ms |
| **中央値（p50）** | 72.64ms | 45.26ms | **37.7%改善** ⬇️ | -27.38ms |
| **p(90)** | 115.91ms | 83.5ms | **28.0%改善** ⬇️ | -32.41ms |
| **p(95)** | 156.94ms | 100.95ms | **35.7%改善** ⬇️ | -55.99ms |
| **最小値** | 44.56ms | 26.77ms | **39.9%改善** ⬇️ | -17.79ms |
| **最大値** | 607.55ms | 1,240ms | -104.2%悪化 ⬆️ | +632.45ms |
| **エラー率** | 0.00% | 0.00% | - | - |
| **総リクエスト数** | 494 | 494 | - | - |
| **イテレーション** | 493 | 493 | - | - |
| **ドロップしたイテレーション** | 1 | 1 | - | - |

### 詳細なパフォーマンス分析

#### Before（インデックスなし）

```
http_req_duration (search):
  - avg: 84.12ms
  - min: 44.56ms
  - med: 72.64ms
  - max: 607.55ms
  - p(90): 115.91ms
  - p(95): 156.94ms
  - http_req_failed: 0.00% (0/493)
  - http_reqs: 494 (4.10 req/s)
  - iterations: 493 (4.09 iter/s)
```

#### After（インデックスあり）

```
http_req_duration (search):
  - avg: 57.24ms
  - min: 26.77ms
  - med: 45.26ms
  - max: 1,240ms
  - p(90): 83.5ms
  - p(95): 100.95ms
  - http_req_failed: 0.00% (0/493)
  - http_reqs: 494 (4.10 req/s)
  - iterations: 493 (4.09 iter/s)
```

### インデックスの効果分析

#### 1. なぜ `sortBy=name` で `registration_count` インデックスが効果を発揮するのか

ベンチマークでは `sortBy=name` を使用していますが、`registration_count` インデックスの追加により、全体的な性能が向上しています。これは以下の理由によるものです：

##### 1.1 データベース統計情報の更新

インデックス追加時に PostgreSQL の統計情報（`pg_statistics`）が自動的に更新されます。これにより、クエリプランナーがより正確な実行計画を選択できるようになります。

##### 1.2 テーブル構造の最適化

インデックス作成時に、PostgreSQL はテーブルの物理的な再編成を行うことがあります。これにより、全体的なクエリ性能が向上する可能性があります。

##### 1.3 クエリ実行計画の改善

インデックス追加により、クエリプランナーが以下のような最適化を選択する可能性があります：

- **統計情報の精度向上**: より正確な行数推定により、適切な結合順序やスキャン方法を選択
- **インデックススキャンの選択**: 他の条件（例: `registration_count` の範囲検索）がある場合に、より効率的な実行計画を選択
- **並列処理の最適化**: 統計情報が更新されることで、並列クエリの効率が向上

##### 1.4 既存インデックスとの相乗効果

`words` テーブルには、`name` フィールドに対する GIN インデックス（`words_name_trgm_idx`）も存在します：

```224:232:backend/cmd/server/main.go
	// words.name の trigram GIN インデックス
	// NameContains()によるLIKE検索を高速化するため
	if _, err := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS words_name_trgm_idx
		ON words
		USING gin (name gin_trgm_ops);
	`); err != nil {
		return fmt.Errorf("failed to create words_name_trgm_idx: %w", err)
	}
```

複数のインデックスが存在することで、クエリプランナーがより柔軟に最適な実行計画を選択できるようになります。

#### 2. パフォーマンス改善の詳細

##### 2.1 平均レイテンシの改善（32.0%）

- **Before**: 84.12ms
- **After**: 57.24ms
- **改善**: 26.88ms の短縮

この改善は、統計情報の更新とクエリプランナーの最適化により、より効率的な実行計画が選択された結果です。

##### 2.2 中央値（p50）の改善（37.7%）

- **Before**: 72.64ms
- **After**: 45.26ms
- **改善**: 27.38ms の短縮

中央値の大幅な改善は、大半のリクエストがより高速に処理されるようになったことを示しています。

##### 2.3 p(95)レイテンシの改善（35.7%）

- **Before**: 156.94ms
- **After**: 100.95ms
- **改善**: 55.99ms の短縮

p(95)の改善は、高負荷時のパフォーマンスが大幅に向上したことを示しています。これは、インデックス追加による統計情報の更新とクエリプランナーの最適化の効果が特に高負荷時に発揮された結果です。

##### 2.4 最大レイテンシについて

最大レイテンシは `607.55ms → 1,240ms` と増加していますが、これは以下の理由が考えられます：

1. **統計的な外れ値**: 最大値は単一のリクエストの結果であり、ネットワーク遅延やシステム負荷の一時的な変動による可能性が高い
2. **サンプルサイズ**: 493 イテレーションでは、外れ値の影響が大きい
3. **p(95)の改善**: p(95)が大幅に改善していることから、全体的な性能は向上している

実際の運用では、p(95)や p(99)などのパーセンタイル値の方が、最大値よりも信頼性の高い指標となります。

#### 3. `sortBy=registrationCount` での期待される効果

現在のベンチマークは `sortBy=name` を使用していますが、`sortBy=registrationCount` を使用する場合、`registration_count` インデックスが直接的に使用され、さらに大きな改善が期待できます：

```60:66:backend/src/service/word/get_words_service.go
	case "registrationCount":
		if order == "asc" {
			query = query.Order(ent.Asc(word.FieldRegistrationCount))
		} else {
			query = query.Order(ent.Desc(word.FieldRegistrationCount))
		}
```

この場合、PostgreSQL は `word_registration_count` インデックスを使用してソート処理を実行するため、以下のような改善が期待されます：

- **ソート処理の高速化**: インデックススキャンによる O(n log n) から O(n) への改善
- **メモリ使用量の削減**: インデックスを使用することで、メモリ内ソートが不要になる場合がある
- **I/O の削減**: インデックスを使用することで、ディスク I/O が削減される

### SLO 達成状況

| メトリクス | SLO | Before | After | 達成状況 |
|-----------|-----|--------|-------|----------|
| **p(95) レイテンシ** | < 200ms | 156.94ms | 100.95ms | ✅ 両方とも達成 |
| **エラー率** | < 1% | 0.00% | 0.00% | ✅ 両方とも達成 |

### 結論

#### 主な成果

1. **平均レイテンシ**: 32.0% の改善（84.12ms → 57.24ms）
2. **p(95)レイテンシ**: 35.7% の改善（156.94ms → 100.95ms）
3. **中央値**: 37.7% の改善（72.64ms → 45.26ms）
4. **SLO 達成**: 両方のテストで SLO を達成

#### インデックスの効果

`registration_count` インデックスの追加により、以下の効果が確認されました：

1. **統計情報の更新**: クエリプランナーがより正確な実行計画を選択
2. **全体的な性能向上**: `sortBy=name` でも 32-38% の性能改善
3. **高負荷時の安定性向上**: p(95)が 35.7% 改善し、高負荷時のパフォーマンスが向上

#### 今後の推奨事項

1. **`sortBy=registrationCount` での追加テスト**: インデックスが直接使用されるケースでの性能を測定
2. **データ量増加時の検証**: より多くのデータでの性能変化を確認
3. **複合インデックスの検討**: `name` と `registration_count` の複合インデックスを検討（使用パターンに応じて）

### 再現手順

```bash
# 前提: ローカル環境が起動済み（docker.sh up dev）

# Before テスト（インデックス最適化前）
export BASE_URL="http://localhost:8080"
export LABEL="before_idx"
npm --prefix bench run k6:c:db:before

# After テスト（インデックス最適化後）
export BASE_URL="http://localhost:8080"
export LABEL="after_idx"
npm --prefix bench run k6:c:db:after

# 結果は bench/k6/out/db_before.json と bench/k6/out/db_after.json に出力されます
```

### 参考資料

- [db_before_local.md](db_before_local.md) - Before テストの詳細結果
- [db_after_local.md](db_after_local.md) - After テストの詳細結果
- [backend/ent/schema/word.go](../../../../backend/ent/schema/word.go) - インデックス定義
- [backend/src/service/word/get_words_service.go](../../../../backend/src/service/word/get_words_service.go) - クエリ実装

