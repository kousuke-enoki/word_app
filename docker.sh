#!/usr/bin/env bash
set -euo pipefail

# ----------------------------------------
# 使い方メモ
# ----------------------------------------
# bash docker.sh up dev|prod
# bash docker.sh down dev|prod
# bash docker.sh exec backend|frontend dev|prod
# bash docker.sh db [dev|prod]        ← ここを改善
# bash docker.sh import dev|prod
# ----------------------------------------

CMD="${1:-}"          # up / down / exec / db / import
TARGET="${2:-}"       # backend / frontend / dev / prod …
ENV_NAME="${2:-dev}"  # dev をデフォルトに
[[ "$CMD" == "exec" ]] && ENV_NAME="${3:-dev}"

# .env 切替
if [[ "$ENV_NAME" =~ ^(prod|production)$ ]]; then
  ENV_FILE="backend/.env.production"
else
  ENV_NAME="dev"
  ENV_FILE="backend/.env.development"
fi

# 開発時だけ確認しやすいように echo
echo "[INFO] CMD=$CMD  TARGET=$TARGET  ENV=$ENV_NAME  ENV_FILE=$ENV_FILE"

# -------------------------------------------------
# 補助: .env からホスト環境変数へ取り込み
# -------------------------------------------------
if [[ -f "$ENV_FILE" ]]; then
  # shellcheck disable=SC1090
  set -o allexport
  source "$ENV_FILE"
  set +o allexport
fi

# -------------------------------------------------
# コマンド本体
# -------------------------------------------------
case "$CMD" in
  up)
    docker compose --env-file "$ENV_FILE" up db backend frontend
    ;;

  up_d)
    docker compose --env-file "$ENV_FILE" up -d db backend frontend
    ;;

  down)
    docker compose --env-file "$ENV_FILE" down
    ;;

  exec)
    if [[ "$TARGET" =~ ^(backend|frontend)$ ]]; then
      docker compose --env-file "$ENV_FILE" exec "$TARGET" bash
    else
      echo "Usage: $0 exec {backend|frontend} {dev|prod}"
      exit 1
    fi
    ;;

  db)
    # DB_USER / DB_NAME は .env に無ければ postgres / postgres
    DB_USER="${DB_USER:-postgres}"
    DB_NAME="${DB_NAME:-postgres}"
    docker compose --env-file "$ENV_FILE" exec -it db \
      psql -U "$DB_USER" -d "$DB_NAME"
    ;;

  import)
    docker compose --env-file "$ENV_FILE" up -d db
    docker compose --env-file "$ENV_FILE" --profile import run --rm dict-import \
      -file=/data/jmdict.json -workers=4
    ;;

  *)
    cat <<EOF
Usage:
  $0 up {dev|prod}
  $0 down {dev|prod}
  $0 exec {backend|frontend} {dev|prod}
  $0 db [dev|prod]
  $0 import {dev|prod}
EOF
    exit 1
    ;;
esac
