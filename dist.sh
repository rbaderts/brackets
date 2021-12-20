#!/bin/sh

(cd database; docker build . -t brackets_postgres)
(cd frontend; npm install; npm run build; docker build . -t brackets_fe)
(./build.sh prod; docker build . -t brackets_be)

docker-compose up




