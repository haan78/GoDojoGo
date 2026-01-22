#!/bin/sh

docker exec mysql8042 mysqldump -uroot --databases godojogo --routines --triggers --events --single-transaction --set-gtid-purged=OFF > mysql-init/01-schema.sql
