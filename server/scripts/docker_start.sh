#!/bin/sh
docker-compose down
docker-compose --env-file="$ENV.env" -f docker-compose.prod.yml up -d
