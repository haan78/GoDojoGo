#!/bin/bash

##run those commands below before run script
# eval "$(ssh-agent -s)"
# ssh-add "[path of key file]"

set -euo pipefail

set -a
source .env.server
set +a

HOST="root@$SERVERADDR"

ssh -o BatchMode=yes -o ConnectTimeout=10 -i "$SERVERKEY" "$HOST" \
  "mysqldump -u$SERVERMYSQLUSER -p$SERVERMYSQLPASS \
  --databases godojogo \
  --routines --triggers --events \
  --single-transaction \
  --set-gtid-purged=OFF" > mysql-init/01-schema.sql
