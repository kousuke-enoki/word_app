#!/bin/bash

if [ "$1" = "production" ]; then
  ENV_FILE="backend/.env.production"
else
  ENV_FILE="backend/.env.development"
fi

docker compose --env-file $ENV_FILE exec backend bash
