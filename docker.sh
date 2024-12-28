#!/bin/bash

# bash docker.sh up dev
# bash docker.sh down dev
# bash docker.sh exec backend dev
# bash docker.sh exec frontend dev
# bash docker.sh db dev
# bash docker.sh up prod
# bash docker.sh down prod
# bash docker.sh exec backend prod
# bash docker.sh exec frontend prod
# bash docker.sh db prod

if [ "$3" = "production" ]; then
  ENV_FILE="backend/.env.production"
else
  ENV_FILE="backend/.env.development"
fi

case "$1" in
  up)
    docker compose --env-file $ENV_FILE up
    ;;
  down)
    docker compose --env-file $ENV_FILE down
    ;;
  exec)
    case "$2" in
      backend)
        docker compose --env-file $ENV_FILE exec backend bash
        ;;
      frontend)
        docker compose --env-file $ENV_FILE exec frontend bash
        ;;
      *)
        echo "Usage: $0 exec {backend|frontend} {development|production}"
        exit 1
        ;;
    esac
    ;;
  db)
    docker compose --env-file $ENV_FILE exec -it db psql -U postgres -d db
    ;;
  *)
    echo "Usage: $0 {up|down|exec|db} {backend|frontend} {development|production}"
    exit 1
    ;;
esac
