# Word App(dictQuiz) — Go/Gin/Ent × AWS Lambda × PostgreSQL

**操作性**と**再現可能なパフォーマンス**を証明する SaaS ポートフォリオです。

[![CI](https://img.shields.io/badge/CI-pass-green)](#)
[![k6 p95](https://img.shields.io/badge/k6-p95%3C200ms-blue)](#)
[![Coverage](https://img.shields.io/badge/Go%20Coverage-XX%25-informational)](#)
（順次作成予定）

## Live Demo

- URL: https://word-app-opal.vercel.app/
- Test Login でお試し可能です。

## SLO & Results

- **SLO**: p95 < 200ms, Error rate < 1%
- **Method**: k6（login / words search / quiz generation）
- **Results**:  
  ![search p95](docs/perf/search_p95.png)
  ![error rate](docs/perf/error_rate.png)
  ![DB Index Before/After](docs/perf/db_index_before_after.png)
  （グラフなどの表については順次作成予定）

詳細なパフォーマンステスト結果は [docs/perf/README.md](docs/perf/README.md) を参照してください。

## Architecture

詳細なアーキテクチャドキュメントは [docs/arch/README.md](docs/arch/README.md) を参照してください。

### 概要

本アプリケーションは**クリーンアーキテクチャ**を採用し、レイヤー間の依存関係を明確に分離しています。認証・設定・ユーザー管理などのコア機能は既に UseCase 層へ移行済みで、Word/Quiz/Result 機能は移行進行中です。

## Tech Stack

Go (Gin, Ent), PostgreSQL, Docker, AWS Lambda + API Gateway, Vercel, k6, GitHub Actions

## Quick Start (Local)

```bash
# 1) envを準備（例: backend/.env.development, frontend/.env.development）
# 2) build & up
docker compose build
bash docker.sh up dev
# 3) （初回のみ）辞書インポート(かなり時間かかる)
bash docker.sh import dev
```
