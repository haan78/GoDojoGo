#!/bin/bash
set -euo pipefail

MYSQL_ROOT_PASSWORD=$(docker exec mysql8042 printenv MYSQL_ROOT_PASSWORD)

if [ -z "$MYSQL_ROOT_PASSWORD" ]; then
    echo "MYSQL_ROOT_PASSWORD is empty"
    exit 1
fi

docker exec mysql8042 mysqldump \
    -uroot \
    -p"$MYSQL_ROOT_PASSWORD" \
    --databases godojogo \
    --routines \
    --triggers \
    --events \
    --single-transaction \
    --set-gtid-purged=OFF \
    > mysql-init/01-schema.sql

echo "Dump completed successfully"
