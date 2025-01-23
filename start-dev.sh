#!/bin/bash

database_name=db

echo "Starting to prepare environment for development..."

cat .env.example > .env

docker compose -f ./dev-docker-compose.yml up --force-recreate -d

until [ "`docker inspect -f {{.State.Health.Status}} $database_name`"=="healthy" ]; do
    sleep 0.5;
done;

make db-up