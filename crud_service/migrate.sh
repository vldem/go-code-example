#!/bin/sh

export MIGRATION_DIR="./migrations"
export DB_DSN="host=localhost port=6432 user=user password=password dbname=gohw sslmode=disable"

if [ "$1" = "--dryrun" ]; then
    goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" status
elif [ "$1" = "--down" ]
then
    goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" down
else
    goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" up
fi
