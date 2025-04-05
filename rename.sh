#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <app_name> <env_var_pref>"
    exit 1
fi

APP_NAME=$1
ENV_VAR_PREF=$2

mv "cmd/app" "cmd/$APP_NAME"

sed -i "s/APP/$ENV_VAR_PREF/g" deploy/app/.env
sed -i "s/APP/$ENV_VAR_PREF/g" deploy/app/.env.example
sed -i "s/app/$APP_NAME/g" deploy/app/docker-compose.yaml
mv "deploy/app" "deploy/$APP_NAME"

sed -i "s/app/$APP_NAME/g" internal/app/app.go
mv "internal/app/app.go" "internal/app/$APP_NAME.go"
mv "internal/app" "internal/$APP_NAME"

sed -i "s/APP/$ENV_VAR_PREF/g" config/config.go
