#!/usr/bin/env bash

export $(cat configs/.env | grep -v ^# | xargs)
docker build -f Dockerfile.migrations --force-rm -t meli-migrations .
docker run --network="host" --env-file configs/.env -it meli-migrations