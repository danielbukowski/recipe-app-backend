#!/bin/bash

set -e 

DATABASE_NAME=db

echo "Preparing environment for development..."

echo "Running docker compose up..."
docker compose -f ./dev-docker-compose.yml up --force-recreate -d

echo "Checking the health of database..."
until [[ "`docker inspect -f {{.State.Health.Status}} ${DATABASE_NAME}`" == "healthy" ]]; do
    sleep 0.5;
done;

echo "Applying migrations on database..."
goose -dir ./db/migrations postgres postgres://postgres:password@localhost:5432/db?sslmode=disable up

echo "Pinging the api..."
if [[ $(curl -s -m 5 -w "%{http_code}" -o /dev/null -X GET http://localhost:8080/api/v1/health) != 200 ]]; then
    echo "Failed to ping the API"
    docker compose -f ./dev-docker-compose.yml down 
    exit 1
fi

echo "All done successfully!"