# コマンドリスト
# bash docker.sh up dev
# bash docker.sh up prod
# bash docker.sh down dev
# bash docker.sh down prod
# bash docker.sh exec backend dev
# bash docker.sh exec backend prod
# bash docker.sh exec frontend dev
# bash docker.sh exec frontend prod
# bash docker.sh db dev
# bash docker.sh db prod
# bash docker.sh import dev
# bash docker.sh import prod

CMD="$1"            # up / down / exec / db / import
TARGET="$2"         # backend / frontend / dev / prod など
ENV_NAME="$2"       # import コマンドでは dev|prod がここに入る
if [[ "$CMD" == "exec" ]]; then
  ENV_NAME="$3"
fi

# .env ファイルの切替
if [[ "$ENV_NAME" == "prod" || "$ENV_NAME" == "production" ]]; then
  ENV_FILE="backend/.env.production"
else
  ENV_FILE="backend/.env.development"
fi

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
    case "$TARGET" in
      backend|frontend)
        docker compose --env-file "$ENV_FILE" exec "$TARGET" bash
        ;;
      *)
        echo "Usage: $0 exec {backend|frontend} {dev|prod}"
        exit 1
        ;;
    esac
    ;;

  db)
    docker compose --env-file "$ENV_FILE" exec -it db psql -U "${DB_USER:-postgres}" -d "${DB_NAME:-postgres}"
    ;;

  import)
    docker compose --env-file "$ENV_FILE" up -d db
    docker compose --env-file "$ENV_FILE" \
      --profile import run --rm dict-import \
      -file=/data/jmdict.json -workers=4
    ;;

  *)
    cat <<EOF
Usage:
  $0 up {dev|prod}
  $0 down {dev|prod}
  $0 exec {backend|frontend} {dev|prod}
  $0 db {dev|prod}
  $0 import {dev|prod}
EOF
    exit 1
    ;;
esac
