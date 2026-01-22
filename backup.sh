#!/bin/sh

docker exec mysql8042 mysqldump -uroot -ptv11k52 --databases godojogo --routines --triggers --events --single-transaction --set-gtid-purged=OFF > mysql-init/01-schema.sql
