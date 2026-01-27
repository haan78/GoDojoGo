#!/usr/bin
MYSQL_ROOT_PASSWORD=`docker exec mysql8042 printenv MYSQL_ROOT_PASSWORD`
cmd="mysqldump -uroot -p$MYSQL_ROOT_PASSWORD --databases godojogo --routines --triggers --events --single-transaction --set-gtid-purged=OFF"
docker exec mysql8042 $cmd > mysql-init/01-schema.sql
