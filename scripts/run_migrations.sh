#!/usr/bin/env bash

export $(cat ./config/.env | grep -v ^# | xargs)
docker build -f ./Dockerfile.migrations --force-rm -t meli-migrations .
docker run --network="host" --env-file ./config/.env -it meli-migrations