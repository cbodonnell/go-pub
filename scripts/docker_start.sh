#!/bin/sh
sudo docker-compose down
sudo ENV=prod docker-compose --env-file="$ENV.env" -f docker-compose.prod.yml up -d
